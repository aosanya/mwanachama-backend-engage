package routes_test

import (
	"context"
	"net/http"
	"testing"

	engage "github.com/aosanya/mwanachama-backend-engage"
)

func anEventBody(zone string) map[string]any {
	return map[string]any{
		"kind":      "event",
		"title":     "Clean-up day",
		"group_ids": []string{"g-1"},
		"starts_at": "2030-05-01T10:00:00Z",
		"ends_at":   "2030-05-01T12:00:00Z",
		"time_zone": zone,
	}
}

func TestAnEventZoneThatIsNotAnIANANameAnswersBadRequest(t *testing.T) {
	srv := newServer(t, "actor-1")
	for _, zone := range []string{"Mars/Olympus", "+02:00"} {
		if code, body := post(t, srv, "/posts", anEventBody(zone)); code != http.StatusBadRequest {
			t.Errorf("zone %q: %d %v, want 400", zone, code, body)
		}
	}
	if code, body := post(t, srv, "/posts", anEventBody("Europe/Budapest")); code != http.StatusCreated {
		t.Errorf("Europe/Budapest: %d %v, want 201", code, body)
	}
}

func TestChangingAnEventZoneToOneThatIsNotAnIANANameAnswersBadRequest(t *testing.T) {
	pm := newManager(t)
	ev, err := pm.CreatePost(context.Background(), engage.Post{Kind: "event", Title: "Meet", GroupIDs: []string{"g-1"}, StartsAt: "2099-01-01T08:00:00Z", EndsAt: "2099-01-01T09:00:00Z", TimeZone: "Europe/Budapest", PostedBy: "poster"})
	if err != nil {
		t.Fatal(err)
	}
	for zone, want := range map[string]int{"Mars/Olympus": http.StatusBadRequest, "Europe/Berlin": http.StatusOK} {
		code, body := sendAs(t, pm, "poster", http.MethodPatch, "/posts/"+ev.ID+"/event", map[string]any{"time_zone": zone})
		if code != want {
			t.Errorf("zone %q: %d %v, want %d", zone, code, body, want)
		}
	}
}
