// +build !noirc

package bridgemap

import (
	birc "github.com/mspgeek-community/matterbridge/bridge/irc"
)

func init() {
	FullMap["irc"] = birc.New
}
