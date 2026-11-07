package gateway

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matterbridge-org/matterbridge/bridge/config"
)

// OptOutStatus is the relay preference a user set with the opt-in/out commands.
type OptOutStatus int64

const (
	OptIn           OptOutStatus = 0 // relay everything (default)
	OptOut          OptOutStatus = 1 // redact text, attachments and identity
	OptOutMediaOnly OptOutStatus = 2 // redact attachments only
)

// Relay preference commands, intercepted only when UserStorePath is set.
const (
	cmdOptIn        = "!optin"
	cmdOptOut       = "!optout"
	cmdOptOutMedia  = "!optoutmedia"
	cmdHelp         = "!help"
	optOutUsername  = "[Opt-out User]"
	redactedTextMsg = "Redacted Text\n"
)

const commandHelp = `!optout - opt out from all message relaying
!optoutmedia - only opt out from relaying attachments
!optin - opt back into chat relaying
!help - display this message`

const userStoreSchema = `
CREATE TABLE IF NOT EXISTS user_preferences (
	account TEXT    NOT NULL,
	user_id TEXT    NOT NULL,
	opt_out INTEGER NOT NULL,
	PRIMARY KEY (account, user_id)
) WITHOUT ROWID;
`

// userStore keeps the relay preferences of users (UserStorePath) in sqlite.
// Users are identified by the account they write from and their user ID on it.
type userStore struct {
	db *sql.DB

	selectStatus, upsertStatus *sql.Stmt
}

func openUserStore(path string) (*userStore, error) {
	db, err := openSQLite(path, userStoreSchema)
	if err != nil {
		return nil, fmt.Errorf("user store: %w", err)
	}

	s := &userStore{db: db}

	err = s.prepareAll()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("user store %s: %w", path, err), s.close())
	}

	return s, nil
}

func (s *userStore) prepareAll() error {
	var err error

	s.selectStatus, err = prepare(s.db, `SELECT opt_out FROM user_preferences WHERE account = ? AND user_id = ?`)
	if err != nil {
		return err
	}

	s.upsertStatus, err = prepare(s.db, `INSERT INTO user_preferences (account, user_id, opt_out) VALUES (?, ?, ?)
		ON CONFLICT (account, user_id) DO UPDATE SET opt_out = excluded.opt_out`)

	return err
}

func (s *userStore) status(account, userID string) (OptOutStatus, error) {
	var status OptOutStatus

	err := s.selectStatus.QueryRowContext(context.Background(), account, userID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return OptIn, nil
	}

	if err != nil {
		return OptIn, fmt.Errorf("read relay preferences of %s on %s: %w", userID, account, err)
	}

	return status, nil
}

func (s *userStore) setStatus(account, userID string, status OptOutStatus) error {
	_, err := s.upsertStatus.ExecContext(context.Background(), account, userID, status)
	if err != nil {
		return fmt.Errorf("store relay preferences of %s on %s: %w", userID, account, err)
	}

	return nil
}

func (s *userStore) close() error {
	return closeAll(s.db, s.selectStatus, s.upsertStatus)
}

// handleUserPreferences answers relay preference commands and redacts the
// messages of users that opted out. It returns true when msg must not be
// relayed. Without UserStorePath it does nothing.
func (r *Router) handleUserPreferences(msg *config.Message) bool {
	if r.userStore == nil {
		return false
	}

	if r.handleCommand(msg) {
		return true
	}

	err := r.redactOptedOut(msg)
	if err != nil {
		// Relaying could leak the content of a user that opted out.
		r.logger.Errorf("Dropping message %s from %s: %s", msg.ID, msg.Account, err)
		return true
	}

	return false
}

// handleCommand answers the relay preference commands. It returns true when
// the message was a command, which is then not relayed.
func (r *Router) handleCommand(msg *config.Message) bool {
	if msg.Event != "" {
		return false
	}

	switch msg.Text {
	case cmdOptIn:
		r.handleOptOutCmd(msg, OptIn)
	case cmdOptOut:
		r.handleOptOutCmd(msg, OptOut)
	case cmdOptOutMedia:
		r.handleOptOutCmd(msg, OptOutMediaOnly)
	case cmdHelp:
		r.replyCmd(msg, commandHelp)
	default:
		return false
	}

	return true
}

func (r *Router) handleOptOutCmd(msg *config.Message, status OptOutStatus) {
	if msg.UserID == "" {
		r.replyCmd(msg, "Cannot set message relay preferences: this bridge does not report user IDs.")
		return
	}

	err := r.userStore.setStatus(msg.Account, msg.UserID, status)
	if err != nil {
		r.logger.Errorf("Setting relay preferences failed: %s", err)
		r.replyCmd(msg, "Error setting message relay preferences, try again later or contact the moderators.")

		return
	}

	r.replyCmd(msg, "Successfully set message relay preferences.")
}

// replyCmd answers a command in the channel it was sent from.
func (r *Router) replyCmd(msg *config.Message, text string) {
	reply := config.Message{
		Text:     text,
		Channel:  msg.Channel,
		Account:  msg.Account,
		Protocol: msg.Protocol,
		Gateway:  msg.Gateway,
		ParentID: msg.ID,
	}

	_, err := r.getBridge(msg.Account).Send(reply)
	if err != nil {
		r.logger.Errorf("Reply to command %q on %s failed: %s", msg.Text, msg.Account, err)
	}
}

// redactOptedOut redacts the message of a user that opted out of relaying.
func (r *Router) redactOptedOut(msg *config.Message) error {
	if msg.UserID == "" {
		return nil
	}

	status, err := r.userStore.status(msg.Account, msg.UserID)
	if err != nil {
		return err
	}

	redactMessage(msg, status)

	return nil
}

// redactMessage applies the relay preference status to msg.
func redactMessage(msg *config.Message, status OptOutStatus) {
	files := msg.Extra["file"]

	switch status {
	case OptIn:
		return
	case OptOut:
		msg.UserID = ""
		msg.Avatar = ""
		msg.Username = optOutUsername

		if msg.Text != "" || hasFileComment(files) {
			msg.Text = redactedTextMsg
		}
	case OptOutMediaOnly:
		for _, f := range files {
			if fi, ok := f.(config.FileInfo); ok && fi.Comment != "" {
				msg.Text += fi.Comment + "\n"
			}
		}
	}

	if len(files) > 0 {
		msg.Text += fmt.Sprintf("Redacted %d Attachment(s)", len(files))
		delete(msg.Extra, "file")
	}
}

func hasFileComment(files []any) bool {
	for _, f := range files {
		if fi, ok := f.(config.FileInfo); ok && fi.Comment != "" {
			return true
		}
	}

	return false
}
