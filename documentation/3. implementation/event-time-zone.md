# Event time zone

An event or volunteer call stores `time_zone` as a name `time.LoadLocation` loads, such as `Europe/Budapest` or `UTC`.

## Rule

- `CreatePost` and `ChangeEvent` run `zoneLoads` inside `postHolds`. A name that does not load is refused with `ErrInvalidPost`, so a route answers 400 and nothing is stored. The message names the value and gives `Europe/Budapest` as the shape.
- Refused: `Mars/Olympus`, `+02:00`, `CEST`, `../etc/passwd`, and spellings only some hosts resolve: `utc`, `europe/budapest`, `./Europe/Budapest`, `Europe//Budapest`, `Factory`, `posixrules`. Allowed: a name spelled exactly as in `zone_names.txt`, for example `UTC`, `Etc/GMT+2`, `Europe/Budapest`, `EST`, `CET`, `GMT`.
- `time.LoadLocation` alone cannot decide this: it reads the host's zoneinfo directory before the embedded data, so a case-insensitive filesystem resolves `utc` and a Linux host does not, and `Location.String()` returns the name it was given, so it cannot tell the spellings apart. `zone_names.txt` holds the names in the toolchain's tzdata (`zoneinfo.zip`), one per line, without `Factory`; `zoneLoads` requires the name to be in it and to load, so the answer is the same on every host.
- `zone_names.txt` is rebuilt after a Go upgrade with `unzip -Z1 "$(go env GOROOT)/lib/time/zoneinfo.zip" | grep -vx Factory | LC_ALL=C sort > zone_names.txt`; `TestTheZoneNameListIsTheToolchainTzdataWithoutFactory` fails when it is behind.
- `Local` and the empty zone are not decided by the list: both are accepted as before.
- An empty zone is allowed on create and means UTC for the holder check-in day. An empty `time_zone` on `PATCH /posts/{postID}/event` keeps the stored zone.
- `time/tzdata` is imported in `zone.go`, so a host without a zoneinfo directory still accepts real names.
- `onTheDay` still falls back to UTC for a zone it cannot load. Only a row stored before this check can hold such a zone; no new row can.

## Proof

`zone_test.go` and `routes/zone_test.go`: refusals on create and change (the stored zone stays), real zones accepted, 400 on the routes, and the holder check-in day compared with the day in the stored zone for five fixed-offset zones, with a guard that at least one zone separates the local day from the UTC day.

## Hub copy

`wakala-api/modules/engage` needs the same change (`zone.go`, the call in `postHolds`, the tests) until ISS-105 replaces the copy with an import. Its skipped probe `TestSlotEventDates_EventTimeZoneMustBeAKnownZone` loses its `t.Skip` then.

## Open for the owner

`time.LoadLocation("Local")` loads, and means the server's own zone. It is accepted today.
