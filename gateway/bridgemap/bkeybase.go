// +build !nokeybase

package bridgemap

import (
	bkeybase "github.com/mspgeek-community/matterbridge/bridge/keybase"
)

func init() {
	FullMap["keybase"] = bkeybase.New
}
