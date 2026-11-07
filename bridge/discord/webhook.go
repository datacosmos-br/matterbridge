package bdiscord

import (
	"bytes"
	"strings"

	"github.com/42wim/matterbridge/bridge/config"
	"github.com/42wim/matterbridge/bridge/helper"
	"github.com/bwmarrin/discordgo"
)

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
		res, err := b.transmitter.SendThread(
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
			msgIds = append(msgIds, res.ID)
		}
	}
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

		_, err := b.transmitter.SendThread(
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

func (b *Bdiscord) webhookSend(msg *config.Message, channelID string, threadID string) (string, error) {
	var (
		res string
		err error
	)

	if msg.Avatar == "" {
		msg.Avatar = b.maybeGetLocalAvatar(msg)
	}

	if msg.Text != "" {
		res, err = b.webhookSendTextOnly(msg, channelID, threadID)
	}

	if err == nil && msg.Extra != nil {
		err = b.webhookSendFilesOnly(msg, channelID, threadID)
	}

	return res, err
}

func (b *Bdiscord) handleEventWebhook(msg *config.Message, channelID string) (string, error) {
	threadID := ""
	if ci, ok := b.channelInfoMap[msg.Channel+b.Account]; ok {
		if ci.Options.ForumID != "" {
			threadID = channelID
			channelID = ci.Options.ForumID
		}
	}

	if msg.Event == config.EventMsgDelete {
		if msg.ID == "" {
			return "", nil
		}

		err := b.transmitter.DeleteThread(channelID, threadID, msg.ID)
		if err != nil {
			b.Log.Errorf("Could not delete message: %s", err)
			return "", err
		}

		b.Log.Infof("Message deleted successfully")
		return "", nil
	}

	if msg.Event != "" && msg.Event != config.EventUserAction && msg.Event != config.EventJoinLeave && msg.Event != config.EventTopicChange {
		return "", nil
	}

	if msg.Text == "" && (msg.Extra == nil || len(msg.Extra["file"]) == 0) {
		b.Log.Debugf("Skipping empty message %#v", msg)
		return "", nil
	}

	if len(msg.Username) > 32 {
		msg.Username = msg.Username[0:32]
	}

	if msg.ID != "" {
		msgIds := strings.Split(msg.ID, ";")
		msgParts := helper.ClipOrSplitMessage(b.replaceUserMentions(msg.Text), MessageLength, b.GetString("MessageClipped"), len(msgIds))
		for len(msgParts) < len(msgIds) {
			msgParts = append(msgParts, "((obsoleted by edit))")
		}
		b.Log.Debugf("Editing webhook message")
		var editErr error = nil
		for i := range msgParts {
			editErr = b.transmitter.EditThread(channelID, threadID, msgIds[i], &discordgo.WebhookParams{
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
	msgID, err := b.webhookSend(msg, channelID, threadID)
	if err != nil {
		b.Log.Errorf("Could not broadcast via webhook for message %#v: %s", msgID, err)
		return "", err
	}
	return msgID, nil
}
