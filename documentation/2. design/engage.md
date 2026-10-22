# Engage

## Objects

| Object | Holds |
|---|---|
| `post` | Its kind (`notice`, `event`, `volunteer_call`); title and body; the groups it goes to (opaque group ids); who posted it (from the caller); posted at; its state (`sent`, `cancelled`, and for a volunteer call `filled`). An event adds starts at, ends at, time zone and place; a volunteer call adds what, when and how many. |
| `rsvp` | An event, a member id from the caller, `going` or `not_going`, answered at. One per member per event; replaced until the event starts. |
| `attendance` | An event, a member id, checked in at, by whom (the member themselves or a holder), and whether it is unticked. |
| `offer` | A volunteer call, a member id from the caller, offered at. One per member per call; withdrawable until the call is filled. |
| `pick` | A volunteer call, an offer, picked at, by whom. |
| `read_marker` | A notice, a member id from the caller, read at. Its absence is what makes a notice unread. |

A member and a group are ids the host passes in. The module never stores a
name, a role or a membership list; it is told on each call which groups the
caller is in.

A post is read-only once sent: its body is never edited in place. An event
may change its time or place, or be cancelled; each says so through comm.

## Unread count

A member's unread count is the number of `sent` notices to any of the
caller's groups with no read marker for that member. Events and volunteer
calls are not counted (decision 346).

## Check-in windows

| Who | When |
|---|---|
| The poster or a holder (host-gated) | On the event's day, in the event's time zone: from the start of the day it starts to the end of the day it ends. Checks in anyone in the event's groups, and unticks anyone. |
| A member, for themselves | Between starts at and ends at only, with no proof (decision 350). |

An untick keeps the attendance row and marks it unticked, so a later
check-in by either path ticks it again.

## Volunteer calls

Members offer while the call is `sent`. The poster picks from the offers, no
more than how many; confirming the picks sets the call `filled`, after which
no offer is taken.

## What comm does

This module calls nothing. Each change it makes is returned to the host,
which asks comm (CM104, CM105) to act:

| Happened here | comm does |
|---|---|
| A post sent | Opens a thread on it in each picked group's chat room; sends an Inbox line to the groups' members |
| An event changed or cancelled | Sends an Inbox line to the event's groups |
| Picks confirmed | Tells each picked member they were picked; tells the others who offered that the call is filled |

Responses to a post live only in comm's thread; this module stores none.

## As built (EN1)

- `engage.blueprint.json` declares the six objects above; the domain's own
  spec is `engage.wakala.json`, and `spec/examples/` carries a club and a
  school so the module is exercised under domains that share none of its
  host's nouns.
- The person on an `rsvp`, `attendance`, `offer` and `read_marker` is
  `participant_id`, because "member" is a host word this module may not use.
- A volunteer call's "what, when and how many" are `task`, `needed_at`
  (free text) and `needed_count`; an event's are `starts_at`, `ends_at`,
  `time_zone` and `place`.
- `engage.operations.json` declares two addresses, `POST /posts`
  (`engage.post.create`, poster bound from the caller) and
  `GET /posts/{postID}` (`engage.post.read`). `AnonymousActions` is empty.
  The rest of the addresses arrive with EN2–EN5.
- `domain_agnostic_test.go` skips an identifier reached through a selector,
  such as the engine's own `spec.Spec.ByRole`: "role" there names a blueprint
  slot, and a name this module declares is still caught where it is declared.
