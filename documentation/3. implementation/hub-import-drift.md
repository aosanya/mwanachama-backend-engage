# Hub import drift

What differs between the hub's in-tree copy of the engage module (`projects/hub-steward/code/wakala-api/modules/engage`) and this repository, found on 2026-11-19 before the hub drops its copy and imports this module. The method is the one in [hub copy diff](hub-copy-diff.md): `scripts/hub-copy-diff.sh` rewrites the hub's import paths, strips comments and compares every `*.go` and `*.json` file; the two trees were also compared by file list and by the hub's commit history for the folder (a32f3cc8, ec77c144, b1f9348c, a37f9023, 20050ce0, 26b4f365, a4fab15f, 46ff9ba7).

The result before this change was one file only in the hub, four items only here, and no file that differs beyond the import path apart from `post.go`.

## Only in the hub

| Item | What it is | Disposition |
|---|---|---|
| `routes/build_cover_test.go` | Hub commit 46ff9ba7: tests that `Build`, `BuildFor` and `RoutesFor` declare the same addresses, that every split route is gated when nothing is anonymous, and that a refused action answers 403 and names an `engage.` action | Moved here with the import path rewritten; passes unchanged; no pinned count touched |

No other file, declared JSON file, address, type or exported function exists only in the hub. The files `engage.blueprint.json`, `engage.operations.json`, `engage.wakala.json`, `spec/examples/*.json` and `models/*.go` are identical after the import rewrite.

## Only in this repository

| Item | What it is | Disposition |
|---|---|---|
| `zone.go`, `zone_test.go`, `zone_names.txt`, `routes/zone_test.go` | The event time zone rule (see [event time zone](event-time-zone.md)) | Kept; the hub gains it when it imports this module |
| the `zoneLoads` call in `postHolds` (`post.go`) | Refuses a `time_zone` that is not spelled exactly as in `zone_names.txt`, on create and on change, for every post kind | Kept; the only difference in a shared file |
| `cmd/ddl`, `scripts/`, `documentation/`, `CLAUDE.md`, `Makefile`, `go.mod`, `LICENSE` | Repository tooling and documents the hub copy never carried | Kept |

## Differing beyond the import path

| Item | Hub | This repository | Disposition |
|---|---|---|---|
| `post.go`, `postHolds` | no zone check; `Mars/Olympus` and `+02:00` are stored | refuses them with `ErrInvalidPost` | This repository is the later behaviour and the hub has no rule that contradicts it, so nothing is merged into this side. Effect on the hub at import: see the next section |

No two behaviours conflict for the same operation in the sense that needs an owner decision; the hub side is the older of the two and carries a known-defect marker for exactly this gap.

## Hub-only wiring that is not part of the module

Not moved; listed so the hub's import step keeps it.

| Item | Where in the hub | Note |
|---|---|---|
| HTTP mount for the engage routes, the per-instance prefix `/instances/{instanceID}/engage` and the caller role gate | `internal/api/http/engage_routes.go` | uses `routes.Mount`, `routes.Table.Shape`, `routes.RoutesFor`, all present here |
| Post-kind permission split (`engage.post.create` into `create_notice`, `create_event`, `create_volunteer_call`) | `internal/api/http/engage_post_kind.go` | hub gate in front of the create address |
| Permission catalogue entries for the `engage.` actions and the control scan that groups them | `modules/agency/models/permissions_catalog.json`, `modules/agency/control_scan.go` | stays with the agency module |
| Mounting through `operation_routes.go`, the post manager cache and its building, group lookups for a post, hidden rows across modules, consent fan-out, the account router's `PostManager` | `internal/api/http/operation_routes.go`, `internal/agencies/registry.go`, `caches.go`, `hidden_rows.go`, `consent_withdrawal.go`, `test_copy.go`, `internal/accounts/router.go`, `internal/api/http/membership.go` | uses `SpecForMount`, `Provision`, `NewPostManager`, `UseGroups`, `UseGroupExists`, `Groups`, `GroupExists`, `Hidden`, `Operations`, all present here |
| Hub tests that import the module | `internal/api/http/misuse_removedmember_sweep_test.go`, `internal/agencies/consent_fanout_db_test.go`, the `misuse_engage_*` and `engage_authors_test.go` files | only the import path changes |

## What the hub's import step must do

- Change every `wakalaserver/modules/engage` and `wakalaserver/modules/engage/routes` import to `github.com/aosanya/mwanachama-backend-engage` and `.../routes`, then delete `modules/engage`.
- `internal/api/http/slot_event_dates_test.go` holds `TestSlotEventDates_EventTimeZoneMustBeAKnownZone` behind a `t.Skip` for the defect this repository already fixes; the skip is removed in the hub at import.
- Run the hub's tests on the memory backend after the import; a hub test that stores a zone this repository now refuses would fail and is a hub fixture to correct, not a reason to loosen the rule here.

## Result after the merge

`scripts/hub-copy-diff.sh` prints `diff post.go` (the `zoneLoads` call) and the `Only in backend` lines for `zone_test.go` and `routes/zone_test.go` (`zone.go` is listed as well); it prints no `Only in hub` line. Its exit code stays 1 only because of the `post.go` difference, which is the time zone rule the hub lacks.
