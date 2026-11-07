package bwhatsapp

import (
	"io"
	"testing"

	"github.com/matterbridge-org/matterbridge/bridge"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var (
	testOwn   = types.NewJID("5511999990000", types.DefaultUserServer)
	testPeer  = types.NewJID("5511999990001", types.DefaultUserServer)
	testAdmin = types.NewJID("5511999990002", types.DefaultUserServer)
	testGroup = types.NewJID("120363000000000001", types.GroupServer)
)

func newTestBridge() *Bwhatsapp {
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	b := new(Bwhatsapp)
	b.Config = new(bridge.Config)
	b.Bridge = new(bridge.Bridge)
	b.Log = logrus.NewEntry(logger)
	b.Account = "whatsapp.test"
	b.Protocol = "whatsapp"
	b.Remote = make(chan config.Message, 1)

	return b
}

// newContextInfo builds a ContextInfo, leaving empty arguments unset.
func newContextInfo(stanzaID, participant string) *waE2E.ContextInfo {
	ci := new(waE2E.ContextInfo)

	if stanzaID != "" {
		ci.StanzaID = &stanzaID
	}

	if participant != "" {
		ci.Participant = &participant
	}

	return ci
}

// newMessageKey builds a MessageKey, leaving empty arguments unset.
func newMessageKey(remoteJID, id, participant string, fromMe bool) *waCommon.MessageKey {
	key := new(waCommon.MessageKey)
	key.FromMe = &fromMe

	if remoteJID != "" {
		key.RemoteJID = &remoteJID
	}

	if id != "" {
		key.ID = &id
	}

	if participant != "" {
		key.Participant = &participant
	}

	return key
}

func newMessageInfo(chat, sender types.JID) types.MessageInfo {
	var info types.MessageInfo

	info.Chat = chat
	info.Sender = sender
	info.IsGroup = chat.Server == types.GroupServer

	return info
}

func newProtocolMessageEvent(info types.MessageInfo, protocolMessage *waE2E.ProtocolMessage) *events.Message {
	evt := new(events.Message)
	evt.Info = info
	evt.Message = new(waE2E.Message)
	evt.Message.ProtocolMessage = protocolMessage

	return evt
}

func TestGetParentIdFromCtx(t *testing.T) {
	testcases := map[string]struct {
		input *waE2E.ContextInfo
		want  string
	}{
		"no context info":        {nil, ""},
		"empty context info":     {newContextInfo("", ""), ""},
		"participant only":       {newContextInfo("", testPeer.String()), ""},
		"stanza ID only":         {newContextInfo("ABCD", ""), ""},
		"stanza and participant": {newContextInfo("ABCD", testPeer.String()), testPeer.String() + "/ABCD"},
		"device participant":     {newContextInfo("ABCD", "5511999990001:3@s.whatsapp.net"), testPeer.String() + "/ABCD"},
		"invalid participant":    {newContextInfo("ABCD", "a.b.c@s.whatsapp.net"), ""},
	}

	for name, tc := range testcases {
		assert.Equalf(t, tc.want, getParentIdFromCtx(tc.input), "This testcase failed: %s", name)
	}
}

func TestGetRevokedMessageID(t *testing.T) {
	groupInfo := newMessageInfo(testGroup, testAdmin)
	directInfo := newMessageInfo(testPeer, testPeer)

	testcases := map[string]struct {
		info    types.MessageInfo
		key     *waCommon.MessageKey
		want    string
		wantErr bool
	}{
		"group admin revokes someone else's message": {
			groupInfo, newMessageKey(testGroup.String(), "ABCD", testPeer.String(), false), testPeer.String() + "/ABCD", false,
		},
		"group member revokes own message": {
			groupInfo, newMessageKey(testGroup.String(), "ABCD", "", true), testAdmin.String() + "/ABCD", false,
		},
		"direct chat revoke": {
			directInfo, newMessageKey(testOwn.String(), "ABCD", "", true), testPeer.String() + "/ABCD", false,
		},
		"missing key":         {directInfo, nil, "", true},
		"missing message ID":  {directInfo, newMessageKey(testOwn.String(), "", "", true), "", true},
		"invalid participant": {groupInfo, newMessageKey(testGroup.String(), "ABCD", "a.b.c@s.whatsapp.net", false), "", true},
	}

	for name, tc := range testcases {
		got, err := getRevokedMessageID(tc.info, tc.key)
		if tc.wantErr {
			assert.Errorf(t, err, "This testcase failed: %s", name)
		} else {
			assert.NoErrorf(t, err, "This testcase failed: %s", name)
		}

		assert.Equalf(t, tc.want, got, "This testcase failed: %s", name)
	}
}

// A revoke in a direct chat carries no participant in its key, and its key
// names the chat from the sender's point of view (our own JID). The deletion
// must be relayed for the peer's message, in the chat it was received from.
func TestHandleMessageRelaysDirectChatRevoke(t *testing.T) {
	b := newTestBridge()

	revoke := new(waE2E.ProtocolMessage)
	revoke.Type = waE2E.ProtocolMessage_REVOKE.Enum()
	revoke.Key = newMessageKey(testOwn.String(), "ABCD", "", true)

	b.handleMessage(newProtocolMessageEvent(newMessageInfo(testPeer, testPeer), revoke))

	require.Len(t, b.Remote, 1)

	rmsg := <-b.Remote
	assert.Equal(t, config.EventMsgDelete, rmsg.Event)
	assert.Equal(t, testPeer.String()+"/ABCD", rmsg.ID)
	assert.Equal(t, testPeer.String(), rmsg.Channel)
}

// A protocol message without a type is not a revoke and must be ignored.
func TestHandleMessageIgnoresUntypedProtocolMessage(t *testing.T) {
	b := newTestBridge()

	b.handleMessage(newProtocolMessageEvent(newMessageInfo(testPeer, testPeer), new(waE2E.ProtocolMessage)))

	assert.Empty(t, b.Remote)
}

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
