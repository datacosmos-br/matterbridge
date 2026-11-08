package bdiscord

import (
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
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
