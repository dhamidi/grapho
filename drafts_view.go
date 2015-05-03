package grapho

type Draft struct {
	Id string
}

func (self *Draft) fromEvent(event *PostDraftedEvent) *Draft {
	return &Draft{
		Id: event.Id,
	}
}

type AllDraftsView struct {
	byId map[string]*Draft
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
}
