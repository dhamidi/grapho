package grapho

import (
	"fmt"
	"regexp"
	"strings"
)

type PostId string

var (
	ErrMalformedPostId = fmt.Errorf("PostId: malformed (allowed: A-Za-z0-9_-)")
	validPostId        = regexp.MustCompile(`^[-_A-Za-z0-9]+$`)
)

func NewPostId(id string) (PostId, error) {
	id = strings.TrimSpace(id)
	if validPostId.MatchString(id) {
		return PostId(id), nil
	} else {
		return PostId(""), ErrMalformedPostId
	}
}

func (self PostId) String() string { return string(self) }
