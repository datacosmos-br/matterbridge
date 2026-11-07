package teamsmeet

// Receiver accepts Microsoft Graph change notifications for Teams meetings.
type Receiver struct {
	ClientState string
}

func (r *Receiver) Validate(clientState string) bool {
	return r.ClientState == "" || r.ClientState == clientState
}
