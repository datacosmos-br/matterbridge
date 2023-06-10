// +build !noapi

package bridgemap

import (
	"github.com/mspgeek-community/matterbridge/bridge/api"
)

func init() {
	FullMap["api"] = api.New
}
