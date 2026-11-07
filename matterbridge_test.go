package main

import (
	"path/filepath"
	"testing"

	"github.com/matterbridge-org/matterbridge/gateway/dify"
	"github.com/stretchr/testify/require"
)

func TestLinkIdentity(t *testing.T) {
	path := filepath.Join(t.TempDir(), "identities.sqlite")

	require.NoError(t, linkIdentity(path, "api.local:user1=alice@corp"))

	store, err := dify.OpenIdentityStore(path)
	require.NoError(t, err)

	defer func() { require.NoError(t, store.Close()) }()

	corp, linked, err := store.Resolve("api.local", "user1")
	require.NoError(t, err)
	require.True(t, linked)
	require.Equal(t, "alice@corp", corp)

	require.ErrorContains(t, linkIdentity("", "api.local:user1=alice@corp"), "IdentityStorePath")
	require.Error(t, linkIdentity(path, "not-a-link"))
}
