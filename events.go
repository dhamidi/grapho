package grapho

type Event interface {
	EventType() string
}
type Events []Event
