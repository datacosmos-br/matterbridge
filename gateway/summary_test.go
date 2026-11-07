package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/matterbridge-org/matterbridge/bridge/api"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway/bridgemap"
	"github.com/matterbridge-org/matterbridge/gateway/dify"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	difyConfig = `
[general]
DifyURL=%q
DifyAPIKey=%q
IdentityStorePath=%q
DifySummaryMessages=2

[api.local]
BindAddress="127.0.0.1:0"
[irc.freenode]
server=""

[[gateway]]
    name = "bridge1"
    enable=true

    [[gateway.inout]]
    account = "api.local"
    channel = "api"

    [[gateway.inout]]
    account = "irc.freenode"
    channel = "#wimtesting"
`
	difyKey      = "app-secret"
	difyCorpUser = "alice@corp"
	difyGateway  = "bridge1"
	ircAccount   = "irc.freenode"
)

// difyRun is one workflow run request received by the test Dify server.
type difyRun struct {
	auth, idempotency, contentType string
	body                           difyRunBody
}

// difyRunBody is the JSON body of a workflow run request.
type difyRunBody struct {
	Inputs       map[string]string `json:"inputs"`
	ResponseMode string            `json:"response_mode"`
	User         string            `json:"user"`
}

// difyServer is a Dify workflow API answering every run with status and body.
type difyServer struct {
	*httptest.Server

	mu   sync.Mutex
	runs []difyRun
}

func newDifyServer(t *testing.T, status int, body string) *difyServer {
	t.Helper()

	s := &difyServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/workflows/run" {
			http.NotFound(w, r)
			return
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		run := difyRun{auth: r.Header.Get("Authorization"), idempotency: r.Header.Get("Idempotency-Key"), contentType: r.Header.Get("Content-Type")}

		err = json.Unmarshal(raw, &run.body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.mu.Lock()
		s.runs = append(s.runs, run)
		s.mu.Unlock()

		w.WriteHeader(status)
		_, _ = io.WriteString(w, body)
	}))
	t.Cleanup(s.Close)

	return s
}

func (s *difyServer) recorded() []difyRun {
	s.mu.Lock()
	defer s.mu.Unlock()

	return append([]difyRun(nil), s.runs...)
}

func newDifyRouter(t *testing.T, difyURL, identityPath string) (*Router, error) {
	t.Helper()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	cfg := config.NewConfigFromString(logger, fmt.Appendf(nil, difyConfig, difyURL, difyKey, identityPath))

	return NewRouter(logger, cfg, bridgemap.FullMap)
}

// linkedDifyRouter returns a router whose identity store links optOutUser.
func linkedDifyRouter(t *testing.T, difyURL string) *Router {
	t.Helper()

	path := filepath.Join(t.TempDir(), "identities.sqlite")

	r, err := newDifyRouter(t, difyURL, path)
	require.NoError(t, err)
	require.NotNil(t, r.summary)
	t.Cleanup(func() { require.NoError(t, r.summary.identities.Close()) })

	require.NoError(t, r.summary.identities.Link(dify.Link{Account: apiAccount, UserID: optOutUser, CorpUser: difyCorpUser}))

	return r
}

func gatewayMessage(text, id string) *config.Message {
	msg := userMessage(text)
	msg.ID = id
	msg.Gateway = difyGateway

	return msg
}

// apiReplies returns the texts the api bridge was asked to send.
func apiReplies(t *testing.T, r *Router) []string {
	t.Helper()

	apiBridge, ok := r.getBridge(apiAccount).Bridger.(*api.API)
	require.True(t, ok)

	apiBridge.Lock()
	defer apiBridge.Unlock()

	var texts []string

	for _, v := range apiBridge.Messages.Values() {
		reply, isMsg := v.(config.Message)
		require.True(t, isMsg)

		texts = append(texts, reply.Text)
	}

	return texts
}

func waitReply(t *testing.T, r *Router) string {
	t.Helper()

	var replies []string

	require.Eventually(t, func() bool {
		replies = apiReplies(t, r)
		return len(replies) > 0
	}, 5*time.Second, 10*time.Millisecond)
	require.Len(t, replies, 1)

	return replies[0]
}

func TestSummaryRelaysWorkflowOutput(t *testing.T) {
	srv := newDifyServer(t, http.StatusOK, `{"workflow_run_id":"w1","data":{"status":"succeeded","outputs":{"text":"Alice said hello."}}}`)
	r := linkedDifyRouter(t, srv.URL)
	gw := r.Gateways[difyGateway]

	for i, text := range []string{"first", "second", "third"} {
		msg := gatewayMessage(text, strconv.Itoa(i))
		gw.recordTranscript(msg)
	}

	// Messages of channels outside the gateway are not recorded.
	gw.recordTranscript(&config.Message{Text: "elsewhere", Account: ircAccount, Channel: "#other", Username: "bob"})

	require.True(t, r.handleSummaryCommand(gatewayMessage(dify.DefaultSummaryCommand, "42")))
	require.Equal(t, "Alice said hello.", waitReply(t, r))

	runs := srv.recorded()
	require.Len(t, runs, 1)
	require.Equal(t, "Bearer "+difyKey, runs[0].auth)
	require.Equal(t, "application/json", runs[0].contentType)
	require.Equal(t, difyGateway+" api 42", runs[0].idempotency)
	require.Equal(t, difyRunBody{
		Inputs: map[string]string{
			"channel":      difyGateway,
			"transcript":   "alice: second\nalice: third\n",
			"requested_by": difyCorpUser,
		},
		ResponseMode: "blocking",
		User:         difyCorpUser,
	}, runs[0].body)
}

func TestSummaryUnlinkedUser(t *testing.T) {
	srv := newDifyServer(t, http.StatusOK, `{}`)
	r := linkedDifyRouter(t, srv.URL)

	msg := gatewayMessage(dify.DefaultSummaryCommand, "1")
	msg.UserID = "stranger"

	require.True(t, r.handleSummaryCommand(msg))
	require.Equal(t, []string{summaryNotLinkedMsg}, apiReplies(t, r))
	require.Empty(t, srv.recorded())
}

func TestSummaryWorkflowFailures(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"client error", http.StatusBadRequest, `{"code":"invalid_param"}`},
		{"server error", http.StatusInternalServerError, `oops`},
		{"malformed", http.StatusOK, `not json`},
		{"missing text", http.StatusOK, `{"data":{"status":"succeeded","outputs":{}}}`},
		{"failed run", http.StatusOK, `{"data":{"status":"failed","error":"boom"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := newDifyServer(t, tc.status, tc.body)
			r := linkedDifyRouter(t, srv.URL)

			require.True(t, r.handleSummaryCommand(gatewayMessage(dify.DefaultSummaryCommand, "7")))
			require.Equal(t, summaryFailedMsg, waitReply(t, r))
			require.Len(t, srv.recorded(), 1)
		})
	}
}

func TestSummaryDisabledWithoutDifyURL(t *testing.T) {
	r := maketestRouter(testconfig)
	require.Nil(t, r.summary)

	for _, gw := range r.Gateways {
		require.Nil(t, gw.transcript)
		gw.recordTranscript(gatewayMessage("hello", "1"))
	}

	msg := gatewayMessage(dify.DefaultSummaryCommand, "1")
	want := *msg

	require.False(t, r.handleSummaryCommand(msg))
	require.Equal(t, want, *msg)
}

func TestSummaryInvalidConfigFailsAtStartup(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	cfg := config.NewConfigFromString(logger, []byte(`
[general]
DifyURL="http://127.0.0.1:1"
IdentityStorePath="/nonexistent/identities.sqlite"

[api.local]
BindAddress="127.0.0.1:0"

[[gateway]]
    name = "bridge1"
    enable=true

    [[gateway.inout]]
    account = "api.local"
    channel = "api"
`))

	r, err := NewRouter(logger, cfg, bridgemap.FullMap)
	require.ErrorContains(t, err, "DifyAPIKey is empty")
	require.Nil(t, r)
}
