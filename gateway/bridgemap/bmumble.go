// +build !nomumble

package bridgemap

import (
	bmumble "github.com/mspgeek-community/matterbridge/bridge/mumble"
)

func init() {
	FullMap["mumble"] = bmumble.New
}
