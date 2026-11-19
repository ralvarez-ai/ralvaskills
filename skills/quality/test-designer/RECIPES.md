# Test Designer — Recipes

Illustrative skeletons. Adapt names to the stack; the shape is what matters.

## 1. Replace a sleep with a condition poll

Wrong: the test passes only if 2 s is enough.

```python
start_worker()
time.sleep(2)
assert queue.empty()
```

Right: wait on the observable state, bounded by one ceiling.

```python
READY_TIMEOUT = 10.0

def wait_until(condition, timeout=READY_TIMEOUT, interval=0.05):
    deadline = time.monotonic() + timeout
    while time.monotonic() < deadline:
        if condition():
            return
        time.sleep(interval)
    raise TimeoutError("condition not met")

start_worker()
wait_until(queue.empty)
```

Better still: have the worker expose a `done` event or awaitable and await it; poll only when no signal exists.

## 2. Replace a sleep with a fake clock

Production code takes a clock; the test advances it.

```go
type Clock interface{ Now() time.Time }

type FakeClock struct{ t time.Time }
func (f *FakeClock) Now() time.Time          { return f.t }
func (f *FakeClock) Advance(d time.Duration) { f.t = f.t.Add(d) }

func TestSessionExpires(t *testing.T) {
	clk := &FakeClock{t: time.Unix(0, 0)}
	s := NewSession(clk, 30*time.Minute)
	clk.Advance(31 * time.Minute)
	if !s.Expired() { t.Fatal("want expired") }
}
```

No wall-clock time passes; the test is instant and exact at the boundary.

## 3. Test retry/backoff as a sequence of delays

```python
def test_backoff_delays():
    clock = FakeClock()
    call_with_retry(always_fails, clock=clock, attempts=4)
    assert clock.sleeps == [1, 2, 4]
```

Assert what was requested, not how long it took.

## 4. Readiness gate for e2e

```bash
# in the e2e job, before any test
until curl -fsS http://localhost:8080/healthz >/dev/null; do sleep 0.2; done
```

The loop is bounded by the CI job timeout and exits the moment the service is ready. In browser tests use the framework's auto-waiting assertions instead of `waitForTimeout`.

## 5. Layout and commands

```text
tests/
  logic/          # in-process, fakes, no network or real clock
  integration/    # one real dependency, readiness-gated
  e2e/            # full journeys against a running system
```

```yaml
# Taskfile.yml (names follow module:verb)
test:logic:        { cmds: ["<runner> tests/logic"] }
test:integration:  { cmds: ["<runner> tests/integration"] }
test:e2e:          { cmds: ["<runner> tests/e2e"] }
```

CI runs `test:logic`, then `test:integration`, then `test:e2e`, each its own job.

## 6. Finding existing waits

Start a cleanup by listing offenders, then classify each with the section 9 triage:

```bash
grep -rnE "sleep\(|time\.Sleep|setTimeout\(|waitForTimeout|Thread\.sleep" tests/
```

Every hit becomes an awaited signal, a fake clock, or a deadline-bounded poll.
