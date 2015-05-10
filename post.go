package grapho

import "time"

type Post struct{}

func NewPost() *Post { return &Post{} }

type PostDraftedEvent struct {
	Id        string
	Title     string
	Body      string
	DraftedAt time.Time
}

func (self *PostDraftedEvent) EventType() string { return "post/drafted" }
func init()                                      { RegisterEvent(&PostDraftedEvent{}) }

type DraftPostCommand struct {
	PostId PostId
	Title  string
	Body   string
	Now    time.Time
}

func (self *Post) Draft(params *DraftPostCommand) (Events, error) {
	return Events{
		&PostDraftedEvent{
			Id:        params.PostId.String(),
			Title:     params.Title,
			Body:      params.Body,
			DraftedAt: params.Now,
		},
	}, nil
}
