// +build !novk

package bridgemap

import (
	bvk "github.com/mspgeek-community/matterbridge/bridge/vk"
)

func init() {
	FullMap["vk"] = bvk.New
}
