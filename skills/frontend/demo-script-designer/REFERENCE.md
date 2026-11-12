# Demo Script Reference

Tables moved out of [SKILL.md](SKILL.md).

## Demo progression template (§3)

| # | Focus | New capability | Approach |
|---|---|---|---|
| 1 | Basics | Single element access | Explicit reference |
| 2 | Integration | Multiple related elements | Grouped query |
| 3 | Abstraction | Natural language, no explicit refs | System infers intent |
| 4 | Dynamics | Live data, change over time | Continuous monitoring |
| 5 | Accessibility | Different access methods | Voice/touch/alternative interface |
| 6 | Automation | System behavior without command | Observe automatic logic |
| 7 | Diagnostics | Troubleshooting scenario | Bundle related elements |
| 8 | Knowledge | Search/lookup capability | Find answers in resources |
| 9 | Control | Issue commands | Write operations + verification |

## Anti-patterns (§4)

| Anti-pattern | Wrong | Right |
|---|---|---|
| Explanation/query mismatch | Query mentions `YI-100, SI-100, PIT-100`; explanation says "we didn't mention specific tags this time" | "Three different data points from two systems in parallel" |
| Pointing to invisible | "[Point to field `PIT-100` on screen]" when not labeled | "[Point to the pressure gauge display]" |
| Robotic narration | "System processes query, retrieves values from database, returns formatted response." | "Watch the system grab the values and hand them back as one answer." |
| Redundant demos | Demo 1 status query → Demo 2 same status query | Demo 1 status query → Demo 2 combine status with related data |
| Weak bridge | "Now let's do a different demo." | "Now that we know the status, let's check the pressure to see what it's delivering." |
