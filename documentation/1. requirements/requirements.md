# Requirements

## Problem

An agency needs to tell its groups things, gather them, and ask for help:
a notice, a meeting or event, a call for volunteers. Chat carries talk but
no structure: nobody can see who is coming, check people in, or pick who
goes. The Engage card on the member app's Home needs something to open and
an unread notice count to show.

## Decisions

Settled by the owner in the member app's
[requirements.md](https://github.com/aosanya/mwanachama-frontend-member/blob/master/documentation/1.%20requirements/requirements.md),
2026-10-22.

| # | Rule |
|---|---|
| 346 | Engage holds notices, events and meetings, and volunteer calls; donations show as a placeholder until payments are taken up (decision 227). Engage opens with Notices, then Coming up, then Volunteer; the Engage card on Home shows how many notices are unread, and notices get no Home card of their own (decision 328 stands). |
| 347 | Engage lives in its own module, `mwanachama-backend-engage` (a public repository), built to the declared-domain standard and mounted per agency like voting (decision 341); comm opens each post's group thread and sends its Inbox lines. |
| 348 | Who may post each kind (a notice, an event or meeting, a volunteer call) is a power per kind, ticked per role (decision 263); owners and agents hold it when no role does (decision 270); a leader reaches the groups inside theirs (decision 266). Holders post from the studio or the member app. |
| 349 | A post goes to the groups picked and is read-only; members respond in a thread on it in the group's chat room, and each post is also an Inbox line. |
| 350 | An event or meeting has a date, time, place and groups; members say I'm going or Not going beforehand; on the day the poster or a holder checks people in, and members may check themselves in during the event's time with no proof; a holder may untick anyone. |
| 351 | A volunteer call says what, when and how many; members tap I can help; the poster picks who goes and they are told, and the rest are told it is filled. |

## What the module must do

- Hold a post of one kind (notice, event or volunteer call), sent to the
  groups picked, read-only once sent.
- Keep a read marker per member per notice and give a member their unread
  notice count across the groups they are in.
- For an event: hold its date, time, place and groups; take a member's I'm
  going or Not going before it starts, changeable until then.
- Check people in: the poster or a holder on the day, a member themselves
  during the event's time with no proof; a holder may untick anyone.
- For a volunteer call: hold what, when and how many; take a member's offer
  of help; let the poster pick who goes, no more than how many, and mark the
  call filled.
- Say what happened (a post sent, an event changed or cancelled, a member
  picked, a call filled) so the host can have comm open threads and send
  Inbox lines.

## Out of scope for now

Donations, until payments are taken up (decision 227); the module holds no
donation object. Proof of attendance (codes, location, photos).

## Identity

The host (`mwanachama-wakala-api`) resolves who is calling, which groups they
are in and whether they hold the power to post each kind (decision 348). This
module receives opaque member and group ids and never learns anything else
about a person.
