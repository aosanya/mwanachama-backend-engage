package engage

import "errors"

var (
	ErrPostNotFound   = errors.New("post not found")
	ErrInvalidPost    = errors.New("invalid post")
	ErrNoPoster       = errors.New("a post needs a caller to stand behind it")
	ErrNoReader       = errors.New("a read marker needs a caller")
	ErrNotANotice     = errors.New("only a sent notice is marked read")
	ErrNoResponder    = errors.New("an rsvp needs a caller")
	ErrInvalidRSVP    = errors.New("invalid rsvp")
	ErrNotAnOpenEvent = errors.New("only a sent event takes an rsvp")
	ErrEventStarted   = errors.New("the event has started")
)
