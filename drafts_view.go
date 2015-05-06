package grapho

import (
	"sort"
	"time"
)

type Draft struct {
	Id        string
	Title     string
	Body      string
	DraftedAt time.Time
}

func (self *Draft) fromEvent(event *PostDraftedEvent) *Draft {
	return &Draft{
		Id:        event.Id,
		Title:     event.Title,
		Body:      event.Body,
		DraftedAt: event.DraftedAt,
	}
}

type DraftsById []*Draft

func (self DraftsById) Len() int           { return len(self) }
func (self DraftsById) Less(i, j int) bool { return self[i].Id < self[j].Id }
func (self DraftsById) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }

type AllDraftsView struct {
	byId  map[string]*Draft
	index DraftsById
}

func NewAllDraftsView() *AllDraftsView {
	return &AllDraftsView{
		byId: map[string]*Draft{},
	}
}

func (self *AllDraftsView) HandleEvent(event Event) error {
	switch evt := event.(type) {
	case *PostDraftedEvent:
		draft := new(Draft).fromEvent(evt)
		self.putDraft(draft)
	}
	return nil
}

func (self *AllDraftsView) List() (DraftsById, error) {
	return self.index, nil
}

func (self *AllDraftsView) Show(id string) (*Draft, error) {
	draft, found := self.byId[id]
	if found {
		return draft, nil
	} else {
		return nil, ErrNotFound
	}
}

func (self *AllDraftsView) putDraft(draft *Draft) {
	self.byId[draft.Id] = draft
	self.index = append(self.index, draft)
	sort.Sort(self.index)
}
