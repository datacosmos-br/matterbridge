// +build !nonctalk

package bridgemap

import (
	btalk "github.com/mspgeek-community/matterbridge/bridge/nctalk"
)

func init() {
	FullMap["nctalk"] = btalk.New
}
