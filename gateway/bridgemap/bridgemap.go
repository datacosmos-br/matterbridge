package bridgemap

import (
	"github.com/mspgeek-community/matterbridge/bridge"
)

var (
	FullMap           = map[string]bridge.Factory{}
	UserTypingSupport = map[string]struct{}{}
)
