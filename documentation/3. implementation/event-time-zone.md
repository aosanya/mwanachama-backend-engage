# Event time zone

A post stores `time_zone` as a time zone name spelled exactly as in `zone_names.txt`, such as `Europe/Budapest` or `UTC`. The holder check-in day of an event or volunteer call is read in that zone; the check itself covers every post kind.

## Rule

- `CreatePost` and `ChangeEvent` run `zoneLoads` inside `postHolds`, before the kind is looked at, so a notice with a bad `time_zone` is refused as well as an event or a volunteer call. A name that is not accepted is refused with `ErrInvalidPost`, so a route answers 400 and nothing is stored. The message names the value and gives `Europe/Budapest` as the shape.
- Refused: `Mars/Olympus`, `+02:00`, `CEST`, `../etc/passwd`, and spellings only some hosts resolve: `utc`, `europe/budapest`, `./Europe/Budapest`, `Europe//Budapest`, `Factory`, `posixrules`. Allowed: a name spelled exactly as in `zone_names.txt`, for example `UTC`, `Etc/GMT+2`, `Europe/Budapest`, `EST`, `CET`, `GMT`.
- `time.LoadLocation` alone cannot decide this: it reads the host's zoneinfo directory before the embedded data, so a case-insensitive filesystem resolves `utc` and a Linux host does not, and `Location.String()` returns the name it was given, so it cannot tell the spellings apart. `zone_names.txt` holds the names in the toolchain's tzdata (`zoneinfo.zip`), one per line, without `Factory`; `zoneLoads` requires the name to be in it and to load, so the answer is the same on every host.
- `zone_names.txt` is rebuilt after a Go upgrade with `unzip -Z1 "$(go env GOROOT)/lib/time/zoneinfo.zip" | grep -vx Factory | LC_ALL=C sort > zone_names.txt`; `TestTheZoneNameListIsTheToolchainTzdataWithoutFactory` fails when it is behind.
- `Local` and the empty zone are not decided by the list: both are accepted as before.
- An empty zone is allowed on create and means UTC for the holder check-in day. An empty `time_zone` on `PATCH /posts/{postID}/event` keeps the stored zone.
- `time/tzdata` is imported in `zone.go`, so a host without a zoneinfo directory still accepts real names.
- `onTheDay` still falls back to UTC for a zone it cannot load. Only a row stored before this check can hold such a zone; no new row can.

## Proof

`zone_test.go` and `routes/zone_test.go`: refusals on create and change (the stored zone stays), real zones accepted, 400 on the routes, and the holder check-in day compared with the day in the stored zone for five fixed-offset zones, with a guard that at least one zone separates the local day from the UTC day.

Red first: of the tests added with the MUN696 check, four failed before it and passed after (`TestAnEventWithAZoneThatIsNotAnIANANameIsRefused` and `TestChangingAnEventToAZoneThatIsNotAnIANANameIsRefusedAndKeepsTheOldZone` in `zone_test.go`, `TestAnEventZoneThatIsNotAnIANANameAnswersBadRequest` and `TestChangingAnEventZoneToOneThatIsNotAnIANANameAnswersBadRequest` in `routes/zone_test.go`). `TestAnEventWithARealOrEmptyZoneIsAccepted` and `TestAHolderCheckInDayIsTheDayInTheStoredZone` passed before the check too; they guard what must keep working. `TestAZoneSpelledInAWayOnlySomeHostsResolveIsRefusedOnCreateAndChange` was red first on a case-insensitive host (macOS) for all six spellings, and passes on any host once `zone_names.txt` decides. `TestEveryListedZoneNameLoads` and `TestTheZoneNameListIsTheToolchainTzdataWithoutFactory` guard the list.

## Hub copy

`wakala-api/modules/engage` needs the same change (`zone.go`, the call in `postHolds`, the tests) until ISS-105 replaces the copy with an import. Its skipped probe `TestSlotEventDates_EventTimeZoneMustBeAKnownZone` loses its `t.Skip` then.

## Open for the owner

`time.LoadLocation("Local")` loads, and means the server's own zone. It is accepted today.
