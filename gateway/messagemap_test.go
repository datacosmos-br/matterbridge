package gateway

import (
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

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	cfg := config.NewConfigFromString(logger, fmt.Appendf(nil, persistentConfig, storePath))

	return NewRouter(logger, cfg, bridgemap.FullMap)
}

func closeMessageMaps(t *testing.T, r *Router) {
	t.Helper()

	for _, gw := range r.Gateways {
		require.NoError(t, gw.messages.close())
	}
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
	require.NoError(t, r.Gateways["bridge1"].messages.add("irc 1", persistentDestinations))
	closeMessageMaps(t, r)
	require.FileExists(t, filepath.Join(dir, "bridge1.sqlite"))

	r, err = newPersistentRouter(t, dir)

	require.NoError(t, err)
	defer closeMessageMaps(t, r)

	ids, found, err := r.Gateways["bridge1"].messages.get("irc 1")
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, persistentDestinations, ids)
}

func TestPersistentMessageMapCanonicalLookup(t *testing.T) {
	r, err := newPersistentRouter(t, t.TempDir())

	require.NoError(t, err)
	defer closeMessageMaps(t, r)

	gw := r.Gateways["bridge1"]
	require.NoError(t, gw.messages.add("irc 1", persistentDestinations))

	require.Equal(t, "irc 1", gw.FindCanonicalMsgID("irc", "1"))
	require.Equal(t, "irc 1", gw.FindCanonicalMsgID(discordProtocol, "200"))
	require.Empty(t, gw.FindCanonicalMsgID(discordProtocol, "300"))
}

func TestPersistentMessageMapDistinguishesChannels(t *testing.T) {
	r, err := newPersistentRouter(t, t.TempDir())

	require.NoError(t, err)
	defer closeMessageMaps(t, r)

	gw := r.Gateways["bridge1"]
	require.NoError(t, gw.messages.add("irc 1", persistentDestinations))

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
