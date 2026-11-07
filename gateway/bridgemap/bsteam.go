// +build !nosteam

package bridgemap

import (
	bsteam "github.com/mspgeek-community/matterbridge/bridge/steam"
)

func init() {
	FullMap["steam"] = bsteam.New
}
