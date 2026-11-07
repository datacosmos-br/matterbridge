package gateway

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway/dify"
)

// Replies to the summary command, see docs/dify.md.
const (
	summaryNoUserIDMsg  = "Cannot summarize: this bridge does not report user IDs."
	summaryNoMsgIDMsg   = "Cannot summarize: this bridge does not identify messages."
	summaryNotLinkedMsg = "Cannot summarize: your account is not linked to a corporate user, contact the moderators."
	summaryFailedMsg    = "Summary failed, try again later or contact the moderators."
)

// summarizer answers the summary command through Dify (DifyURL).
type summarizer struct {
	settings   dify.Settings
	client     *dify.Client
	identities *dify.IdentityStore
}

// newSummarizer returns nil when DifyURL is unset, and fails on invalid
// Dify settings or an identity store that cannot be opened.
func newSummarizer(general *config.Protocol) (*summarizer, error) {
	settings, enabled, err := dify.LoadSettings(general)
	if err != nil || !enabled {
		return nil, err
	}

	identities, err := dify.OpenIdentityStore(settings.IdentityPath)
	if err != nil {
		return nil, err
	}

	return &summarizer{settings: settings, client: dify.NewClient(settings), identities: identities}, nil
}

// isSource reports whether msg was received on a sending channel of gw.
func (gw *Gateway) isSource(msg *config.Message) bool {
	if msg.Protocol == apiProtocol && gw.Name != msg.Gateway {
		return false
	}

	channel, ok := gw.Channels[getChannelID(msg)]

	return ok && strings.Contains(channel.Direction, "in")
}

// recordTranscript keeps msg in the transcript of gw when summaries are enabled.
func (gw *Gateway) recordTranscript(msg *config.Message) {
	if gw.transcript == nil || msg.Event != "" || msg.Text == "" || !gw.isSource(msg) {
		return
	}

	gw.transcript.Add(dify.Line{Username: msg.Username, Text: msg.Text})
}

// handleSummaryCommand answers the summary command. It returns true when msg
// was the command, which is then not relayed.
func (r *Router) handleSummaryCommand(msg *config.Message) bool {
	if r.summary == nil || msg.Event != "" || strings.TrimSpace(msg.Text) != r.summary.settings.SummaryCommand {
		return false
	}

	var gateways []*Gateway

	for _, gw := range r.Gateways {
		if gw.isSource(msg) {
			gateways = append(gateways, gw)
		}
	}

	if len(gateways) == 0 {
		return false
	}

	if msg.UserID == "" {
		r.replyCmd(msg, summaryNoUserIDMsg)
		return true
	}

	key, ok := idempotencyKey(msg)
	if !ok {
		r.replyCmd(msg, summaryNoMsgIDMsg)
		return true
	}

	corpUser, linked, err := r.summary.identities.Resolve(msg.Account, msg.UserID)
	if err != nil {
		r.logger.Errorf("Summary requested by %s on %s: %s", msg.UserID, msg.Account, err)
		r.replyCmd(msg, summaryFailedMsg)

		return true
	}

	if !linked {
		r.replyCmd(msg, summaryNotLinkedMsg)
		return true
	}

	for _, gw := range gateways {
		req := dify.SummaryRequest{
			Channel:        gw.Name,
			Transcript:     dify.Format(gw.transcript.Lines()),
			RequestedBy:    corpUser,
			IdempotencyKey: gw.Name + " " + key,
		}

		go r.summarize(*msg, req)
	}

	return true
}

// summarize runs the workflow and answers the requesting message with its output.
func (r *Router) summarize(msg config.Message, req dify.SummaryRequest) {
	text, err := r.summary.client.Summarize(context.Background(), req)
	if err != nil {
		r.logger.Errorf("Summary of %s requested by %s failed: %s", req.Channel, req.RequestedBy, err)
		r.replyCmd(&msg, summaryFailedMsg)

		return
	}

	r.replyCmd(&msg, text)
}

// idempotencyKey identifies the requesting message: its protocol and ID, or,
// for bridges that do not report message IDs (api), a digest of its sender,
// channel and receive time.
func idempotencyKey(msg *config.Message) (string, bool) {
	if msg.ID != "" {
		return msg.Protocol + " " + msg.ID, true
	}

	if msg.Timestamp.IsZero() {
		return "", false
	}

	sum := sha256.Sum256(fmt.Appendf(nil, "%s\x00%s\x00%s\x00%s",
		msg.Account, msg.Channel, msg.UserID, msg.Timestamp.UTC().Format(time.RFC3339Nano)))

	return msg.Protocol + " " + hex.EncodeToString(sum[:]), true
}
