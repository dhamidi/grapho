package grapho

type Draft struct{}

func (self *Draft) fromEvent(event *PostDraftedEvent) *Draft {
	return &Draft{}
}

type AllDraftsView struct{}

func NewAllDraftsView() *AllDraftsView {
	return &AllDraftsView{}
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
	return &Draft{}, nil
}

func (self *AllDraftsView) putDraft(draft *Draft) {

}
