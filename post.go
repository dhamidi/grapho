package grapho

type Post struct{}

func NewPost() *Post { return &Post{} }

type PostDraftedEvent struct {
	Id    string
	Title string
	Body  string
}

func (self *PostDraftedEvent) EventType() string { return "post/drafted" }
func init()                                      { RegisterEvent(&PostDraftedEvent{}) }

func (self *Post) Draft(slug, title, body string) (Events, error) {
	return Events{
		&PostDraftedEvent{
			Id:    slug,
			Title: title,
			Body:  body,
		},
	}, nil
}
