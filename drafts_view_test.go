package grapho

import (
	"fmt"
	"testing"
	"time"
)

func Test_DraftsView_AddsDraft(t *testing.T) {
	event := &PostDraftedEvent{
		Id:    "slug",
		Title: "a draft",
		Body:  "for a post",
	}

	view := NewAllDraftsView()
	if err := view.HandleEvent(event); err != nil {
		t.Fatal(err)
	}

	draft, err := view.Show("slug")
	if err != nil {
		t.Fatal(err)
	}

	if draft == nil {
		t.Fatalf("Draft %q not found", "slug")
	}
}

func Test_DraftsView_ReturnsErrNotFound_WhenDraftDoesNotExist(t *testing.T) {
	view := NewAllDraftsView()
	_, err := view.Show("does-not-exist")
	if err != ErrNotFound {
		t.Fatalf("err = %#v; want %s", err, ErrNotFound)
	}
}

func Test_DraftsView_ListsAllDrafts_InAlphabeticalOrder_ById(t *testing.T) {
	ids := []string{"b", "c", "a"}
	view := NewAllDraftsView()
	for i, id := range ids {
		if err := view.HandleEvent(&PostDraftedEvent{
			Id:    id,
			Title: fmt.Sprintf("draft-%d", i),
			Body:  "body",
		}); err != nil {
			t.Error(err)
		}
	}

	drafts, err := view.List()
	if err != nil {
		t.Error(err)
	}
	for i, want := range []string{"a", "b", "c"} {
		if got := drafts[i].Id; got != want {
			t.Errorf("drafts[%d].Id = %q; want %q", i, got, want)
		}
	}
}

func Test_DraftsView_Draft_TracksDateDraftedAt(t *testing.T) {
	now := time.Now()
	event := &PostDraftedEvent{
		Id:        "slug",
		Title:     "A post",
		Body:      "yeah",
		DraftedAt: now,
	}
	view := NewAllDraftsView()
	if err := view.HandleEvent(event); err != nil {
		t.Error(err)
	}

	draft, err := view.Show("slug")
	if err != nil {
		t.Error(err)
	}
	if got := draft.DraftedAt; !got.Equal(now) {
		t.Errorf("draft.DraftedAt = %q; want %q", got, now)
	}
}

func Test_DraftsView_Draft_TracksDraftTitle(t *testing.T) {
	now := time.Now()
	event := &PostDraftedEvent{
		Id:        "slug",
		Title:     "A post",
		Body:      "yeah",
		DraftedAt: now,
	}
	view := NewAllDraftsView()
	if err := view.HandleEvent(event); err != nil {
		t.Error(err)
	}

	draft, err := view.Show("slug")
	if err != nil {
		t.Error(err)
	}
	if got, want := draft.Title, "A post"; got != want {
		t.Errorf("draft.DraftedAt = %q; want %q", got, want)
	}

}

func Test_DraftsView_Draft_TracksDraftBody(t *testing.T) {
	now := time.Now()
	event := &PostDraftedEvent{
		Id:        "slug",
		Title:     "A post",
		Body:      "yeah",
		DraftedAt: now,
	}
	view := NewAllDraftsView()
	if err := view.HandleEvent(event); err != nil {
		t.Error(err)
	}

	draft, err := view.Show("slug")
	if err != nil {
		t.Error(err)
	}
	if got, want := draft.Body, "yeah"; got != want {
		t.Errorf("draft.DraftedAt = %q; want %q", got, want)
	}

}
