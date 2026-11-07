// Package dify answers on-demand group summary requests through a Dify
// workflow, for users linked to a corporate identity.
//
// Matterbridge does not run any model itself: it sends the recent transcript
// of a gateway to the Dify workflow API and relays the workflow output.
package dify

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/matterbridge-org/matterbridge/bridge/config"
)

// Defaults of the optional Dify settings, see docs/dify.md.
const (
	DefaultWorkflowTimeout = 30 * time.Second
	DefaultSummaryCommand  = "!summary"
	DefaultSummaryMessages = 50
)

// ErrInvalidSettings reports invalid Dify settings.
var ErrInvalidSettings = errors.New("invalid Dify settings")

// Settings are the validated Dify settings of the [general] section.
type Settings struct {
	URL             *url.URL
	APIKey          string
	WorkflowTimeout time.Duration
	SummaryCommand  string
	SummaryMessages int
	IdentityPath    string
}

// LoadSettings validates the Dify settings of general. It returns false when
// DifyURL is unset, which disables summaries entirely.
func LoadSettings(general *config.Protocol) (Settings, bool, error) {
	if general.DifyURL == "" {
		return Settings{}, false, nil
	}

	u, err := url.Parse(general.DifyURL)
	if err != nil {
		return Settings{}, false, fmt.Errorf("DifyURL: %w", err)
	}

	if (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return Settings{}, false, fmt.Errorf("%w: DifyURL %q: must be an absolute http(s) URL", ErrInvalidSettings, general.DifyURL)
	}

	if general.DifyAPIKey == "" {
		return Settings{}, false, fmt.Errorf("%w: DifyURL is set but DifyAPIKey is empty", ErrInvalidSettings)
	}

	if general.IdentityStorePath == "" {
		return Settings{}, false, fmt.Errorf("%w: DifyURL is set but IdentityStorePath is empty", ErrInvalidSettings)
	}

	s := Settings{
		URL:             u,
		APIKey:          general.DifyAPIKey,
		WorkflowTimeout: DefaultWorkflowTimeout,
		SummaryCommand:  DefaultSummaryCommand,
		SummaryMessages: DefaultSummaryMessages,
		IdentityPath:    general.IdentityStorePath,
	}

	switch {
	case general.DifyWorkflowTimeout < 0:
		return Settings{}, false, fmt.Errorf("%w: DifyWorkflowTimeout %d: must not be negative", ErrInvalidSettings, general.DifyWorkflowTimeout)
	case general.DifyWorkflowTimeout > 0:
		s.WorkflowTimeout = time.Duration(general.DifyWorkflowTimeout) * time.Second
	}

	switch {
	case general.DifySummaryMessages < 0:
		return Settings{}, false, fmt.Errorf("%w: DifySummaryMessages %d: must not be negative", ErrInvalidSettings, general.DifySummaryMessages)
	case general.DifySummaryMessages > 0:
		s.SummaryMessages = general.DifySummaryMessages
	}

	if general.DifySummaryCommand != "" {
		s.SummaryCommand = general.DifySummaryCommand
	}

	return s, true, nil
}
