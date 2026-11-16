package engage_test

import (
	"archive/zip"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	engage "github.com/aosanya/mwanachama-backend-engage"
)

var unknownZones = []string{"Mars/Olympus", "+02:00", "CEST", "../etc/passwd"}

var hostDependentZones = []string{"utc", "europe/budapest", "./Europe/Budapest", "Europe//Budapest", "Factory", "posixrules"}

func TestAZoneSpelledInAWayOnlySomeHostsResolveIsRefusedOnCreateAndChange(t *testing.T) {
	ctx := context.Background()
	pm := newTestManager(t)
	in := anEvent()
	in.TimeZone = "Europe/Budapest"
	ev, err := pm.CreatePost(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	for _, zone := range hostDependentZones {
		bad := anEvent()
		bad.TimeZone = zone
		if _, err := pm.CreatePost(ctx, bad); !errors.Is(err, engage.ErrInvalidPost) {
			t.Errorf("create with zone %q: err = %v, want ErrInvalidPost", zone, err)
		}
		if _, err := pm.ChangeEvent(ctx, engage.EventChange{PostID: ev.ID, ChangedBy: "actor-2", TimeZone: zone}); !errors.Is(err, engage.ErrInvalidPost) {
			t.Errorf("change to zone %q: err = %v, want ErrInvalidPost", zone, err)
		}
	}
	if got, _ := pm.GetPost(ctx, ev.ID); got.TimeZone != "Europe/Budapest" {
		t.Errorf("stored zone = %q, want Europe/Budapest", got.TimeZone)
	}
}

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
	for _, zone := range []string{"Europe/Budapest", "Africa/Nairobi", "UTC", "Etc/GMT+2", "Etc/GMT-14", "EST", "CET", "GMT", ""} {
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

func listedZoneNames(t *testing.T) []string {
	t.Helper()
	raw, err := os.ReadFile("zone_names.txt")
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(raw))
}

func TestEveryListedZoneNameLoads(t *testing.T) {
	names := listedZoneNames(t)
	if len(names) < 400 {
		t.Fatalf("zone_names.txt holds %d names, want the whole tzdata list", len(names))
	}
	for _, name := range names {
		if _, err := time.LoadLocation(name); err != nil {
			t.Errorf("listed zone %q does not load: %v", name, err)
		}
	}
}

func TestTheZoneNameListIsTheToolchainTzdataWithoutFactory(t *testing.T) {
	archive, err := zip.OpenReader(filepath.Join(runtime.GOROOT(), "lib", "time", "zoneinfo.zip"))
	if err != nil {
		t.Skipf("toolchain zoneinfo.zip not readable: %v", err)
	}
	defer archive.Close()
	var want []string
	for _, f := range archive.File {
		if f.Name != "Factory" {
			want = append(want, f.Name)
		}
	}
	got := listedZoneNames(t)
	sort.Strings(want)
	sort.Strings(got)
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("zone_names.txt differs from the toolchain tzdata: %d names listed, %d in the toolchain", len(got), len(want))
	}
}
