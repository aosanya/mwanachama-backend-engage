# Hub copy diff

`mwanachama-backend-engage` is the custodian of the engage module. The hub (`projects/hub-steward/code/wakala-api/modules/engage`) ported it on 2026-11-07 (hub a32f3cc8) and changed it afterwards. Every change made in that copy exists here, so the two differ by the import path only (hub issue ISS-1 of the engage project, the engage part of ISS-105).

## Command

```
scripts/hub-copy-diff.sh [path-to-hub-modules-engage]
```

The default path is `../projects/hub-steward/code/wakala-api/modules/engage`, relative to this repository's parent folder. The exit code is 1 when the output has any line starting with `diff `, `Binary files ` or `Only in hub`, and 0 otherwise, so `Only in backend` lines alone (files only this repository holds) leave it at 0. The script and `scripts/hubdiff` are the ones of `mwanachama-backend-forms`, pointed at `modules/engage`; the script here carries the fix of hub issue I13 (the exit rule), and it already compared `*.json` files.

## What it does

1. Copies the `*.go` and `*.json` files of both trees into two temporary folders, leaving out `cmd/`, `documentation/`, `scripts/` and `.git/`.
2. Rewrites the hub's imports: `wakalaserver/modules/shared/` becomes `github.com/aosanya/mwanachama-backend-shared/` and `wakalaserver/modules/engage` becomes this module's path.
3. Runs every Go file on both sides through `scripts/hubdiff`, which parses it, drops all comments and prints it with `gofmt`, so comments and layout do not show as differences.
4. Runs `diff -r backend hub` and prints its output.
5. Exits 1 if that output has a `diff ` line (a file differs, JSON files included), a `Binary files ` line or an `Only in hub` line (the hub has something this repository lacks); otherwise exits 0.

## Result of 2026-11-16

After MUN696 the hub copy lacks the time zone rule, so the script prints a real difference today and exits 1: `diff post.go` (the `zoneLoads` check in `postHolds`, three lines here that the hub lacks), plus the `Only in backend` lines for `zone.go`, `zone_test.go` and `routes/zone_test.go`. The `Only in backend` lines alone would not fail it.

Until the hub copy gets the MUN696 change the output is not empty. Once it does, the output is the three `Only in backend` lines and the script exits 0. Every other file, tests and the declared JSON files included, is identical after the rewrite, and the hub holds no file this repository lacks.

## The hub commits carried over

One commit here per group of hub commits, each naming the hub hashes: a32f3cc8 (the event fixture starts thirty days ahead, so the event tests no longer fail once a fixed date passes); ec77c144, b1f9348c, a37f9023 (the consent gate and the withdrawal call); 20050ce0 (a withdrawn person's answers and unpicked offers are hidden, restored and listed); 26b4f365 (a post goes only to groups of its own agency); a4fab15f (the author of a post, pick and change comes from the caller).

## Result of 2026-11-19

The hub gained `routes/build_cover_test.go` (hub 46ff9ba7) after the last result; it now exists here too. The script prints `diff post.go` and the `Only in backend` lines for `zone.go`, `zone_test.go` and `routes/zone_test.go`, and no `Only in hub` line. The full account, with the hub wiring that stays in the hub, is in [hub import drift](hub-import-drift.md).
