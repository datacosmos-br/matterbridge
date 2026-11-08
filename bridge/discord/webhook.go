package bdiscord

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
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

// forumThreadTarget returns the channel owning the webhook and the thread to
// post in for the bridged channel channelID. Forum posts are threads: when
// ForumID is set, the webhook belongs to the forum and channelID is the thread.
// Otherwise the webhook belongs to channelID itself and there is no thread.
func forumThreadTarget(options config.ChannelOptions, channelID string) (string, string) {
	if options.ForumID == "" {
		return channelID, ""
	}

	return options.ForumID, channelID
}

// webhookTarget returns the channel owning the webhook and the thread to post
// in for the bridged channel channelName (whose Discord ID is channelID).
func (b *Bdiscord) webhookTarget(channelName string, channelID string) (string, string) {
	b.channelsMutex.RLock()
	defer b.channelsMutex.RUnlock()

	if ci, ok := b.channelInfoMap[channelName+b.Account]; ok {
		return forumThreadTarget(ci.Options, channelID)
	}

	return channelID, ""
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

func (b *Bdiscord) webhookSendTextOnly(msg *config.Message, channelID string, threadID string) (string, error) {
	msgParts := helper.ClipOrSplitMessage(msg.Text, MessageLength, b.GetString("MessageClipped"), b.GetInt("MessageSplitMaxCount"))
	msgIds := []string{}
	for _, msgPart := range msgParts {
		res, err := b.transmitter.Send(
			channelID,
			threadID,
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

func (b *Bdiscord) webhookSendFilesOnly(msg *config.Message, channelID string, threadID string) error {
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
			threadID,
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
// uploads (from slack, telegram or mattermost) and forum threads.
// Returns messageID and error.
func (b *Bdiscord) webhookSend(msg *config.Message, channelID string, threadID string) (string, error) {
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
		res, err = b.webhookSendTextOnly(msg, channelID, threadID)
	}

	if err == nil && msg.Extra != nil {
		err = b.webhookSendFilesOnly(msg, channelID, threadID)
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

	webhookChannelID, threadID := b.webhookTarget(msg.Channel, channelID)

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
			editErr = b.transmitter.Edit(webhookChannelID, threadID, msgIds[i], &discordgo.WebhookParams{ //nolint:exhaustruct_v5
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

	msgID, err := b.webhookSend(msg, webhookChannelID, threadID)
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

// deleteMessage deletes every part of a message in the bridged channel
// channelName (whose Discord ID is channelID): parts sent through the webhook
// are deleted with the webhook token (no "Manage Messages" permission needed),
// the other parts are deleted by the bot user. A failed part does not stop the
// remaining parts from being deleted; all failures are returned together.
func (b *Bdiscord) deleteMessage(channelName string, channelID string, msgID string) error {
	webhookIDs, botIDs := partitionDeleteIDs(msgID, b.webhookMessages)
	webhookChannelID, threadID := b.webhookTarget(channelName, channelID)

	var errs []error

	for _, id := range webhookIDs {
		err := b.transmitter.Delete(webhookChannelID, threadID, id)
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

var errInvalidWebhookURL = errors.New("failed to parse WebhookURL")

// registerChannelWebhooks registers the WebhookURL of each bridged channel with
// the transmitter and returns the channels whose webhooks are found or created
// automatically (AutoWebhooks).
func (b *Bdiscord) registerChannelWebhooks() ([]string, error) {
	var webhookChannelIDs []string

	for _, channel := range b.Channels {
		channelID := b.getChannelID(channel.Name) // note(qaisjp): this readlocks channelsMutex
		// The webhook of a forum thread belongs to the forum (see forumThreadTarget).
		channelID, _ = forumThreadTarget(channel.Options, channelID)

		// If a WebhookURL was not explicitly provided for this channel,
		// there are two options: just a regular bot message (ugly) or this is should be webhook sent
		if channel.Options.WebhookURL == "" {
			// If it should be webhook sent, we should enforce this via the transmitter
			// Threads of the same forum share the forum's webhook, so add it once.
			if b.useAutoWebhooks && !slices.Contains(webhookChannelIDs, channelID) {
				webhookChannelIDs = append(webhookChannelIDs, channelID)
			}

			continue
		}

		whID, whToken, ok := b.splitURL(channel.Options.WebhookURL)
		if !ok {
			return nil, fmt.Errorf("%w %#v for channel %#v", errInvalidWebhookURL, channel.Options.WebhookURL, channel.ID)
		}

		b.transmitter.AddWebhook(channelID, &discordgo.Webhook{ //nolint:exhaustruct_v5
			ID:        whID,
			Token:     whToken,
			GuildID:   b.guildID,
			ChannelID: channelID,
		})
	}

	return webhookChannelIDs, nil
}
