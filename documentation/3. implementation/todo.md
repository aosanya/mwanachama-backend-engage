# Task board — pending

| ID | Task | Status | Depends on |
|---|---|---|---|
| EN3 | 🧭 **Events and I'm going / Not going** (member app decisions 346, 350): date, time, place and groups; an RSVP from a caller in one of the groups, changeable until it starts; a change of time or place, or a cancellation, returned to the host for comm. 20:55 run 331845: `PUT /posts/{postID}/rsvp` (engage.rsvp.answer), one answer per person per event, changeable until `starts_at`, from a caller in the event's groups (e9e0a74 merged as 87e33cb; engage.md As built EN3; go test green). Left: changing an event's time or place and cancelling it, returned to the host for comm — waits on the owner: who may change or cancel (the poster, a holder of the event-posting power, or both). | 📋 | EN1 |
| EN4 | 🧭 **Check-in by a holder, self check-in inside the event's time, untick** (member app decision 350): a holder checks in anyone in the event's groups on the event's day; a member checks themselves in only between starts at and ends at, with no proof; a holder unticks anyone, and a later check-in ticks again. | ⏸️ | EN3 |
