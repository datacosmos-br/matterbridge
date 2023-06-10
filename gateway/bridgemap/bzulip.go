// +build !nozulip

package bridgemap

import (
	bzulip "github.com/mspgeek-community/matterbridge/bridge/zulip"
)

func init() {
	FullMap["zulip"] = bzulip.New
}
