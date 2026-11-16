# Event time zone

An event or volunteer call stores `time_zone` as a name `time.LoadLocation` loads, such as `Europe/Budapest` or `UTC`.

## Rule

- `CreatePost` and `ChangeEvent` run `zoneLoads` inside `postHolds`. A name that does not load is refused with `ErrInvalidPost`, so a route answers 400 and nothing is stored. The message names the value and gives `Europe/Budapest` as the shape.
- Refused: `Mars/Olympus`, `+02:00`, `CEST`, `../etc/passwd`. Allowed: any IANA name, and `UTC`.
- An empty zone is allowed on create and means UTC for the holder check-in day. An empty `time_zone` on `PATCH /posts/{postID}/event` keeps the stored zone.
- `time/tzdata` is imported in `zone.go`, so a host without a zoneinfo directory still accepts real names.
- `onTheDay` still falls back to UTC for a zone it cannot load. Only a row stored before this check can hold such a zone; no new row can.

## Proof

`zone_test.go` and `routes/zone_test.go`: refusals on create and change (the stored zone stays), real zones accepted, 400 on the routes, and the holder check-in day compared with the day in the stored zone for five fixed-offset zones, with a guard that at least one zone separates the local day from the UTC day.

## Hub copy

`wakala-api/modules/engage` needs the same change (`zone.go`, the call in `postHolds`, the tests) until ISS-105 replaces the copy with an import. Its skipped probe `TestSlotEventDates_EventTimeZoneMustBeAKnownZone` loses its `t.Skip` then.

## Open for the owner

`time.LoadLocation("Local")` loads, and means the server's own zone. It is accepted today.
