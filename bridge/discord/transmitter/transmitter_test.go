package transmitter

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookMessageURI(t *testing.T) {
	tr := New(nil, "guild", "title", false)
	tr.AddWebhook("channel", &discordgo.Webhook{ //nolint:exhaustruct_v5
		ID:        "hookid",
		Type:      discordgo.WebhookTypeIncoming,
		GuildID:   "guild",
		ChannelID: "channel",
		Name:      "title",
		Token:     "hooktoken",
	})

	uri, err := tr.webhookMessageURI("channel", "", "1234")
	require.NoError(t, err)
	assert.Equal(t, discordgo.EndpointWebhookToken("hookid", "hooktoken")+"/messages/1234", uri)

	uri, err = tr.webhookMessageURI("channel", "5678", "1234")
	require.NoError(t, err)
	assert.Equal(t, discordgo.EndpointWebhookToken("hookid", "hooktoken")+"/messages/1234?thread_id=5678", uri,
		"A message in a thread (e.g. a forum post) should be addressed through the thread_id parameter")
}

func TestWebhookMessageURIWithoutWebhook(t *testing.T) {
	tr := New(nil, "guild", "title", false)

	_, err := tr.webhookMessageURI("unknown-channel", "5678", "1234")
	assert.ErrorIs(t, err, ErrWebhookNotFound)
}

func TestDeleteWithoutWebhook(t *testing.T) {
	tr := New(nil, "guild", "title", false)

	err := tr.Delete("unknown-channel", "", "1234")
	assert.ErrorIs(t, err, ErrWebhookNotFound)
}
