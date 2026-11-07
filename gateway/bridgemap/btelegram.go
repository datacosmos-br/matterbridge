// +build !notelegram

package bridgemap

import (
	btelegram "github.com/mspgeek-community/matterbridge/bridge/telegram"
)

func init() {
	FullMap["telegram"] = btelegram.New
}
