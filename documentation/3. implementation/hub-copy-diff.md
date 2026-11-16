# Hub copy diff

`mwanachama-backend-engage` is the custodian of the engage module. The hub (`projects/hub-steward/code/wakala-api/modules/engage`) ported it on 2026-11-07 (hub a32f3cc8) and changed it afterwards. Every change made in that copy exists here, so the two differ by the import path only (hub issue ISS-1 of the engage project, the engage part of ISS-105).

## Command

```
scripts/hub-copy-diff.sh [path-to-hub-modules-engage]
```

The default path is `../projects/hub-steward/code/wakala-api/modules/engage`, relative to this repository's parent folder. The exit code is 1 when the trees differ in any way and 0 when they are identical. The script and `scripts/hubdiff` are the ones of `mwanachama-backend-forms`, pointed at `modules/engage`.

## What it does

1. Copies the `*.go` and `*.json` files of both trees into two temporary folders, leaving out `cmd/`, `documentation/`, `scripts/` and `.git/`.
2. Rewrites the hub's imports: `wakalaserver/modules/shared/` becomes `github.com/aosanya/mwanachama-backend-shared/` and `wakalaserver/modules/engage` becomes this module's path.
3. Runs every Go file on both sides through `scripts/hubdiff`, which parses it, drops all comments and prints it with `gofmt`, so comments and layout do not show as differences.
4. Runs `diff -r backend hub`.

## Result of 2026-11-16

The output is empty: every file, tests and the declared JSON files included, is identical after the rewrite, and neither tree holds a file the other lacks.

## The hub commits carried over

One commit here per group of hub commits, each naming the hub hashes: a32f3cc8 (the event fixture starts thirty days ahead, so the event tests no longer fail once a fixed date passes); ec77c144, b1f9348c, a37f9023 (the consent gate and the withdrawal call); 20050ce0 (a withdrawn person's answers and unpicked offers are hidden, restored and listed); 26b4f365 (a post goes only to groups of its own agency); a4fab15f (the author of a post, pick and change comes from the caller).
