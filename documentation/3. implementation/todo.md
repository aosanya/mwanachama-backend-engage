# Task board — pending

| ID | Task | Status | Depends on |
|---|---|---|---|
| EN3 | 🧭 **Events: I'm going / Not going, and check-in** (member app decisions 346, 350): date, time, place and groups; an RSVP from a caller in one of the groups, changeable until it starts; a change of time or place, or a cancellation, returned to the host for comm; a holder checks in anyone in the event's groups on the event's day, a member checks themselves in only between starts at and ends at with no proof, a holder unticks anyone and a later check-in ticks again. **Built:** `PUT /posts/{postID}/rsvp` (engage.rsvp.answer), one answer per person per event, changeable until `starts_at`, from a caller in the event's groups (e9e0a74 merged as 87e33cb; engage.md As built EN3). **Left:** changing and cancelling an event, returned to the host for comm, waits on the owner (who may change or cancel: the poster, a holder of the event-posting power, or both); check-in and untick. Absorbs EN4. | ⏸️ | EN1 |
