package engage

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/aosanya/mwanachama-backend-shared/spec"
	"github.com/aosanya/mwanachama-backend-shared/specstore"

	"github.com/aosanya/mwanachama-backend-engage/models"
)

type (
	Post       = models.Post
	RSVP       = models.RSVP
	Attendance = models.Attendance
	Offer      = models.Offer
	Pick       = models.Pick
	ReadMarker = models.ReadMarker
)

type PostManager interface {
	CreatePost(ctx context.Context, in Post) (Post, error)
	GetPost(ctx context.Context, postID string) (Post, error)
}

type postManager struct {
	db  *gorm.DB
	st  *store
	now func() string
}

func NewPostManager(db *gorm.DB, s *spec.Spec) (PostManager, error) {
	if db == nil {
		return nil, fmt.Errorf("NewPostManager: db must not be nil")
	}
	st, err := newStore(db, s, map[string]any{
		objectPost:       models.Post{},
		objectRSVP:       models.RSVP{},
		objectAttendance: models.Attendance{},
		objectOffer:      models.Offer{},
		objectPick:       models.Pick{},
		objectReadMarker: models.ReadMarker{},
	})
	if err != nil {
		return nil, fmt.Errorf("NewPostManager: %w", err)
	}
	return &postManager{db: db, st: st, now: nowStamp}, nil
}

func nowStamp() string { return time.Now().UTC().Format(specstore.TimeLayout) }
