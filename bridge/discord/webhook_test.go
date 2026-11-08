package bdiscord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/matterbridge-org/matterbridge/bridge"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/bridge/discord/transmitter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newWebhookMessages(t *testing.T, webhookSent ...string) *lru.Cache[string, struct{}] {
	t.Helper()

	webhookMessages, err := lru.New[string, struct{}](webhookMessageCacheSize)
	require.NoError(t, err)

	for _, id := range webhookSent {
		webhookMessages.Add(id, struct{}{})
	}

	return webhookMessages
}

func TestPartitionDeleteIDs(t *testing.T) {
	testcases := map[string]struct {
		webhookSent        []string
		msgID              string
		expectedWebhookIDs []string
		expectedBotIDs     []string
	}{
		"empty message ID": { //nolint:exhaustruct_v5
			webhookSent: []string{"1"},
			msgID:       "",
		},
		"message sent through the webhook": { //nolint:exhaustruct_v5
			webhookSent:        []string{"1"},
			msgID:              "1",
			expectedWebhookIDs: []string{"1"},
		},
		"message sent by the bot user (e.g. a reply with webhooks enabled)": { //nolint:exhaustruct_v5
			webhookSent:    []string{"1"},
			msgID:          "2",
			expectedBotIDs: []string{"2"},
		},
		"split message sent through the webhook": { //nolint:exhaustruct_v5
			webhookSent:        []string{"1", "2", "3"},
			msgID:              "1;2;3",
			expectedWebhookIDs: []string{"1", "2", "3"},
		},
		"split message sent by the bot user": { //nolint:exhaustruct_v5
			msgID:          "4;5",
			expectedBotIDs: []string{"4", "5"},
		},
		"each part is routed on its own": {
			webhookSent:        []string{"1", "3"},
			msgID:              "1;2;3",
			expectedWebhookIDs: []string{"1", "3"},
			expectedBotIDs:     []string{"2"},
		},
	}

	for testname, testcase := range testcases {
		webhookIDs, botIDs := partitionDeleteIDs(testcase.msgID, newWebhookMessages(t, testcase.webhookSent...))
		assert.Equalf(t, testcase.expectedWebhookIDs, webhookIDs, "Should have routed the expected IDs through the webhook for testcase %s", testname)
		assert.Equalf(t, testcase.expectedBotIDs, botIDs, "Should have routed the expected IDs through the bot user for testcase %s", testname)
	}
}

func TestForumThreadTarget(t *testing.T) {
	testcases := map[string]struct {
		options                  config.ChannelOptions
		channelID                string
		expectedWebhookChannelID string
		expectedThreadID         string
	}{
		"regular channel": { //nolint:exhaustruct_v5
			channelID:                "1",
			expectedWebhookChannelID: "1",
		},
		"regular channel with a webhook URL": { //nolint:exhaustruct_v5
			options:                  config.ChannelOptions{WebhookURL: "https://discord.com/api/webhooks/1234/abcd"}, //nolint:exhaustruct_v5
			channelID:                "1",
			expectedWebhookChannelID: "1",
		},
		"forum thread": {
			options:                  config.ChannelOptions{ForumID: "2"}, //nolint:exhaustruct_v5
			channelID:                "1",
			expectedWebhookChannelID: "2",
			expectedThreadID:         "1",
		},
	}

	for testname, testcase := range testcases {
		webhookChannelID, threadID := forumThreadTarget(testcase.options, testcase.channelID)
		assert.Equalf(t, testcase.expectedWebhookChannelID, webhookChannelID, "Should have used the expected webhook channel for testcase %s", testname)
		assert.Equalf(t, testcase.expectedThreadID, threadID, "Should have used the expected thread for testcase %s", testname)
	}
}

func TestRegisterChannelWebhooksOfForumThreads(t *testing.T) {
	b := &Bdiscord{ //nolint:exhaustruct_v5
		Config: &bridge.Config{Bridge: &bridge.Bridge{Channels: map[string]config.ChannelInfo{ //nolint:exhaustruct_v5
			"thread 11 of forum 1": {Name: "ID:11", Options: config.ChannelOptions{ForumID: "1"}}, //nolint:exhaustruct_v5
			"thread 12 of forum 1": {Name: "ID:12", Options: config.ChannelOptions{ForumID: "1"}}, //nolint:exhaustruct_v5
			"thread 21 of forum 2": {Name: "ID:21", Options: config.ChannelOptions{ //nolint:exhaustruct_v5
				WebhookURL: "https://discord.com/api/webhooks/1234/abcd",
				ForumID:    "2",
			}},
		}}},
		useAutoWebhooks: true,
		transmitter:     transmitter.New(nil, "guild", "title", true),
	}

	webhookChannelIDs, err := b.registerChannelWebhooks()
	require.NoError(t, err)
	assert.Equal(t, []string{"1"}, webhookChannelIDs, "Threads of the same forum should share the forum's auto webhook")

	assert.True(t, b.transmitter.HasWebhook("1234"))
	assert.True(t, b.transmitter.AddWebhook("2", &discordgo.Webhook{ID: "5678"}), //nolint:exhaustruct_v5
		"The WebhookURL of a forum thread should be registered under the forum ID")
	assert.False(t, b.transmitter.AddWebhook("21", &discordgo.Webhook{ID: "5678"}), //nolint:exhaustruct_v5
		"The WebhookURL of a forum thread should not be registered under the thread ID")
}
