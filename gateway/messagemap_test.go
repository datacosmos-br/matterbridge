package gateway

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway/bridgemap"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// persistentConfig reuses one discord bridge for two channels (#342).
const persistentConfig = `
[general]
PersistentMessageStorePath=%q
%s

[irc.freenode]
server=""
[discord.test]
server=""

[[gateway]]
    name = "bridge1"
    enable=true

    [[gateway.inout]]
    account = "irc.freenode"
    channel = "#wimtesting"

    [[gateway.inout]]
    account = "discord.test"
    channel = "general"

    [[gateway.inout]]
    account = "discord.test"
    channel = "general2"
`

func newPersistentRouter(t *testing.T, storePath string) (*Router, error) {
	t.Helper()

	return newPersistentRouterWith(t, storePath, "")
}

// newPersistentRouterWith adds the extra [general] settings to persistentConfig.
func newPersistentRouterWith(t *testing.T, storePath, general string) (*Router, error) {
	t.Helper()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	cfg := config.NewConfigFromString(logger, fmt.Appendf(nil, persistentConfig, storePath, general))

	return NewRouter(logger, cfg, bridgemap.FullMap)
}

func closeRouter(t *testing.T, r *Router) {
	t.Helper()

	require.NoError(t, r.Close())
}

const (
	discordProtocol = "discord"
	discordName     = "test"
	generalChannel  = "generaldiscord.test"
	general2Channel = "general2discord.test"
)

var persistentDestinations = []*BrMsgID{
	{Protocol: discordProtocol, DestName: discordName, ChannelID: generalChannel, ID: discordProtocol + " 100"},
	{Protocol: discordProtocol, DestName: discordName, ChannelID: general2Channel, ID: discordProtocol + " 200"},
}

func TestPersistentMessageMapSurvivesRestart(t *testing.T) {
	dir := t.TempDir()

	r, err := newPersistentRouter(t, dir)
	require.NoError(t, err)
	require.IsType(t, &sqliteMessageMap{}, r.Gateways["bridge1"].messages)
	require.NoError(t, r.Gateways["bridge1"].messages.add("irc 1", nil, persistentDestinations))
	closeRouter(t, r)
	require.FileExists(t, filepath.Join(dir, "bridge1.sqlite"))

	r, err = newPersistentRouter(t, dir)

	require.NoError(t, err)

	defer closeRouter(t, r)

	ids, found, err := r.Gateways["bridge1"].messages.get("irc 1")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, persistentDestinations, ids)
}

func TestPersistentMessageMapCanonicalLookup(t *testing.T) {
	r, err := newPersistentRouter(t, t.TempDir())

	require.NoError(t, err)

	defer closeRouter(t, r)

	gw := r.Gateways["bridge1"]
	require.NoError(t, gw.messages.add("irc 1", nil, persistentDestinations))

	require.Equal(t, "irc 1", gw.FindCanonicalMsgID("irc", "1", ""))
	require.Equal(t, "irc 1", gw.FindCanonicalMsgID(discordProtocol, "200", general2Channel))
	require.Empty(t, gw.FindCanonicalMsgID(discordProtocol, "300", general2Channel))
	// The same message ID in another channel is a different message (protocols
	// such as Telegram reuse IDs across chats).
	require.Empty(t, gw.FindCanonicalMsgID(discordProtocol, "200", generalChannel))
}

func TestPersistentMessageMapDistinguishesChannels(t *testing.T) {
	r, err := newPersistentRouter(t, t.TempDir())

	require.NoError(t, err)

	defer closeRouter(t, r)

	gw := r.Gateways["bridge1"]
	require.NoError(t, gw.messages.add("irc 1", nil, persistentDestinations))

	dest := gw.Bridges["discord.test"]
	require.Equal(t, "100", gw.getDestMsgID("irc 1", dest, gw.Channels[generalChannel]))
	require.Equal(t, "200", gw.getDestMsgID("irc 1", dest, gw.Channels[general2Channel]))
}

func TestPersistentMessageMapOpenFailure(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-directory")
	require.NoError(t, os.WriteFile(file, nil, 0o600))

	r, err := newPersistentRouter(t, file)
	require.Error(t, err)
	require.Nil(t, r)
}

func TestMessageMapDefaultsToMemory(t *testing.T) {
	r := maketestRouter(testconfig)
	require.IsType(t, &lruMessageMap{}, r.Gateways["bridge1"].messages)
}

// destinationsOf returns one destination relayed for canonical message n.
func destinationsOf(n int) []*BrMsgID {
	return []*BrMsgID{{
		Protocol: discordProtocol, DestName: discordName, ChannelID: generalChannel,
		ID: fmt.Sprintf("%s %d", discordProtocol, 100+n),
	}}
}

func countRows(t *testing.T, m *sqliteMessageMap, table string) int {
	t.Helper()

	var n int
	require.NoError(t, m.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM "+table).Scan(&n))

	return n
}

func TestPersistentMessageMapPrunesOldest(t *testing.T) {
	r, err := newPersistentRouterWith(t, t.TempDir(), "PersistentMessageStoreSize=2")
	require.NoError(t, err)

	defer closeRouter(t, r)

	gw := r.Gateways["bridge1"]
	m, ok := gw.messages.(*sqliteMessageMap)
	require.True(t, ok)

	require.NoError(t, m.add("irc 1", nil, destinationsOf(1)))
	require.NoError(t, m.add("irc 2", nil, destinationsOf(2)))
	// Re-adding irc 1 makes irc 2 the oldest.
	require.NoError(t, m.add("irc 1", nil, destinationsOf(1)))
	require.NoError(t, m.add("irc 3", nil, destinationsOf(3)))

	_, found, err := m.get("irc 2")
	require.NoError(t, err)
	require.False(t, found)
	require.Empty(t, gw.FindCanonicalMsgID(discordProtocol, "102", generalChannel))

	require.Equal(t, "irc 1", gw.FindCanonicalMsgID(discordProtocol, "101", generalChannel))
	require.Equal(t, "irc 3", gw.FindCanonicalMsgID(discordProtocol, "103", generalChannel))
	require.Equal(t, 2, countRows(t, m, "messages"))
	require.Equal(t, 2, countRows(t, m, "downstream"))
}

func TestPersistentMessageMapDefaultSize(t *testing.T) {
	r, err := newPersistentRouter(t, t.TempDir())
	require.NoError(t, err)

	defer closeRouter(t, r)

	m, ok := r.Gateways["bridge1"].messages.(*sqliteMessageMap)
	require.True(t, ok)
	require.Equal(t, messageCacheSize, m.size)
}

func TestPersistentMessageMapRejectsNegativeSize(t *testing.T) {
	r, err := newPersistentRouterWith(t, t.TempDir(), "PersistentMessageStoreSize=-1")
	require.ErrorContains(t, err, "size must be positive")
	require.Nil(t, r)
}

func TestPersistentMessageMapMigratesSeq(t *testing.T) {
	dir := t.TempDir()

	db, err := openSQLite(messageMapPath(dir, "bridge1"), `
CREATE TABLE messages (canonical TEXT PRIMARY KEY, destinations TEXT NOT NULL) WITHOUT ROWID;
INSERT INTO messages VALUES ('irc 1', '[]');`)
	require.NoError(t, err)
	require.NoError(t, db.Close())

	r, err := newPersistentRouterWith(t, dir, "PersistentMessageStoreSize=1")
	require.NoError(t, err)

	defer closeRouter(t, r)

	m, ok := r.Gateways["bridge1"].messages.(*sqliteMessageMap)
	require.True(t, ok)
	require.NoError(t, m.add("irc 2", nil, destinationsOf(2)))

	_, found, err := m.get("irc 1")
	require.NoError(t, err)
	require.False(t, found)
	require.Equal(t, 1, countRows(t, m, "messages"))
}

func TestNewRouterClosesStoresOnFailure(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "not-a-directory")
	require.NoError(t, os.WriteFile(file, nil, 0o600))

	r, err := newPersistentRouterWith(t, dir, fmt.Sprintf("UserStorePath=%q", filepath.Join(file, "users.sqlite")))
	require.ErrorContains(t, err, "user store")
	require.Nil(t, r)

	// sqlite removes the WAL file when the last connection closes.
	require.FileExists(t, filepath.Join(dir, "bridge1.sqlite"))
	require.NoFileExists(t, filepath.Join(dir, "bridge1.sqlite-wal"))
}
