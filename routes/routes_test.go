package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	engage "github.com/aosanya/mwanachama-backend-engage"
	"github.com/aosanya/mwanachama-backend-engage/routes"
)

type callerKey struct{}

func callerFrom(ctx context.Context) string {
	id, _ := ctx.Value(callerKey{}).(string)
	return id
}

func newManager(t *testing.T) engage.PostManager {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	s, err := engage.SpecFor("routestest")
	if err != nil {
		t.Fatal(err)
	}
	if err := engage.Provision(db, s); err != nil {
		t.Fatal(err)
	}
	pm, err := engage.NewPostManager(db, s)
	if err != nil {
		t.Fatal(err)
	}
	return pm
}

func newServer(t *testing.T, caller string) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	for _, rt := range routes.RoutesFor(newManager(t), routes.Mount{Caller: callerFrom}) {
		mux.HandleFunc(rt.Pattern(""), rt.Handler)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), callerKey{}, caller)))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func post(t *testing.T, srv *httptest.Server, path string, body any) (int, map[string]any) {
	t.Helper()
	raw, _ := json.Marshal(body)
	res, err := http.Post(srv.URL+path, "application/json", bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var out map[string]any
	_ = json.NewDecoder(res.Body).Decode(&out)
	return res.StatusCode, out
}

func TestNoRouteIsAnonymous(t *testing.T) {
	split := routes.Split(newManager(t), routes.Mount{})
	if len(split.Anonymous) != 0 {
		t.Fatalf("%d routes answer without a caller; every address is gated", len(split.Anonymous))
	}
	if len(routes.AnonymousActions) != 0 {
		t.Fatalf("AnonymousActions = %v, want none", routes.AnonymousActions)
	}
}

func aNotice(postedBy string) map[string]any {
	return map[string]any{
		"kind":      "notice",
		"title":     "Water off on Friday",
		"group_ids": []string{"g-1"},
		"posted_by": postedBy,
		"state":     "cancelled",
	}
}

func TestThePosterComesFromTheCallerNotTheBody(t *testing.T) {
	code, p := post(t, newServer(t, "actor-1"), "/posts", aNotice("someone-else"))
	if code != http.StatusCreated || p["posted_by"] != "actor-1" || p["state"] != "sent" {
		t.Fatalf("create post = %d %v", code, p)
	}
}

func TestAPostWithNoCallerIsUnauthorized(t *testing.T) {
	code, _ := post(t, newServer(t, ""), "/posts", aNotice(""))
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", code)
	}
}

func TestAnUnknownPostAnswersNotFound(t *testing.T) {
	res, err := http.Get(newServer(t, "actor-1").URL + "/posts/nope")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", res.StatusCode)
	}
}
