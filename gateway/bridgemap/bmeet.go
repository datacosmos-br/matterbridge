// +build !nomeet

package bridgemap

import (
	bmeet "github.com/42wim/matterbridge/bridge/meet"
)

func init() {
	FullMap["meet"] = bmeet.New
}
