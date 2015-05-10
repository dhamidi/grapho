package grapho

import (
	"log"
	"os"
	"testing"
)

func diskBackend() EventStorage {
	os.Remove("_test/events.log")
	onDisk, err := NewEventsOnDisk("_test/events.log")
	if err != nil {
		log.Fatal(err)
	}

	return onDisk
}

func memoryBackend() EventStorage {
	result, _ := NewEventsInMemory()
	return result
}

var backend = diskBackend

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
		&PostDraftedEvent{
			Id:    "slug-1",
			Title: "title",
			Body:  "body",
		},
		&PostDraftedEvent{
			Id:    "slug-2",
			Title: "title",
			Body:  "body",
		},
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
		&PostDraftedEvent{
			Id:    "slug-3",
			Title: "title",
			Body:  "body",
		},
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

func Test_Store_IsPersistent(t *testing.T) {
	storage := backend()
	store := NewEventStore(storage)
	events := Events{
		&PostDraftedEvent{
			Id:    "slug-1",
			Title: "title",
			Body:  "body",
		},
		&PostDraftedEvent{
			Id:    "slug-2",
			Title: "title",
			Body:  "body",
		},
	}
	if err := store.Store(events); err != nil {
		t.Fatal(err)
	}

	store.Close()

	onDisk, err := NewEventsOnDisk("_test/events.log")
	if err != nil {
		log.Fatal(err)
	}

	store = NewEventStore(onDisk)
	count := 0
	counter := func(Event) error { count++; return nil }

	if err := store.ReplayFunc(counter); err != nil {
		t.Fatal(err)
	}

	if want := len(events); count != want {
		t.Errorf("count = %d; want = %d", count, want)
	}
}

func TestMain(m *testing.M) {
	log.Printf("EventsOnDisk")
	backend = diskBackend
	failed := m.Run() != 0
	log.Printf("EventsInMemory")
	backend = memoryBackend
	failed = failed || m.Run() != 0

	if failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
