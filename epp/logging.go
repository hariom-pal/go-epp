package epp

import "time"

// EventType identifies a safe SDK event. Events never include raw XML or secrets.
type EventType string

const (
	EventConnect EventType = "connect"
	EventClose   EventType = "close"
	EventLogin   EventType = "login"
	EventLogout  EventType = "logout"
	EventCommand EventType = "command"
)

// Event contains redaction-safe SDK telemetry for application logging.
type Event struct {
	Type       EventType
	Command    string
	Duration   time.Duration
	ResultCode int
	ClientTRID string
	ServerTRID string
	Err        error
}

// Logger receives redaction-safe SDK events.
type Logger interface {
	EPPEvent(Event)
}
