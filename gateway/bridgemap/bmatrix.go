// +build !nomatrix

package bridgemap

import (
	bmatrix "github.com/mspgeek-community/matterbridge/bridge/matrix"
)

func init() {
	FullMap["matrix"] = bmatrix.New
}
