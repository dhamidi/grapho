package grapho

type Event interface {
	EventType() string
}

type EventHandler interface {
	HandleEvent(e Event) error
}

type Events []Event

func (self Events) ApplyTo(handler EventHandler) {
	for _, e := range self {
		handler.HandleEvent(e)
	}
}
