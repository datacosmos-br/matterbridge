// +build !noxmpp

package bridgemap

import (
	bxmpp "github.com/mspgeek-community/matterbridge/bridge/xmpp"
)

func init() {
	FullMap["xmpp"] = bxmpp.New
}
