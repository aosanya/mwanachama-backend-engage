package engage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	engage "github.com/aosanya/mwanachama-backend-engage"
)

var unknownZones = []string{"Mars/Olympus", "+02:00", "CEST", "../etc/passwd"}

func TestAnEventWithAZoneThatIsNotAnIANANameIsRefused(t *testing.T) {
	pm := newTestManager(t)
	for _, zone := range unknownZones {
		in := anEvent()
		in.TimeZone = zone
		if _, err := pm.CreatePost(context.Background(), in); !errors.Is(err, engage.ErrInvalidPost) {
			t.Errorf("zone %q: err = %v, want ErrInvalidPost", zone, err)
		}
	}
}

func TestAnEventWithARealOrEmptyZoneIsAccepted(t *testing.T) {
	pm := newTestManager(t)
	for _, zone := range []string{"Europe/Budapest", "Africa/Nairobi", "UTC", ""} {
		in := anEvent()
		in.TimeZone = zone
		p, err := pm.CreatePost(context.Background(), in)
		if err != nil || p.TimeZone != zone {
			t.Errorf("zone %q: post = %+v, err = %v", zone, p, err)
		}
	}
}

func TestChangingAnEventToAZoneThatIsNotAnIANANameIsRefusedAndKeepsTheOldZone(t *testing.T) {
	ctx := context.Background()
	pm := newTestManager(t)
	in := anEvent()
	in.TimeZone = "Europe/Budapest"
	ev, err := pm.CreatePost(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	for _, zone := range unknownZones {
		if _, err := pm.ChangeEvent(ctx, engage.EventChange{PostID: ev.ID, ChangedBy: "actor-2", TimeZone: zone}); !errors.Is(err, engage.ErrInvalidPost) {
			t.Errorf("zone %q: err = %v, want ErrInvalidPost", zone, err)
		}
	}
	if got, _ := pm.GetPost(ctx, ev.ID); got.TimeZone != "Europe/Budapest" {
		t.Errorf("stored zone = %q, want Europe/Budapest", got.TimeZone)
	}
	if ch, err := pm.ChangeEvent(ctx, engage.EventChange{PostID: ev.ID, ChangedBy: "actor-2", TimeZone: "Africa/Nairobi"}); err != nil || ch.Post.TimeZone != "Africa/Nairobi" {
		t.Errorf("change to a real zone = %+v, %v", ch, err)
	}
	if ch, err := pm.ChangeEvent(ctx, engage.EventChange{PostID: ev.ID, ChangedBy: "actor-2", Place: "Hall"}); err != nil || ch.Post.TimeZone != "Africa/Nairobi" {
		t.Errorf("change with no zone = %+v, %v", ch, err)
	}
}

func TestAHolderCheckInDayIsTheDayInTheStoredZone(t *testing.T) {
	ctx := context.Background()
	pm := newTestManager(t)
	pm.UseGroups(groupsFixed(map[string][]string{"p-1": {"g-1"}}))
	now := time.Now().UTC()
	start := now.Add(-10 * time.Hour)
	differs := 0
	for _, zone := range []string{"Etc/GMT-12", "Etc/GMT-6", "Etc/GMT+6", "Etc/GMT+12", "UTC"} {
		loc, err := time.LoadLocation(zone)
		if err != nil {
			t.Fatal(err)
		}
		sameLocalDay := start.In(loc).Format(time.DateOnly) == now.In(loc).Format(time.DateOnly)
		sameUTCDay := start.Format(time.DateOnly) == now.Format(time.DateOnly)
		if sameLocalDay != sameUTCDay {
			differs++
		}
		in := anEvent()
		in.StartsAt = start.Format(time.RFC3339)
		in.EndsAt = now.Add(time.Hour).Format(time.RFC3339)
		in.TimeZone = zone
		ev, err := pm.CreatePost(ctx, in)
		if err != nil {
			t.Fatalf("zone %s: %v", zone, err)
		}
		_, err = pm.CheckIn(ctx, ev.ID, "p-1", "holder-1")
		if sameLocalDay && err != nil {
			t.Errorf("zone %s: same local day but check-in err = %v", zone, err)
		}
		if !sameLocalDay && !errors.Is(err, engage.ErrOutsideEvent) {
			t.Errorf("zone %s: other local day but check-in err = %v, want ErrOutsideEvent", zone, err)
		}
	}
	if differs == 0 {
		t.Fatal("no zone separates the local day from the UTC day, so the test proves nothing")
	}
}
