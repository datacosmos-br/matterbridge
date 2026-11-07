package gateway

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/matterbridge-org/matterbridge/bridge"
	"github.com/matterbridge-org/matterbridge/bridge/config"
	"github.com/matterbridge-org/matterbridge/gateway/dify"
	"github.com/matterbridge-org/matterbridge/gateway/samechannel"
	"github.com/sirupsen/logrus"
)

var (
	errGatewayWithoutName = errors.New("Gateway without name found")
	errDuplicateGateway   = errors.New("Gateway with name already exists")
)

type Router struct {
	config.Config
	sync.RWMutex

	BridgeMap        map[string]bridge.Factory
	Gateways         map[string]*Gateway
	Message          chan config.Message
	MattermostPlugin chan config.Message

	// userStore holds user relay preferences when UserStorePath is set
	// (nil otherwise: no command is intercepted and nothing is redacted).
	userStore *userStore

	// deleteEchoes holds the "<protocol> <ID>" of the messages the gateways
	// recently sent a delete for (DeleteEchoWindow).
	deleteEchoes *expirable.LRU[string, struct{}]
	// summary answers the Dify summary command when DifyURL is set
	// (nil otherwise: no command is intercepted and nothing is recorded).
	summary *summarizer

	logger *logrus.Entry
}

// NewRouter initializes a new Matterbridge router for the specified configuration and
// sets up all required gateways. When it fails, the stores it already opened
// are closed and their close errors are joined to the original error.
func NewRouter(rootLogger *logrus.Logger, cfg config.Config, bridgeMap map[string]bridge.Factory) (*Router, error) {
	logger := rootLogger.WithFields(logrus.Fields{"prefix": "router"})

	echoWindow, err := deleteEchoWindow(cfg.BridgeValues().General.DeleteEchoWindow)
	if err != nil {
		return nil, err
	}

	r := &Router{
		deleteEchoes:     expirable.NewLRU[string, struct{}](messageCacheSize, nil, echoWindow),
		Config:           cfg,
		BridgeMap:        bridgeMap,
		Message:          make(chan config.Message),
		MattermostPlugin: make(chan config.Message),
		Gateways:         make(map[string]*Gateway),
		logger:           logger,
	}

	err = r.openGateways(rootLogger, cfg)
	if err != nil {
		return nil, errors.Join(err, r.Close())
	}

	summary, err := newSummarizer(&cfg.BridgeValues().General)
	if err != nil {
		return nil, errors.Join(err, r.Close())
	}

	if summary != nil {
		r.summary = summary
		for _, gw := range r.Gateways {
			gw.transcript = dify.NewTranscript(summary.settings.SummaryMessages)
		}
	}

	return r, nil
}

// Close closes the message map of every gateway and the user store. Call it
// only once no message is routed anymore.
func (r *Router) Close() error {
	var errs []error

	for name, gw := range r.Gateways {
		err := gw.messages.close()
		if err != nil {
			errs = append(errs, fmt.Errorf("gateway %s: %w", name, err))
		}
	}

	if r.userStore != nil {
		err := r.userStore.close()
		if err != nil {
			errs = append(errs, fmt.Errorf("user store: %w", err))
		}
	}

	if r.summary != nil {
		err := r.summary.identities.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("identity store: %w", err))
		}
	}

	return errors.Join(errs...)
}

// Start will connect all gateways belonging to this router and subsequently route messages
// between them.
func (r *Router) Start() error {
	// Deprecating MediaServerUpload. Remove in future v2.1 release
	deprecatedValue, _ := r.GetString("MediaServerUpload")
	if deprecatedValue != "" {
		r.logger.Fatal("MediaServerUpload config option has been deprecated. You should either remove this option from your configuration, or help us document it.")
	}

	m := make(map[string]*bridge.Bridge)
	if len(r.Gateways) == 0 {
		return fmt.Errorf("no [[gateway]] configured. See https://github.com/42wim/matterbridge/wiki/How-to-create-your-config for more info")
	}
	for _, gw := range r.Gateways {
		r.logger.Infof("Parsing gateway %s", gw.Name)
		if len(gw.Bridges) == 0 {
			return fmt.Errorf("no bridges configured for gateway %s. See https://github.com/42wim/matterbridge/wiki/How-to-create-your-config for more info", gw.Name)
		}
		for _, br := range gw.Bridges {
			m[br.Account] = br
		}
	}
	for _, br := range m {
		r.logger.Infof("Starting bridge: %s ", br.Account)
		err := br.Connect()
		if err != nil {
			e := fmt.Errorf("Bridge %s failed to start: %v", br.Account, err)
			if r.disableBridge(br, e) {
				continue
			}
			return e
		}
		err = br.JoinChannels()
		if err != nil {
			e := fmt.Errorf("Bridge %s failed to join channel: %v", br.Account, err)
			if r.disableBridge(br, e) {
				continue
			}
			return e
		}
	}
	// remove unused bridges
	for _, gw := range r.Gateways {
		for i, br := range gw.Bridges {
			if br.Bridger == nil {
				r.logger.Errorf("removing failed bridge %s", i)
				delete(gw.Bridges, i)
			}
		}
	}
	go r.handleReceive()
	//go r.updateChannelMembers()
	return nil
}

// openGateways creates the gateways and the user store; on error, the stores
// opened so far stay referenced by r for NewRouter to close.
func (r *Router) openGateways(rootLogger *logrus.Logger, cfg config.Config) error {
	sgw := samechannel.New(cfg)
	gwconfigs := append(sgw.GetConfig(), cfg.BridgeValues().Gateway...)

	for idx := range gwconfigs {
		entry := &gwconfigs[idx]
		if entry.Name == "" {
			return errGatewayWithoutName
		}

		if !entry.Enable {
			r.logger.Warnf("Gateway %s defined but disabled", entry.Name)
			continue
		}

		if _, ok := r.Gateways[entry.Name]; ok {
			return fmt.Errorf("%w: %s", errDuplicateGateway, entry.Name)
		}

		gw, err := New(rootLogger, entry, r)
		if err != nil {
			return err
		}

		r.Gateways[entry.Name] = gw
	}

	if path := cfg.BridgeValues().General.UserStorePath; path != "" {
		store, err := openUserStore(path)
		if err != nil {
			return err
		}

		r.userStore = store
	}

	return nil
}

// disableBridge returns true and empties a bridge if we have IgnoreFailureOnStart configured
// otherwise returns false
func (r *Router) disableBridge(br *bridge.Bridge, err error) bool {
	if r.BridgeValues().General.IgnoreFailureOnStart {
		r.logger.Error(err)
		// setting this bridge empty
		*br = bridge.Bridge{
			Log: br.Log,
		}
		return true
	}
	return false
}

func (r *Router) getBridge(account string) *bridge.Bridge {
	for _, gw := range r.Gateways {
		if br, ok := gw.Bridges[account]; ok {
			return br
		}
	}
	return nil
}

func (r *Router) handleReceive() {
	for msg := range r.Message {
		r.routeMessage(msg)
	}
}

// routeMessage relays one message received from a bridge.
func (r *Router) routeMessage(msg config.Message) {
	r.handleEventGetChannelMembers(&msg)
	r.handleEventFailure(&msg)
	r.handleEventRejoinChannels(&msg)

	// Set message protocol based on the account it came from
	msg.Protocol = r.getBridge(msg.Account).Protocol

	if r.handleUserPreferences(&msg) || r.handleSummaryCommand(&msg) {
		return
	}

	if r.isDeleteEcho(&msg) {
		r.logger.Debugf("Dropping delete of %s %s from %s: echo of a delete sent by the gateway", msg.Protocol, msg.ID, msg.Account)
		return
	}

	filesHandled := false

	for _, gw := range r.Gateways {
		// record all the message ID's of the different bridges
		var msgIDs []*BrMsgID

		if gw.ignoreMessage(&msg) {
			continue
		}

		msg.Timestamp = time.Now()
		gw.modifyMessage(&msg)
		gw.recordTranscript(&msg)

		if !filesHandled {
			gw.handleFiles(&msg)

			filesHandled = true
		}

		canonical, err := gw.canonicalMsgID(&msg)
		if err != nil {
			r.logger.Errorf("Looking up message map of %s %s failed, dropping it: %s", msg.Protocol, msg.ID, err)
			continue
		}

		for _, br := range gw.Bridges {
			msgIDs = append(msgIDs, gw.handleMessage(&msg, br, canonical)...)
		}

		// Only record a message that is unknown both as a source and as a
		// destination: edits and deletes, from any bridge, keep resolving to
		// the record of the original message.
		if msg.ID != "" && canonical == "" {
			err := gw.messages.add(msg.Protocol+" "+msg.ID, gw.sourceMsgID(&msg), msgIDs)
			if err != nil {
				r.logger.Errorf("Storing message map of %s %s failed: %s", msg.Protocol, msg.ID, err)
			}
		}
	}
}

// updateChannelMembers sends every minute an GetChannelMembers event to all bridges.
func (r *Router) updateChannelMembers() {
	// TODO sleep a minute because slack can take a while
	// fix this by having actually connectionDone events send to the router
	time.Sleep(time.Minute)
	for {
		for _, gw := range r.Gateways {
			for _, br := range gw.Bridges {
				// only for slack now
				if br.Protocol != "slack" {
					continue
				}
				r.logger.Debugf("sending %s to %s", config.EventGetChannelMembers, br.Account)
				if _, err := br.Send(config.Message{Event: config.EventGetChannelMembers}); err != nil {
					r.logger.Errorf("updateChannelMembers: %s", err)
				}
			}
		}
		time.Sleep(time.Minute)
	}
}

var errNegativeDeleteEchoWindow = errors.New("DeleteEchoWindow must not be negative")

// defaultDeleteEchoWindow is the DeleteEchoWindow used when it is unset.
const defaultDeleteEchoWindow = 10 * time.Second

// deleteEchoWindow validates the DeleteEchoWindow setting (seconds, 0 for the
// default).
func deleteEchoWindow(seconds int) (time.Duration, error) {
	if seconds < 0 {
		return 0, fmt.Errorf("%w, got %d", errNegativeDeleteEchoWindow, seconds)
	}

	if seconds == 0 {
		return defaultDeleteEchoWindow, nil
	}

	return time.Duration(seconds) * time.Second, nil
}

// rememberDelete records that a delete of message ID on protocol is being
// sent, so that the delete event the bridge may emit for it is dropped.
func (r *Router) rememberDelete(protocol, id string) {
	r.deleteEchoes.Add(protocol+" "+id, struct{}{})
}

// isDeleteEcho reports whether msg is a delete event for a message the
// gateways sent a delete for within DeleteEchoWindow.
func (r *Router) isDeleteEcho(msg *config.Message) bool {
	if msg.Event != config.EventMsgDelete || msg.ID == "" {
		return false
	}

	return r.deleteEchoes.Contains(msg.Protocol + " " + msg.ID)
}
