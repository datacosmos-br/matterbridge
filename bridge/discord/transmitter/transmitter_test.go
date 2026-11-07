package transmitter

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookMessageURI(t *testing.T) {
	tr := New(nil, "guild", "title", false)
	tr.AddWebhook("channel", &discordgo.Webhook{
		ID:        "hookid",
		Type:      discordgo.WebhookTypeIncoming,
		GuildID:   "guild",
		ChannelID: "channel",
		Name:      "title",
		Token:     "hooktoken",
	})

	uri, err := tr.webhookMessageURI("channel", "1234")
	require.NoError(t, err)
	assert.Equal(t, discordgo.EndpointWebhookToken("hookid", "hooktoken")+"/messages/1234", uri)
}

func TestDeleteWithoutWebhook(t *testing.T) {
	tr := New(nil, "guild", "title", false)

	err := tr.Delete("unknown-channel", "1234")
	assert.ErrorIs(t, err, ErrWebhookNotFound)
}
