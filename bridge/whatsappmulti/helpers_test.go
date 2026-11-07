package bwhatsapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/types"
)

var (
	reactionOwnJID   = types.NewJID("48111222333", types.DefaultUserServer)
	reactionPeerJID  = types.NewJID("48444555666", types.DefaultUserServer)
	reactionGroupJID = types.NewJID("48111222333-1549986983", types.GroupServer)
)

// reactionInfo returns the info of a reaction sent by sender in chat.
func reactionInfo(chat, sender types.JID) types.MessageInfo {
	return types.MessageInfo{
		MessageSource: types.MessageSource{Chat: chat, Sender: sender},
	}
}

func TestGetReactionParentID(t *testing.T) {
	testcases := map[string]struct {
		key      *waCommon.MessageKey
		info     types.MessageInfo
		expected string
	}{
		"group reaction to another participant's message": {
			key: &waCommon.MessageKey{
				RemoteJID:   new(reactionGroupJID.String()),
				FromMe:      new(false),
				ID:          new("MSGID"),
				Participant: new("48777888999@s.whatsapp.net"),
			},
			info:     reactionInfo(reactionGroupJID, reactionPeerJID),
			expected: "48777888999@s.whatsapp.net/MSGID",
		},
		"group reaction to the reactor's own message": {
			key: &waCommon.MessageKey{
				RemoteJID: new(reactionGroupJID.String()),
				FromMe:    new(true),
				ID:        new("MSGID"),
			},
			info:     reactionInfo(reactionGroupJID, reactionPeerJID),
			expected: "48444555666@s.whatsapp.net/MSGID",
		},
		"direct chat reaction to a message sent by the bridge": {
			key: &waCommon.MessageKey{
				RemoteJID: new(reactionOwnJID.String()),
				FromMe:    new(false),
				ID:        new("MSGID"),
			},
			info:     reactionInfo(reactionPeerJID, reactionPeerJID),
			expected: "48111222333@s.whatsapp.net/MSGID",
		},
		"direct chat reaction to the reactor's own message": {
			key: &waCommon.MessageKey{
				RemoteJID: new(reactionOwnJID.String()),
				FromMe:    new(true),
				ID:        new("MSGID"),
			},
			info:     reactionInfo(reactionPeerJID, reactionPeerJID),
			expected: "48444555666@s.whatsapp.net/MSGID",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			parentID, err := getReactionParentID(tc.key, tc.info, reactionOwnJID)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, parentID)
		})
	}
}

func TestGetReactionParentIDErrors(t *testing.T) {
	testcases := map[string]struct {
		key      *waCommon.MessageKey
		expected error
	}{
		"nil key": {
			expected: errReactionWithoutMessage,
		},
		"empty message ID": {
			key: &waCommon.MessageKey{
				RemoteJID:   new(reactionGroupJID.String()),
				FromMe:      new(false),
				ID:          new(""),
				Participant: new("48777888999@s.whatsapp.net"),
			},
			expected: errReactionWithoutMessage,
		},
		"group reaction without author": {
			key: &waCommon.MessageKey{
				RemoteJID: new(reactionGroupJID.String()),
				FromMe:    new(false),
				ID:        new("MSGID"),
			},
			expected: errReactionWithoutAuthor,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			parentID, err := getReactionParentID(tc.key, reactionInfo(reactionGroupJID, reactionPeerJID), reactionOwnJID)
			require.ErrorIs(t, err, tc.expected)
			assert.Empty(t, parentID)
		})
	}
}

func TestGetReactionParentIDUnparsableParticipant(t *testing.T) {
	key := &waCommon.MessageKey{
		RemoteJID:   new(reactionGroupJID.String()),
		FromMe:      new(false),
		ID:          new("MSGID"),
		Participant: new("48777888999:device@s.whatsapp.net"),
	}

	parentID, err := getReactionParentID(key, reactionInfo(reactionGroupJID, reactionPeerJID), reactionOwnJID)
	require.ErrorContains(t, err, "failed to parse reacted message author")
	assert.Empty(t, parentID)
}
