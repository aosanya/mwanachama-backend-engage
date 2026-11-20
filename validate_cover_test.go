package engage_test

import (
	"encoding/json"
	"strings"
	"testing"

	engage "github.com/aosanya/mwanachama-backend-engage"
)

func aSentNotice() engage.Post {
	return engage.Post{Kind: "notice", Title: "Water cut on Friday", State: "sent"}
}

func TestCheckAcceptsASentNoticeByValueAndByPointer(t *testing.T) {
	s, err := engage.SpecFor("test")
	if err != nil {
		t.Fatal(err)
	}
	p := aSentNotice()
	if err := engage.Check(s, engage.ObjectPost, p); err != nil {
		t.Fatalf("value: %v", err)
	}
	if err := engage.Check(s, engage.ObjectPost, &p); err != nil {
		t.Fatalf("pointer: %v", err)
	}
}

func TestCheckRefusesAPostWithABlankTitleOrKindOrState(t *testing.T) {
	s, err := engage.SpecFor("test")
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]func(*engage.Post){
		"title": func(p *engage.Post) { p.Title = "  \t" },
		"kind":  func(p *engage.Post) { p.Kind = "" },
		"state": func(p *engage.Post) { p.State = "" },
	}
	for field, blank := range cases {
		t.Run(field, func(t *testing.T) {
			p := aSentNotice()
			blank(&p)
			err := engage.Check(s, engage.ObjectPost, p)
			if err == nil || err.Error() != field+" is required" {
				t.Fatalf("want %q, got %v", field+" is required", err)
			}
		})
	}
}

func TestCheckRefusesAKindOrStateOutsideTheDeclaredValues(t *testing.T) {
	s, err := engage.SpecFor("test")
	if err != nil {
		t.Fatal(err)
	}
	p := aSentNotice()
	p.Kind = "poll"
	err = engage.Check(s, engage.ObjectPost, p)
	if err == nil || !strings.Contains(err.Error(), `kind is "poll"`) || !strings.Contains(err.Error(), "notice, event, volunteer_call") {
		t.Fatalf("kind refusal should name the value and the declared list, got %v", err)
	}

	p = aSentNotice()
	p.State = "archived"
	err = engage.Check(s, engage.ObjectPost, p)
	if err == nil || !strings.Contains(err.Error(), `state is "archived"`) || !strings.Contains(err.Error(), "sent, cancelled, filled") {
		t.Fatalf("state refusal should name the value and the declared list, got %v", err)
	}
}

func TestCheckAcceptsEveryDeclaredKindAndState(t *testing.T) {
	s, err := engage.SpecFor("test")
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"notice", "event", "volunteer_call"} {
		for _, state := range []string{"sent", "cancelled", "filled"} {
			p := engage.Post{Kind: kind, Title: "t", State: state}
			if err := engage.Check(s, engage.ObjectPost, p); err != nil {
				t.Fatalf("%s/%s: %v", kind, state, err)
			}
		}
	}
}

func TestCheckRefusesAnObjectTheDomainDoesNotFillAndANonStruct(t *testing.T) {
	s, err := engage.SpecFor("test")
	if err != nil {
		t.Fatal(err)
	}
	if err := engage.Check(s, "ballot", aSentNotice()); err == nil || !strings.Contains(err.Error(), `"ballot"`) {
		t.Fatalf("unknown object: got %v", err)
	}
	for _, v := range []any{"notice", 7, map[string]string{"title": "x"}, []engage.Post{aSentNotice()}} {
		if err := engage.Check(s, engage.ObjectPost, v); err == nil || !strings.Contains(err.Error(), "want a struct") {
			t.Fatalf("%T: want a struct refusal, got %v", v, err)
		}
	}
}

func TestOperationsIsTheDeclaredRouteTableAsJSON(t *testing.T) {
	var doc map[string]any
	if err := json.Unmarshal(engage.Operations(), &doc); err != nil {
		t.Fatalf("operations must be valid JSON: %v", err)
	}
	if len(doc) == 0 {
		t.Fatal("operations declared nothing")
	}
}

func TestSpecForMountGivesASecondBoardItsOwnTables(t *testing.T) {
	plain, err := engage.SpecFor("agy1")
	if err != nil {
		t.Fatal(err)
	}
	other, err := engage.SpecForMount("agy1", "second")
	if err != nil {
		t.Fatal(err)
	}
	if plain.Instance != "agy1" || other.Instance != "agy1" {
		t.Fatalf("instance not stamped: %q, %q", plain.Instance, other.Instance)
	}
	if other.Mount != "second" {
		t.Fatalf("mount not stamped: %q", other.Mount)
	}
	a, ok := plain.ByRole(engage.ObjectPost)
	b, ok2 := other.ByRole(engage.ObjectPost)
	if !ok || !ok2 {
		t.Fatal("post role missing")
	}
	if plain.TableFor(a) == other.TableFor(b) {
		t.Fatalf("two mounts of one instance must not share a table, both %q", plain.TableFor(a))
	}
}

func TestSpecForRefusesAnInstanceNameTheEngineRejects(t *testing.T) {
	if _, err := engage.SpecFor("has_underscore"); err == nil {
		t.Fatal("an instance carrying an underscore must be refused")
	}
}
