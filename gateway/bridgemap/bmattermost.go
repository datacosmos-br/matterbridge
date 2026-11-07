// +build !nomattermost

package bridgemap

import (
	bmattermost "github.com/mspgeek-community/matterbridge/bridge/mattermost"
)

func init() {
	FullMap["mattermost"] = bmattermost.New
}
