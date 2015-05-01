package grapho

import (
	"log"
	"os"
	"testing"
)

func backend() EventStorage {
	os.Remove("_test/events.log")
	onDisk, err := NewEventsOnDisk("_test/events.log")
	if err != nil {
		log.Fatal(err)
	}

	return onDisk
}

func Test_Store_StoresEvents(t *testing.T) {
	store := NewEventStore(backend())
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

func Test_Store_ReplayFunc_onlySeesSnapshot(t *testing.T) {
	storage := backend()
	store := NewEventStore(storage)
	events := Events{
		&PostDraftedEvent{"slug-1", "title", "body"},
		&PostDraftedEvent{"slug-2", "title", "body"},
	}
	if err := store.Store(events); err != nil {
		t.Fatal(err)
	}

	view, err := storage.View()
	if err != nil {
		t.Fatal(err)
	}
	defer view.Close()

	shouldNotShowUp := Events{
		&PostDraftedEvent{"slug-3", "title", "body"},
	}

	if err := store.Store(shouldNotShowUp); err != nil {
		t.Fatal(err)
	}

	count := 0
	if err := view.ForEach(func(event Event) error {
		count++
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if want := len(events); count != want {
		t.Fatalf("len(events) = %d; want %d", count, want)
	}
}
