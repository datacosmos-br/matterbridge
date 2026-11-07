package gateway

import (
	"fmt"
	"io"
	"maps"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/matterbridge-org/matterbridge/bridge"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway/bridgemap"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	testUser      = "user"
	deletedSource = config.EventMsgDelete + ":src-1"
)

// recordingProtocol is a bridge protocol that keeps what the gateway sends to
// it and answers new messages with a fresh ID.
const recordingProtocol = "rec"

type recordingBridge struct {
	*bridge.Config

	mu     sync.Mutex
	sent   []config.Message
	nextID int
}

func (b *recordingBridge) Send(msg config.Message) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.sent = append(b.sent, msg)

	if msg.Event == config.EventMsgDelete || msg.ID != "" {
		return msg.ID, nil
	}

	b.nextID++

	return fmt.Sprintf("%s-%d", b.Account, b.nextID), nil
}

func (b *recordingBridge) Connect() error                       { return nil }
func (b *recordingBridge) JoinChannel(config.ChannelInfo) error { return nil }
func (b *recordingBridge) Disconnect() error                    { return nil }
func (b *recordingBridge) SanitizeNick(*config.Message) error   { return nil }

func (b *recordingBridge) take() []config.Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	sent := b.sent
	b.sent = nil

	return sent
}

const deleteConfig = `
[general]
%s

[rec.a]
server=""
[rec.b]
server=""
[rec.c]
server=""

[[gateway]]
    name = "gw"
    enable = true

    [[gateway.inout]]
    account = "rec.a"
    channel = "a"

    [[gateway.inout]]
    account = "rec.b"
    channel = "b"

    [[gateway.inout]]
    account = "rec.c"
    channel = "c"
`

func newDeleteRouter(t *testing.T, general string) (*Router, map[string]*recordingBridge) {
	t.Helper()

	recorders := make(map[string]*recordingBridge)
	factories := maps.Clone(bridgemap.FullMap)
	factories[recordingProtocol] = func(cfg *bridge.Config) bridge.Bridger {
		b := &recordingBridge{Config: cfg}
		recorders[cfg.Account] = b

		return b
	}

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	cfg := config.NewConfigFromString(logger, fmt.Appendf(nil, deleteConfig, general))

	r, err := NewRouter(logger, cfg, factories)
	require.NoError(t, err)
	t.Cleanup(func() {
		for _, gw := range r.Gateways {
			require.NoError(t, gw.messages.close())
		}
	})

	return r, recorders
}

func recMsg(account, channel, id, event, text string) config.Message {
	return config.Message{Account: account, Channel: channel, ID: id, Event: event, Text: text, Username: testUser}
}

// relayNew relays a new message from rec.a and returns its copies on rec.b
// and rec.c.
func relayNew(t *testing.T, r *Router, recs map[string]*recordingBridge) (string, string) {
	t.Helper()

	r.routeMessage(recMsg("rec.a", "a", "src-1", "", "hello"))

	require.Len(t, recs["rec.a"].take(), 0)

	return "rec.b-1", "rec.c-1"
}

// sentIDs returns the event and target ID of each message sent.
func sentIDs(msgs []config.Message) []string {
	ids := make([]string, 0, len(msgs))
	for _, msg := range msgs {
		ids = append(ids, msg.Event+":"+msg.ID)
	}

	return ids
}

func TestDeleteFromDestinationReachesSourceAndOtherDestinations(t *testing.T) {
	r, recs := newDeleteRouter(t, "")
	idB, idC := relayNew(t, r, recs)
	recs["rec.b"].take()
	recs["rec.c"].take()

	r.routeMessage(recMsg("rec.b", "b", idB, config.EventMsgDelete, config.EventMsgDelete))

	require.Equal(t, []string{deletedSource}, sentIDs(recs["rec.a"].take()))
	require.Equal(t, []string{"msg_delete:" + idC}, sentIDs(recs["rec.c"].take()))
	require.Empty(t, recs["rec.b"].take())
}

func TestEditFromDestinationDoesNotReachSource(t *testing.T) {
	r, recs := newDeleteRouter(t, "")
	idB, idC := relayNew(t, r, recs)
	recs["rec.b"].take()
	recs["rec.c"].take()

	r.routeMessage(recMsg("rec.b", "b", idB, "", "edited"))

	require.Empty(t, recs["rec.a"].take())
	require.Equal(t, []string{":" + idC}, sentIDs(recs["rec.c"].take()))

	// The edit did not replace the record: a later delete still reaches the
	// source.
	r.routeMessage(recMsg("rec.c", "c", idC, config.EventMsgDelete, config.EventMsgDelete))
	require.Equal(t, []string{deletedSource}, sentIDs(recs["rec.a"].take()))
}

func TestEditFromSourceKeepsUpstreamBehaviour(t *testing.T) {
	r, recs := newDeleteRouter(t, "")
	idB, idC := relayNew(t, r, recs)
	recs["rec.b"].take()
	recs["rec.c"].take()

	r.routeMessage(recMsg("rec.a", "a", "src-1", "", "edited"))

	require.Equal(t, []string{":" + idB}, sentIDs(recs["rec.b"].take()))
	require.Equal(t, []string{":" + idC}, sentIDs(recs["rec.c"].take()))
	require.Empty(t, recs["rec.a"].take())
}

func TestEchoedDeleteIsDropped(t *testing.T) {
	r, recs := newDeleteRouter(t, "")
	idB, idC := relayNew(t, r, recs)
	recs["rec.b"].take()
	recs["rec.c"].take()

	r.routeMessage(recMsg("rec.b", "b", idB, config.EventMsgDelete, config.EventMsgDelete))
	recs["rec.a"].take()
	recs["rec.c"].take()

	// rec.a and rec.c report the deletes the gateway made.
	r.routeMessage(recMsg("rec.a", "a", "src-1", config.EventMsgDelete, config.EventMsgDelete))
	r.routeMessage(recMsg("rec.c", "c", idC, config.EventMsgDelete, config.EventMsgDelete))

	for _, rec := range recs {
		require.Empty(t, rec.take())
	}
}

func TestDeleteAfterEchoWindowIsRelayed(t *testing.T) {
	r, recs := newDeleteRouter(t, "DeleteEchoWindow=1")
	idB, idC := relayNew(t, r, recs)

	r.routeMessage(recMsg("rec.b", "b", idB, config.EventMsgDelete, config.EventMsgDelete))
	recs["rec.a"].take()
	recs["rec.c"].take()

	time.Sleep(1100 * time.Millisecond)

	// Past the window, a delete of the source is a new delete again.
	r.routeMessage(recMsg("rec.a", "a", "src-1", config.EventMsgDelete, config.EventMsgDelete))
	require.Equal(t, []string{"msg_delete:" + idC}, sentIDs(recs["rec.c"].take()))
}

func TestDeleteEchoWindowValidation(t *testing.T) {
	window, err := deleteEchoWindow(0)
	require.NoError(t, err)
	require.Equal(t, defaultDeleteEchoWindow, window)

	window, err = deleteEchoWindow(3)
	require.NoError(t, err)
	require.Equal(t, "3s", window.String())

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	cfg := config.NewConfigFromString(logger, fmt.Appendf(nil, deleteConfig, "DeleteEchoWindow=-1"))
	r, err := NewRouter(logger, cfg, bridgemap.FullMap)
	require.ErrorContains(t, err, "DeleteEchoWindow")
	require.Nil(t, r)
}

func TestDeletePropagationAcrossPersistentReopen(t *testing.T) {
	dir := t.TempDir()
	general := fmt.Sprintf("PersistentMessageStorePath=%q", dir)

	r, recs := newDeleteRouter(t, general)
	idB, idC := relayNew(t, r, recs)
	require.NoError(t, r.Gateways["gw"].messages.close())

	r, recs = newDeleteRouter(t, general)
	require.IsType(t, &sqliteMessageMap{}, r.Gateways["gw"].messages)

	r.routeMessage(recMsg("rec.b", "b", idB, "", "edited"))
	require.Empty(t, recs["rec.a"].take())
	require.Equal(t, []string{":" + idC}, sentIDs(recs["rec.c"].take()))

	r.routeMessage(recMsg("rec.b", "b", idB, config.EventMsgDelete, config.EventMsgDelete))
	require.Equal(t, []string{deletedSource}, sentIDs(recs["rec.a"].take()))
	require.Equal(t, []string{"msg_delete:" + idC}, sentIDs(recs["rec.c"].take()))
}

// legacyMessageMapSchema is the messages table written before the source
// column existed.
const legacyMessageMapSchema = `
CREATE TABLE messages (canonical TEXT PRIMARY KEY, destinations TEXT NOT NULL) WITHOUT ROWID;
INSERT INTO messages VALUES ('rec old', '[]');
`

func TestPersistentMessageMapAddsSourceColumn(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gw.sqlite")
	db, err := openSQLite(path, legacyMessageMapSchema)
	require.NoError(t, err)
	require.NoError(t, db.Close())

	m, err := openSQLiteMessageMap(path, messageCacheSize)

	require.NoError(t, err)
	defer func() { require.NoError(t, m.close()) }()

	_, found, err := m.source("rec old")
	require.NoError(t, err)
	require.False(t, found)

	src := &BrMsgID{Protocol: recordingProtocol, DestName: "a", ChannelID: "arec.a", ID: "rec new"}
	require.NoError(t, m.add("rec new", src, nil))

	source, found, err := m.source("rec new")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, src, source)
}
