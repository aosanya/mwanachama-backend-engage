# Engage

## Objects

| Object | Holds |
|---|---|
| `post` | Its kind (`notice`, `event`, `volunteer_call`); title and body; the groups it goes to (opaque group ids); who posted it (from the caller); posted at; its state (`sent`, `cancelled`, and for a volunteer call `filled`). An event adds starts at, ends at, time zone and place; a volunteer call adds what, a start time with an optional end, and how many. |
| `rsvp` | An event, a member id from the caller, `going` or `not_going`, answered at. One per member per event; replaced until the event starts. |
| `attendance` | An event, a member id, checked in at, by whom (the member themselves or a holder), and whether it is unticked. |
| `offer` | A volunteer call, a member id from the caller, offered at. One per member per call; withdrawable while the call is open, until it is picked. |
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

Members offer while the call is `sent`, and withdraw their offer until it is
picked. The poster picks from the offers, no more than how many, and may
unpick. The poster may close the call with fewer picks than how many: with at
least one pick it becomes `filled`, with none it becomes `cancelled`; either
way no offer is taken after (decision 483).

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

## As built (EN2)

- **The caller's groups come from the host, through `UseGroups`.** The
  manager takes a `Groups` function, `func(ctx, callerID) ([]string,
  error)`, the same shape `mwanachama-backend-forms` takes for its waiting
  list. It is asked on every read that needs it and its answer is never
  stored. A manager given none reaches no group: the unread count is 0 and
  every notice answers a mark-read as a missing one.
- `POST /posts/{postID}/read` (`engage.read_marker.create`) records that the
  caller read a notice. The reader is bound `from: caller`; a body naming
  someone else is ignored. Only a `sent` notice is marked (409
  `ErrNotANotice` for an event, a volunteer call or a cancelled notice); a
  notice to none of the caller's groups answers 404 exactly as an unknown
  post; marking again returns the first marker, so `read_at` is when it was
  first read. No caller is 401 `ErrNoReader`.
- `GET /notices/unread` (`engage.notice.count_unread`) answers `{"count": n}`
  for the Home card: `sent` notices to any of the caller's groups with no
  marker of theirs. It reads every sent notice and filters by group in Go,
  because `group_ids` is a JSON column; a Postgres-side filter is a later
  optimisation, not a behaviour change.
- **A sent post is returned to the host.** `CreatePost` answers the stored
  post (id, kind, groups, poster, `posted_at`) and nothing else; the host
  hands it to comm for the thread in each picked group and the Inbox line.
  The module calls nothing.
- **Read-only once sent** holds because no address edits a post's title,
  body or groups; the event changes and cancellation EN3 adds touch time,
  place and state only.

## As built (EN3)

- `PUT /posts/{postID}/rsvp` (`engage.rsvp.answer`) records the caller's
  `going` or `not_going` on an event. The responder is bound `from: caller`
  and the event from the address; a body naming either is overwritten. One
  row per person per event: answering again replaces the answer and moves
  `answered_at`, keeping the row's id.
- Only a `sent` event takes an answer (409 `ErrNotAnOpenEvent` for a notice,
  a volunteer call or a cancelled event). An event to none of the caller's
  groups answers 404 exactly as an unknown post, through the same
  `UseGroups` reach EN2 built. From `starts_at` on, an answer is refused
  with 409 `ErrEventStarted`. No caller is 401 `ErrNoResponder`; an answer
  outside the declared values is 400 `ErrInvalidRSVP`.
- `PATCH /posts/{postID}/event` (`engage.event.change`) changes a `sent`
  event's `starts_at`, `ends_at`, `time_zone` or `place`; an empty field is
  kept, and the result must still end after it starts (400
  `ErrInvalidPost`). `POST /posts/{postID}/cancelled` (`engage.event.cancel`)
  sets the state to `cancelled`, after which RSVPs, check-ins and further
  changes are 409 `ErrNotAnOpenEvent`. Both bind the caller (401
  `ErrNoPoster` without one) and answer `Changed`: the event and the
  `participant_id` of everyone whose answer is `going`, for the host to tell
  them through comm. Decision 482 lets the poster or any holder of Post an
  event in that group do either; the module does not check who, because it
  knows no powers — that is the host's gate on the two actions.
- `POST /posts/{postID}/checkin` (`engage.attendance.check_in_self`) checks
  the caller in to a `sent` event to one of their groups, only from
  `starts_at` until `ends_at` (409 `ErrOutsideEvent` otherwise), with no
  proof. No caller is 401 `ErrNoAttendee`.
- `PUT /posts/{postID}/attendance/{participantID}`
  (`engage.attendance.check_in`) is the holder's check-in: anyone in the
  event's groups (asked of `UseGroups` for that person; anyone else is 404),
  on the day `starts_at` falls on in the event's `time_zone` (UTC when the
  zone does not load). `checked_in_by` is the caller.
  `DELETE /posts/{postID}/attendance/{participantID}`
  (`engage.attendance.untick`) sets `unticked` and keeps the row (404
  `ErrNotCheckedIn` when there is none); a later check-in by either address
  ticks the same row again and moves `checked_in_at`. That only a holder
  reaches these two is the host's gate.

## As built (EN5)

- `POST /posts/{postID}/offers` (`engage.offer.create`) records the caller's
  I can help on a `sent` volunteer call to one of their groups. The person is
  bound `from: caller`. One offer per person per call: offering again answers
  the first offer. A call outside the caller's groups is 404 exactly as an
  unknown post; any other kind of post, or a filled or cancelled call, is 409
  `ErrNotAnOpenCall`; no caller is 401 `ErrNoVolunteer`.
- `POST /posts/{postID}/picks` (`engage.pick.create`) picks one offer by
  `offer_id`. The picker is bound `from: caller` into `picked_by` and must be
  the call's `posted_by` (403 `ErrNotThePoster`); the call comes from the
  address. An offer made to another call is 404 `ErrOfferNotFound`. Once the
  call has `needed_count` picks a further pick is 409 `ErrCallFull`; picking
  the same offer again answers the first pick. The answer is `Picked`: the
  pick and the picked person's `participant_id`, for the host to raise comm's
  `volunteer_picked` (CM105).
- `POST /posts/{postID}/filled` (`engage.post.fill`) is the poster's confirm:
  the call's state becomes `filled`, after which offers and picks are 409
  `ErrNotAnOpenCall`. The answer is `Filled`: the call and every picked
  person's `participant_id`, for the host to raise `volunteer_call_filled`.
  The module calls nothing; the host tells comm.
- **Not built:** withdrawing an offer, unpicking, and listing a call's offers
  for the poster. The row did not name them; withdrawal is promised by the
  offer object's description ("Withdrawable until the call is filled").

## As built (EN6)

Decision 483 (owner, 2026-10-25).

- **Close with fewer picks.** `POST /posts/{postID}/filled` no longer needs
  `needed_count` picks: with one or more it sets `filled`, with none it sets
  `cancelled`. `Filled.Post.State` says which, so the host raises
  `volunteer_call_filled` on a `filled` close; on a `cancelled` one
  `ParticipantIDs` is empty and there is nobody picked to tell.
- **Withdraw.** `DELETE /posts/{postID}/offers/mine`
  (`engage.offer.withdraw`) removes the caller's own offer, the person bound
  `from: caller`, and answers the offer withdrawn. Only on a `sent` call
  (409 `ErrNotAnOpenCall` otherwise), only an offer not picked (409
  `ErrOfferPicked`; the poster unpicks first), 404 `ErrOfferNotFound` when
  the caller has no offer, 404 for a call outside their groups, 401
  `ErrNoVolunteer` with no caller. Offering again afterwards makes a new
  offer.
- **Unpick.** `DELETE /posts/{postID}/picks/{pickID}` (`engage.pick.delete`)
  is the poster's (`from: caller`, 403 `ErrNotThePoster`) on a `sent` call;
  the pick row goes and the offer stays. Answers `Picked` (the pick removed
  and who it was) for the host; an unknown pick or one on another call is 404
  `ErrPickNotFound`. A filled or cancelled call is not unpicked.
- **A start time with an optional end.** `needed_at` (free text) is gone from
  the blueprint and `Post`. A volunteer call's when is `starts_at`, required,
  and `ends_at`, optional and after the start when given (400
  `ErrInvalidPost` otherwise); they are the same two fields an event uses.
  No client sent `needed_at`, so nothing migrates.
- **Not built:** listing a call's offers for the poster, and any notice to a
  volunteer on unpick; the host decides whether comm says anything.
