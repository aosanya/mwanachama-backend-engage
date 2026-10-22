# mwanachama-backend-engage

Runs an agency's posts to its groups: notices, events and meetings, and
volunteer calls. Each post goes to the groups picked and is read-only;
members respond in a thread on it in the group's chat room. An event takes
I'm going / Not going and check-ins on the day; a volunteer call takes offers
of help, and the poster picks who goes. Each member has an unread notice
count. Donations are a placeholder until payments are taken up.

A standalone Go library with no server of its own. `mwanachama-wakala-api`
compiles it in and mounts it per agency under `/instances/{id}/engage`; the
host does identity and the posting-power checks, comm opens each post's
group thread and sends its Inbox lines, and this module sees only opaque
member and group ids.

- Requirements: [documentation/1. requirements/requirements.md](documentation/1.%20requirements/requirements.md)
- Design: [documentation/2. design/engage.md](documentation/2.%20design/engage.md)
- Board: [documentation/3. implementation/todo.md](documentation/3.%20implementation/todo.md)
