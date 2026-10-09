package wire

const (
	StatusOK    = "ok"
	StatusError = "error"
	StatusEvent = "event"
)

type Outgoing struct {
	ID      string `json:"id,omitempty"`
	Type    string `json:"type"`
	Action  string `json:"action"`
	Status  string `json:"status"`
	Content any    `json:"content,omitempty"`
}
