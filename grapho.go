package grapho

import "log"

type Grapho struct {
	config        *Config
	store         *EventStore
	listeners     []EventHandler
	allDraftsView *AllDraftsView
}

func NewGrapho(env string) *Grapho {
	result := &Grapho{
		config: Configurations[env],
	}
	result.setupStore()
	result.setupListeners()
	return result
}

func (self *Grapho) setupStore() {
	self.store = NewEventStore(self.config.Storage())
}

func (self *Grapho) setupListeners() {
	self.allDraftsView = NewAllDraftsView()
	self.listeners = []EventHandler{
		self.allDraftsView,
	}
}

func (self *Grapho) DraftPost(cmd *DraftPostCommand) {
	post := NewPost()
	events, _ := post.Draft(cmd)
	self.handleEvents(events)
}

func (self *Grapho) ShowDraft(draftId string) (*Draft, error) {
	return self.allDraftsView.Show(draftId)
}

func (self *Grapho) ListDrafts() ([]*Draft, error) {
	return self.allDraftsView.List()
}

func (self *Grapho) handleEvents(events Events) {
	err := self.store.Store(events)
	if err != nil {
		log.Printf("grapho.handleEvents: %s", err)
		return
	}

	self.notifyListeners(events)
}

func (self *Grapho) notifyListeners(events Events) {
	for _, listener := range self.listeners {
		events.ApplyTo(listener)
	}
}
