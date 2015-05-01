package grapho

import "reflect"

var registeredEventTypes = map[string]reflect.Type{}

func RegisterEvent(event Event) {
	registeredEventTypes[event.EventType()] = reflect.TypeOf(event)
}

func NewEventForType(typename string) Event {
	typ, found := registeredEventTypes[typename]
	if !found {
		panic("Event " + typename + " not registered")
	}

	return reflect.New(typ.Elem()).Interface().(Event)
}
