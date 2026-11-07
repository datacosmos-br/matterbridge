package dify

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/stretchr/testify/require"
)

const (
	testIDs = "/ids.sqlite"
	testURL = "http://dify"
)

func TestLoadSettings(t *testing.T) {
	_, enabled, err := LoadSettings(&config.Protocol{})
	require.NoError(t, err)
	require.False(t, enabled)

	s, enabled, err := LoadSettings(&config.Protocol{DifyURL: "https://dify.example/", DifyAPIKey: "k", IdentityStorePath: testIDs})
	require.NoError(t, err)
	require.True(t, enabled)
	require.Equal(t, DefaultWorkflowTimeout, s.WorkflowTimeout)
	require.Equal(t, DefaultSummaryCommand, s.SummaryCommand)
	require.Equal(t, DefaultSummaryMessages, s.SummaryMessages)

	s, _, err = LoadSettings(&config.Protocol{
		DifyURL: "http://dify:5001", DifyAPIKey: "k", IdentityStorePath: testIDs,
		DifyWorkflowTimeout: 5, DifySummaryCommand: "!resumo", DifySummaryMessages: 10,
	})
	require.NoError(t, err)
	require.Equal(t, 5*time.Second, s.WorkflowTimeout)
	require.Equal(t, "!resumo", s.SummaryCommand)
	require.Equal(t, 10, s.SummaryMessages)

	for name, general := range map[string]config.Protocol{
		"missing key":      {DifyURL: testURL, IdentityStorePath: testIDs},
		"missing identity": {DifyURL: testURL, DifyAPIKey: "k"},
		"relative url":     {DifyURL: "dify/api", DifyAPIKey: "k", IdentityStorePath: testIDs},
		"negative timeout": {DifyURL: testURL, DifyAPIKey: "k", IdentityStorePath: testIDs, DifyWorkflowTimeout: -1},
		"negative size":    {DifyURL: testURL, DifyAPIKey: "k", IdentityStorePath: testIDs, DifySummaryMessages: -1},
	} {
		_, _, err := LoadSettings(&general)
		require.Error(t, err, name)
	}
}

func TestParseLink(t *testing.T) {
	l, err := ParseLink("whatsapp.corp:5511999990000@s.whatsapp.net=alice@corp")
	require.NoError(t, err)
	require.Equal(t, Link{Account: "whatsapp.corp", UserID: "5511999990000@s.whatsapp.net", CorpUser: "alice@corp"}, l)

	for _, bad := range []string{"", "api.local:user", "user=alice", ":user=alice", "api.local:=alice", "api.local:user="} {
		_, err := ParseLink(bad)
		require.Error(t, err, bad)
	}
}

func TestIdentityStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "identities.sqlite")

	s, err := OpenIdentityStore(path)
	require.NoError(t, err)

	_, linked, err := s.Resolve("api.local", "u1")
	require.NoError(t, err)
	require.False(t, linked)

	require.NoError(t, s.Link(Link{Account: "api.local", UserID: "u1", CorpUser: "alice"}))
	require.NoError(t, s.Link(Link{Account: "api.local", UserID: "u1", CorpUser: "alice2"}))
	require.NoError(t, s.Close())

	s, err = OpenIdentityStore(path)
	require.NoError(t, err)

	defer func() { require.NoError(t, s.Close()) }()

	corp, linked, err := s.Resolve("api.local", "u1")
	require.NoError(t, err)
	require.True(t, linked)
	require.Equal(t, "alice2", corp)
}

func TestTranscriptKeepsMostRecent(t *testing.T) {
	tr := NewTranscript(2)
	require.Empty(t, tr.Lines())

	tr.Add(Line{Username: "a", Text: "1"})
	require.Equal(t, []Line{{"a", "1"}}, tr.Lines())

	tr.Add(Line{Username: "b", Text: "2"})
	tr.Add(Line{Username: "c", Text: "3"})
	require.Equal(t, []Line{{"b", "2"}, {"c", "3"}}, tr.Lines())
	require.Equal(t, "b: 2\nc: 3\n", Format(tr.Lines()))
}
