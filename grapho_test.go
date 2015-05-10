package grapho

import (
	"testing"
	"time"
)

func Test_Grapho_DraftPost_DraftCanBeShown(t *testing.T) {
	app := NewGrapho(getenv("GRAPHO_ENV", "test"))
	now := time.Now()
	app.DraftPost(&DraftPostCommand{
		PostId: "test-id",
		Title:  "test-title",
		Body:   "test-body",
		Now:    now,
	})

	draft, err := app.ShowDraft("test-id")
	if err != nil {
		t.Fatal(err)
	}

	if draft == nil {
		t.Fatal("Draft not found")
	}
}

func Test_Grapho_DraftPost_DraftCanBeListed(t *testing.T) {
	app := NewGrapho(getenv("GRAPHO_ENV", "test"))
	now := time.Now()
	app.DraftPost(&DraftPostCommand{
		PostId: "test-id",
		Title:  "test-title",
		Body:   "test-body",
		Now:    now,
	})

	drafts, err := app.ListDrafts()
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(drafts), 1; got != want {
		t.Errorf("len(drafts) = %d; want %d", got, want)
	}
}
