package grapho

import "testing"

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
