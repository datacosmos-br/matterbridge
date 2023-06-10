// +build whatsappmulti

package bridgemap

import (
	bwhatsapp "github.com/mspgeek-community/matterbridge/bridge/whatsappmulti"
)

func init() {
	FullMap["whatsapp"] = bwhatsapp.New
}
