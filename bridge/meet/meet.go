package bmeet

import (
	"github.com/42wim/matterbridge/bridge"
	"github.com/42wim/matterbridge/bridge/config"
)

type Bmeet struct {
	*bridge.Config
}

func New(cfg *bridge.Config) bridge.Bridger {
	return &Bmeet{Config: cfg}
}

func (b *Bmeet) Connect() error { b.Log.Info("Google Meet bridge connected (transcription-only mode)"); return nil }
func (b *Bmeet) JoinChannel(channel config.ChannelInfo) error { return nil }
func (b *Bmeet) Send(msg config.Message) (string, error) { return "", nil }
func (b *Bmeet) Disconnect() error { return nil }
