// +build !nowhatsapp
// +build !whatsappmulti

package bridgemap

import (
	bwhatsapp "github.com/mspgeek-community/matterbridge/bridge/whatsapp"
)

func init() {
	FullMap["whatsapp"] = bwhatsapp.New
}
