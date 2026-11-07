// +build !nosshchat

package bridgemap

import (
	bsshchat "github.com/mspgeek-community/matterbridge/bridge/sshchat"
)

func init() {
	FullMap["sshchat"] = bsshchat.New
}
