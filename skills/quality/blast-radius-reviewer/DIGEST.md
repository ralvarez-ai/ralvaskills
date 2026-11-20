# Blast Radius Reviewer (digest)

Find what a change breaks elsewhere before it ships, and prove the one fact it is safe because of by running real code.

## Always-on rules
- Listing callers is not the job; find the breakage grep won't show: library source and pinned version, timing and teardown, wire formats, other languages, flags, code three hops away.
- Find the one fact the change is safe because of, then prove it. A writeup that sounds right is worthless.
- Rate each fact on the ladder and say where it stopped: claimed, cited (`file:line`), reasoned, ran (script or test against the real code), reproduced (running app).
- Each risk gets a real likelihood, a real cost and a `file:line`. List cleared risks separately. Never invent a caller or an API.
- Use CodeGraph for callers when `.codegraph/` exists.
- Hand back: what it does, the one safety fact (proven or marked unproven), risks, cleared, cheapest check before merge.

## Hard constraints
- Propose only; do not edit the code under review.
- Never present an unproven fact as proven.
- Strip anything private before the writeup goes public.

Invoke the full skill when: asked for the blast radius of a change, "what could this break", or reviewing a small diff you don't trust yet.

Full guide: invoke the blast-radius-reviewer skill.
