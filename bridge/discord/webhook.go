package bdiscord

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/bridge/helper"
)

// shouldMessageUseWebhooks checks if have a channel specific webhook, if we're not using auto webhooks
func (b *Bdiscord) shouldMessageUseWebhooks(msg *config.Message) bool {
	if b.useAutoWebhooks {
		return true
	}

	b.channelsMutex.RLock()
	defer b.channelsMutex.RUnlock()
	if ci, ok := b.channelInfoMap[msg.Channel+b.Account]; ok {
		if ci.Options.WebhookURL != "" {
			return true
		}
	}
	return false
}

// maybeGetLocalAvatar checks if UseLocalAvatar contains the message's
// account or protocol, and if so, returns the Discord avatar (if exists)
func (b *Bdiscord) maybeGetLocalAvatar(msg *config.Message) string {
	for _, val := range b.GetStringSlice("UseLocalAvatar") {
		if msg.Protocol != val && msg.Account != val {
			continue
		}

		member, err := b.getGuildMemberByNick(msg.Username)
		if err != nil {
			return ""
		}

		return member.User.AvatarURL("")
	}
	return ""
}

func (b *Bdiscord) webhookSendTextOnly(msg *config.Message, channelID string) (string, error) {
	msgParts := helper.ClipOrSplitMessage(msg.Text, MessageLength, b.GetString("MessageClipped"), b.GetInt("MessageSplitMaxCount"))
	msgIds := []string{}
	for _, msgPart := range msgParts {
		res, err := b.transmitter.Send(
			channelID,
			&discordgo.WebhookParams{
				Content:         msgPart,
				Username:        msg.Username,
				AvatarURL:       msg.Avatar,
				AllowedMentions: b.getAllowedMentions(),
			},
		)
		if err != nil {
			return "", err
		} else {
			b.webhookMessages.Add(res.ID, struct{}{})
			msgIds = append(msgIds, res.ID)
		}
	}
	// Exploit that a discord message ID is actually just a large number, so we encode a list of IDs by separating them with ";".
	return strings.Join(msgIds, ";"), nil
}

func (b *Bdiscord) webhookSendFilesOnly(msg *config.Message, channelID string) error {
	for _, f := range msg.Extra["file"] {
		fi := f.(config.FileInfo) //nolint:forcetypeassert
		file := discordgo.File{
			Name:        fi.Name,
			ContentType: "",
			Reader:      bytes.NewReader(*fi.Data),
		}
		content := fi.Comment

		// The resulting ID is thrown away, so messages with file uploads
		// can neither be edited nor deleted through the webhook later.
		_, err := b.transmitter.Send(
			channelID,
			&discordgo.WebhookParams{
				Username:        msg.Username,
				AvatarURL:       msg.Avatar,
				Files:           []*discordgo.File{&file},
				Content:         content,
				AllowedMentions: b.getAllowedMentions(),
			},
		)
		if err != nil {
			b.Log.Errorf("Could not send file %#v for message %#v: %s", file, msg, err)
			return err
		}
	}
	return nil
}

// webhookSend send one or more message via webhook, taking care of file
// uploads (from slack, telegram or mattermost).
// Returns messageID and error.
func (b *Bdiscord) webhookSend(msg *config.Message, channelID string) (string, error) {
	var (
		res string
		err error
	)

	// If avatar is unset, mutate the message to include the local avatar (but only if settings say we should do this)
	if msg.Avatar == "" {
		msg.Avatar = b.maybeGetLocalAvatar(msg)
	}

	// WebhookParams can have either `Content` or `File`.

	// We can't send empty messages.
	if msg.Text != "" {
		res, err = b.webhookSendTextOnly(msg, channelID)
	}

	if err == nil && msg.Extra != nil {
		err = b.webhookSendFilesOnly(msg, channelID)
	}

	return res, err
}

func (b *Bdiscord) handleEventWebhook(msg *config.Message, channelID string) (string, error) {
	// skip events
	if msg.Event != "" &&
		msg.Event != config.EventUserAction &&
		msg.Event != config.EventJoinLeave &&
		msg.Event != config.EventJoin &&
		msg.Event != config.EventLeave &&
		msg.Event != config.EventTopicChange {
		return "", nil
	}

	// skip empty messages
	if msg.Text == "" && (msg.Extra == nil || len(msg.Extra["file"]) == 0) {
		b.Log.Debugf("Skipping empty message %#v", msg)
		return "", nil
	}

	// discord username must be [0..32] max
	if len(msg.Username) > 32 {
		msg.Username = msg.Username[0:32]
	}

	if msg.ID != "" {
		// Exploit that a discord message ID is actually just a large number, and we encode a list of IDs by separating them with ";".
		msgIds := strings.Split(msg.ID, ";")
		msgParts := helper.ClipOrSplitMessage(b.replaceUserMentions(msg.Text), MessageLength, b.GetString("MessageClipped"), len(msgIds))
		for len(msgParts) < len(msgIds) {
			msgParts = append(msgParts, "((obsoleted by edit))")
		}
		b.Log.Debugf("Editing webhook message")
		var editErr error = nil
		for i := range msgParts {
			// In case of split-messages where some parts remain the same (i.e. only a typo-fix in a huge message), this causes some noop-updates.
			// TODO: Optimize away noop-updates of un-edited messages
			editErr = b.transmitter.Edit(channelID, msgIds[i], &discordgo.WebhookParams{
				Content:         msgParts[i],
				Username:        msg.Username,
				AllowedMentions: b.getAllowedMentions(),
			})
			if editErr != nil {
				break
			}
		}
		if editErr == nil {
			return msg.ID, nil
		}
		b.Log.Errorf("Could not edit webhook message(s): %s; sending as new message(s) instead", editErr)
	}

	b.Log.Debugf("Processing webhook sending for message %#v", msg)
	msg.Text = b.replaceUserMentions(msg.Text)
	msgID, err := b.webhookSend(msg, channelID)
	if err != nil {
		b.Log.Errorf("Could not broadcast via webhook for message %#v: %s", msgID, err)
		return "", err
	}
	return msgID, nil
}

// partitionDeleteIDs splits a (possibly ";"-joined) message ID into the parts
// that were sent through the webhook (as recorded in webhookMessages) and the
// parts that were sent by the bot user. Replies are sent by the bot user even
// when webhooks are enabled, so each part has to be deleted through the path
// that actually sent it.
func partitionDeleteIDs(msgID string, webhookMessages *lru.Cache[string, struct{}]) ([]string, []string) {
	var webhookIDs, botIDs []string

	if msgID == "" {
		return webhookIDs, botIDs
	}

	// Exploit that a discord message ID is actually just a large number, and we encode a list of IDs by separating them with ";".
	for id := range strings.SplitSeq(msgID, ";") {
		if webhookMessages.Contains(id) {
			webhookIDs = append(webhookIDs, id)
		} else {
			botIDs = append(botIDs, id)
		}
	}

	return webhookIDs, botIDs
}

// deleteMessage deletes every part of a message: parts sent through the webhook
// are deleted with the webhook token (no "Manage Messages" permission needed),
// the other parts are deleted by the bot user. A failed part does not stop the
// remaining parts from being deleted; all failures are returned together.
func (b *Bdiscord) deleteMessage(channelID string, msgID string) error {
	webhookIDs, botIDs := partitionDeleteIDs(msgID, b.webhookMessages)

	var errs []error

	for _, id := range webhookIDs {
		err := b.transmitter.Delete(channelID, id)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not delete webhook message %s: %w", id, err))
			continue
		}

		b.webhookMessages.Remove(id)
	}

	for _, id := range botIDs {
		err := b.c.ChannelMessageDelete(channelID, id)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not delete message %s: %w", id, err))
		}
	}

	return errors.Join(errs...)
}
