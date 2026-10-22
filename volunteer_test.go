package engage_test

import (
	"context"
	"errors"
	"testing"

	engage "github.com/aosanya/mwanachama-backend-engage"
)

func aCall(needed int) engage.Post {
	return engage.Post{
		Kind:        "volunteer_call",
		Title:       "Stewards wanted",
		GroupIDs:    []string{"g-1"},
		Task:        "Hand out water",
		NeededAt:    "Saturday morning",
		NeededCount: needed,
		PostedBy:    "actor-1",
	}
}

func TestAVolunteerCallIsOfferedPickedAndFilled(t *testing.T) {
	ctx := context.Background()
	pm := newTestManager(t)
	pm.UseGroups(groupsFixed(map[string][]string{"p-1": {"g-1"}, "p-2": {"g-1"}, "p-3": {"g-1"}}))
	call, _ := pm.CreatePost(ctx, aCall(1))
	o1, err := pm.OfferHelp(ctx, call.ID, "p-1")
	if err != nil || o1.OfferedAt == "" {
		t.Fatalf("OfferHelp = %+v, %v", o1, err)
	}
	again, _ := pm.OfferHelp(ctx, call.ID, "p-1")
	if again.ID != o1.ID {
		t.Fatalf("a second offer = %+v, want the first", again)
	}
	o2, _ := pm.OfferHelp(ctx, call.ID, "p-2")
	picked, err := pm.PickOffer(ctx, engage.Pick{PostID: call.ID, OfferID: o1.ID, PickedBy: "actor-1"})
	if err != nil || picked.ParticipantID != "p-1" || picked.Pick.PickedAt == "" {
		t.Fatalf("PickOffer = %+v, %v", picked, err)
	}
	if _, err := pm.PickOffer(ctx, engage.Pick{PostID: call.ID, OfferID: o2.ID, PickedBy: "actor-1"}); !errors.Is(err, engage.ErrCallFull) {
		t.Fatalf("a pick past the count: err = %v", err)
	}
	filled, err := pm.FillCall(ctx, call.ID, "actor-1")
	if err != nil || filled.Post.State != "filled" || len(filled.ParticipantIDs) != 1 || filled.ParticipantIDs[0] != "p-1" {
		t.Fatalf("FillCall = %+v, %v", filled, err)
	}
	if got, _ := pm.GetPost(ctx, call.ID); got.State != "filled" {
		t.Fatalf("state = %q, want filled", got.State)
	}
	if _, err := pm.OfferHelp(ctx, call.ID, "p-3"); !errors.Is(err, engage.ErrNotAnOpenCall) {
		t.Fatalf("an offer to a filled call: err = %v", err)
	}
}

func TestVolunteerActsAreRefusedWhereTheyDoNotBelong(t *testing.T) {
	ctx := context.Background()
	pm := newTestManager(t)
	pm.UseGroups(groupsFixed(map[string][]string{"p-1": {"g-1"}}))
	call, _ := pm.CreatePost(ctx, aCall(2))
	other, _ := pm.CreatePost(ctx, aCall(2))
	ev, _ := pm.CreatePost(ctx, anEvent())
	offer, _ := pm.OfferHelp(ctx, call.ID, "p-1")
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"offer with no caller", second(pm.OfferHelp(ctx, call.ID, "")), engage.ErrNoVolunteer},
		{"offer outside the groups", second(pm.OfferHelp(ctx, call.ID, "p-2")), engage.ErrPostNotFound},
		{"offer to an event", second(pm.OfferHelp(ctx, ev.ID, "p-1")), engage.ErrNotAnOpenCall},
		{"pick by someone else", second(pm.PickOffer(ctx, engage.Pick{PostID: call.ID, OfferID: offer.ID, PickedBy: "p-1"})), engage.ErrNotThePoster},
		{"pick an offer to another call", second(pm.PickOffer(ctx, engage.Pick{PostID: other.ID, OfferID: offer.ID, PickedBy: "actor-1"})), engage.ErrOfferNotFound},
		{"fill by someone else", second(pm.FillCall(ctx, call.ID, "p-1")), engage.ErrNotThePoster},
		{"fill with no caller", second(pm.FillCall(ctx, call.ID, "")), engage.ErrNoPoster},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

func second[T any](_ T, err error) error { return err }
