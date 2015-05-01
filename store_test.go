package grapho

import (
	"log"
	"os"
	"testing"
)

var backend EventStorage

func init() {
	os.Remove("_test/events.log")
	onDisk, err := NewEventsOnDisk("_test/events.log")
	if err != nil {
		log.Fatal(err)
	}
	backend = onDisk
}

func Test_Store_StoresEvents(t *testing.T) {
	store := NewEventStore(backend)
	event := &PostDraftedEvent{
		Id:    "slug",
		Title: "title",
		Body:  "body",
	}
	events := Events{event, event}
	if err := store.Store(events); err != nil {
		t.Fatal(err)
	}

	count := 0
	store.ReplayFunc(func(event Event) error {
		evt, ok := event.(*PostDraftedEvent)
		if !ok {
			t.Errorf("got %T; want %T", event, evt)
		} else {
			count++
		}
		return nil
	})

	if want := len(events); count != want {
		t.Errorf("count = %d; want %d", count, want)
	}
}
