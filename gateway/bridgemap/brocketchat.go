// +build !norocketchat

package bridgemap

import (
	brocketchat "github.com/mspgeek-community/matterbridge/bridge/rocketchat"
)

func init() {
	FullMap["rocketchat"] = brocketchat.New
}
