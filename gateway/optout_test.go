package gateway

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/matterbridge-org/matterbridge/bridge/api"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway/bridgemap"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	userStoreConfig = `
[general]
UserStorePath=%q

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
	apiAccount = "api.local"
	optOutUser = "user1"
	userName   = "alice"
)

func newUserStoreRouter(t *testing.T, path string) (*Router, error) {
	t.Helper()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	cfg := config.NewConfigFromString(logger, fmt.Appendf(nil, userStoreConfig, path))

	r, err := NewRouter(logger, cfg, bridgemap.FullMap)
	if err == nil {
		t.Cleanup(func() { disconnectBridges(t, r) })
	}

	return r, err
}

// disconnectBridges stops the HTTP server of the api bridge of r. The other
// bridges were never connected, so they have nothing to stop.
func disconnectBridges(t *testing.T, r *Router) {
	t.Helper()

	apiBridge, ok := r.getBridge(apiAccount).Bridger.(*api.API)
	require.True(t, ok)
	require.NoError(t, apiBridge.Disconnect())
}

func userMessage(text string) *config.Message {
	return &config.Message{
		Text:     text,
		Channel:  "api",
		Account:  apiAccount,
		Protocol: "api",
		UserID:   optOutUser,
		Username: userName,
		Avatar:   "https://example.invalid/alice.png",
		ID:       "1",
	}
}

func lastReply(t *testing.T, r *Router) string {
	t.Helper()

	apiBridge, ok := r.getBridge(apiAccount).Bridger.(*api.API)
	require.True(t, ok)

	replies := apiBridge.Messages.Values()
	require.NotEmpty(t, replies)

	reply, ok := replies[len(replies)-1].(config.Message)
	require.True(t, ok)

	return reply.Text
}

func TestUserCommandsNotInterceptedWithoutUserStore(t *testing.T) {
	r := maketestRouter(testconfig)
	require.Nil(t, r.userStore)

	for _, text := range []string{cmdOptOut, cmdOptOutMedia, cmdOptIn, cmdHelp} {
		msg := &config.Message{Text: text, Account: "irc.freenode", Channel: "#wimtesting", UserID: optOutUser, Username: userName}
		want := *msg

		require.False(t, r.handleUserPreferences(msg))
		require.Equal(t, want, *msg)
	}
}

func TestUserCommandsSetRelayPreferences(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.sqlite")

	r, err := newUserStoreRouter(t, path)
	require.NoError(t, err)
	require.NotNil(t, r.userStore)

	require.True(t, r.handleUserPreferences(userMessage(cmdHelp)))
	require.Equal(t, commandHelp, lastReply(t, r))

	for _, tc := range []struct {
		command string
		status  OptOutStatus
	}{
		{cmdOptOut, OptOut},
		{cmdOptOutMedia, OptOutMediaOnly},
		{cmdOptIn, OptIn},
		{cmdOptOut, OptOut},
	} {
		require.True(t, r.handleUserPreferences(userMessage(tc.command)), tc.command)
		require.Equal(t, "Successfully set message relay preferences.", lastReply(t, r))

		status, statusErr := r.userStore.status(apiAccount, optOutUser)
		require.NoError(t, statusErr)
		require.Equal(t, tc.status, status, tc.command)
	}

	msg := userMessage("secret")
	require.False(t, r.handleUserPreferences(msg))
	require.Equal(t, optOutUsername, msg.Username)
	require.Equal(t, "Redacted Text\n", msg.Text)
	require.NoError(t, r.userStore.close())

	// The preference survives a restart.
	r, err = newUserStoreRouter(t, path)
	require.NoError(t, err)

	defer func() { require.NoError(t, r.userStore.close()) }()

	status, err := r.userStore.status(apiAccount, optOutUser)
	require.NoError(t, err)
	require.Equal(t, OptOut, status)
}

func TestUserStoreOpenFailure(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-directory")
	require.NoError(t, os.WriteFile(file, nil, 0o600))

	r, err := newUserStoreRouter(t, filepath.Join(file, "users.sqlite"))
	require.Error(t, err)
	require.Nil(t, r)
}

func TestRedactMessage(t *testing.T) {
	newMessage := func() *config.Message {
		msg := userMessage("hello")
		msg.Extra = map[string][]any{"file": {
			config.FileInfo{Name: "a.png", Comment: "look"},
			config.FileInfo{Name: "b.png"},
		}}

		return msg
	}

	for _, tc := range []struct {
		name     string
		status   OptOutStatus
		text     string
		username string
		avatar   bool
		files    int
	}{
		{"opt-in", OptIn, "hello", userName, true, 2},
		{"opt-out", OptOut, "Redacted Text\nRedacted 2 Attachment(s)", optOutUsername, false, 0},
		{"opt-out media", OptOutMediaOnly, "hello\nlook\nRedacted 2 Attachment(s)", userName, true, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			msg := newMessage()
			redactMessage(msg, tc.status)

			require.Equal(t, tc.text, msg.Text)
			require.Equal(t, tc.username, msg.Username)
			require.Equal(t, tc.avatar, msg.Avatar != "")
			require.Equal(t, tc.status != OptOut, msg.UserID != "")
			require.Len(t, msg.Extra["file"], tc.files)
		})
	}

	t.Run("opt-out without text keeps an empty text", func(t *testing.T) {
		msg := userMessage("")
		redactMessage(msg, OptOut)
		require.Empty(t, msg.Text)
		require.Equal(t, optOutUsername, msg.Username)
	})
}
