package grapho

import (
	"bytes"
	"testing"
)

func Test_Post_Draft_succeeds(t *testing.T) {
	subject := NewPost()
	events, err := subject.Draft("slug", "hello world", "Post\nbody\n")
	expected := Events{
		&PostDraftedEvent{
			Id:    "slug",
			Title: "hello world",
			Body:  "Post\nbody\n",
		},
	}

	if err != nil {
		t.Fatal(err)
	}

	if len(events) != len(expected) {
		t.Fatalf("len(events) = %d; want = %d", len(events), len(expected))
	}

	if want, got := asJSON(expected, events); !bytes.Equal(want, got) {
		t.Errorf("Got:\n%s\nWanted:\n%s\n", got, want)
	}

}
