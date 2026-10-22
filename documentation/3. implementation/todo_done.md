# Task board — done

| ID | Task | Done | Notes |
|---|---|---|---|
| EN1 | 🧭 Scaffold the module in the declared-domain shape (decisions 346, 347) | 2026-10-22 | 7bc8d78 (merged as d01efb7): blueprint with post (notice, event, volunteer_call; sent, cancelled, filled), rsvp, attendance, offer, pick, read_marker, each per-person object unique on post and person (person stored as `participant_id`); engage.wakala.json plus club and school examples; specstore PostManager with CreatePost (poster from the caller) and GetPost; engage.operations.json (POST /posts, GET /posts/{postID}), AnonymousActions empty; Provision, cmd/ddl; the six guard tests (engage.md, As built EN1). Verified: go build, vet, test, gofmt clean. Open for the owner: a call's when as free text or a window; a pick's own participant id; a state between sent and filled. |
