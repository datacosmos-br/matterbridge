// +build !nomsteams

package bridgemap

import (
	bmsteams "github.com/mspgeek-community/matterbridge/bridge/msteams"
)

func init() {
	FullMap["msteams"] = bmsteams.New
}
