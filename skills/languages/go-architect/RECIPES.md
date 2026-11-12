# Go Architecture Recipes

Reference implementations from [SKILL.md](SKILL.md). Currently: a `go.work` multi-module workspace with two bounded contexts — `order` (a normal context) and `checkout` (an orchestrator that coordinates `order`, inventory, and payment), see §16 for the rules; and an external program behind a port (§4), at the end of this file.

## Multi-module workspace layout

```
myworkspace/
├── go.work
│
├── domain/                              # ONLY for value objects genuinely shared across contexts
│   └── money/
│       └── money.go
│
├── pkg/                                  # cross-cutting, shared by 2+ services — each its own module
│   ├── problem/                          # RFC 9457 mapper
│   ├── validation/
│   └── obs/
│
├── internal/
│   ├── order/                            # bounded context — own module
│   │   ├── domain/
│   │   │   └── order/
│   │   │       ├── order.go
│   │   │       ├── repository.go         # secondary port — lives in domain (DIP)
│   │   │       └── errors.go
│   │   ├── app/
│   │   │   └── order/
│   │   │       └── place_order.go
│   │   └── adapters/
│   │       ├── primary/
│   │       │   └── http/
│   │       │       └── handler.go
│   │       └── secondary/
│   │           └── postgres/
│   │               └── order_repo.go
│   │
│   └── checkout/                          # orchestrator context — own module
│       ├── app/
│       │   ├── ports.go                   # its OWN ports — not order's
│       │   └── checkout.go
│       └── adapters/
│           ├── primary/
│           │   └── http/
│           │       └── handler.go
│           └── secondary/
│               ├── orderclient/            # implements app.OrderCreator via HTTP
│               ├── inventoryclient/
│               └── paymentclient/
│
├── app/                                    # composition roots — one per deployable
│   ├── order/wire.go
│   └── checkout/wire.go
│
└── cmd/
    ├── order/main.go                        # thin entrypoint
    └── checkout/main.go
```

Every arrow points inward: `domain` imports nothing; `app` imports only `domain`; `adapters/*` import `app`+`domain` and satisfy their interfaces with compile-time assertions; `checkout` never imports `order` — only HTTP, decoded through the same `pkg/problem` mapper both services share.

## `domain/money/money.go` — shared value object (deliberate Shared Kernel)

```go
package money

type Money struct {
    Amount   int64 // minor units (cents)
    Currency string
}

func (m Money) IsZero() bool { return m.Amount == 0 }
```

## `internal/order/domain/order/` — pure business logic + the port

```go
// order.go
package order

func New(customer CustomerID, total money.Money) (Order, error) {
    if total.IsZero() {
        return Order{}, ErrEmptyOrder
    }
    return Order{ID: newID(), Customer: customer, Total: total, Status: StatusPending}, nil
}

type Order struct {
    ID       ID
    Customer CustomerID
    Total    money.Money
    Status   Status
}
```

```go
// repository.go — secondary port; DIP requires this live in domain, not app
package order

import "context"

type Repository interface {
    Save(ctx context.Context, o Order) error
    FindByID(ctx context.Context, id ID) (Order, error)
}
```

```go
// errors.go
package order

import "errors"

var (
    ErrEmptyOrder = errors.New("order: total cannot be zero")
    ErrNotFound   = errors.New("order: not found")
)
```

## `internal/order/app/order/place_order.go` — thin orchestration

No explicit primary port here: `PlaceOrderService` has exactly one production caller (the HTTP handler below), so per [hexagonal-architecture §7](../../design/hexagonal-architecture/SKILL.md#7-when-not-to-use-it) an interface would be ceremony without payoff. Introduce one (`PlaceOrderer`, declared here since this is the layer offering the capability) only when a second caller genuinely needs to swap it — e.g. a handler test injecting a fake.

```go
package order

type PlaceOrderCommand struct {
    Customer domain.CustomerID
    Total    money.Money
}

type PlaceOrderService struct {
    repo domain.Repository // depends on the port, not a concrete adapter
}

func NewPlaceOrderService(repo domain.Repository) *PlaceOrderService {
    return &PlaceOrderService{repo: repo}
}

func (s *PlaceOrderService) PlaceOrder(ctx context.Context, cmd PlaceOrderCommand) (domain.ID, error) {
    o, err := domain.New(cmd.Customer, cmd.Total)
    if err != nil {
        return domain.ID{}, fmt.Errorf("place order: %w", err)
    }
    if err := s.repo.Save(ctx, o); err != nil {
        return domain.ID{}, fmt.Errorf("place order: %w", err)
    }
    return o.ID, nil
}
```

## `internal/order/adapters/secondary/postgres/order_repo.go` — driven adapter

Default: `sqlc`-generated queries (`db.Queries`, built from `query.sql` + the migrated schema — see §12). The repo adapter is a thin translation layer between `db.Order` (generated row struct) and `domain.Order` (the aggregate).

```go
package postgres

var _ domain.Repository = (*OrderRepo)(nil) // makes the contract checkable at compile time

type OrderRepo struct{ queries *db.Queries }

func (r *OrderRepo) Save(ctx context.Context, o domain.Order) error {
    _, err := r.queries.SaveOrder(ctx, db.SaveOrderParams{
        ID: o.ID, Customer: o.Customer, Total: o.Total.Amount, Status: o.Status,
    })
    return err
}

func (r *OrderRepo) FindByID(ctx context.Context, id domain.ID) (domain.Order, error) {
    row, err := r.queries.GetOrderByID(ctx, id)
    if err != nil {
        return domain.Order{}, err
    }
    return domain.Order{ID: row.ID, Customer: row.Customer, Total: domain.NewMoney(row.Total), Status: row.Status}, nil
}
```

Alternative (`sqlx` + `//go:embed`, for query shapes sqlc can't express as a static named query):

```go
type OrderRepo struct{ db *sqlx.DB }

//go:embed queries/save_order.sql
var saveOrderSQL string

//go:embed queries/find_order_by_id.sql
var findOrderByIDSQL string

func (r *OrderRepo) Save(ctx context.Context, o domain.Order) error {
    _, err := r.db.ExecContext(ctx, saveOrderSQL, o.ID, o.Customer, o.Total.Amount, o.Status)
    return err
}

func (r *OrderRepo) FindByID(ctx context.Context, id domain.ID) (domain.Order, error) {
    var o domain.Order
    err := r.db.GetContext(ctx, &o, findOrderByIDSQL, id)
    return o, err
}
```

## `internal/order/adapters/primary/http/handler.go` — driving adapter

RFC 9457 is applied exactly at this boundary — the domain and app layers never format a wire response.

```go
package http

type Handler struct{ svc *orderapp.PlaceOrderService }

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
    var cmd orderapp.PlaceOrderCommand
    if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
        problem.Write(w, problem.BadRequest("invalid request body"))
        return
    }
    id, err := h.svc.PlaceOrder(r.Context(), cmd)
    if err != nil {
        problem.WriteFromError(w, err) // maps domain/app errors → RFC 9457
        return
    }
    json.NewEncoder(w).Encode(struct {
        ID string `json:"id"`
    }{id.String()})
}
```

## `internal/checkout/app/ports.go` — the orchestrator's own ports

Unlike `order`'s optional primary port, these are load-bearing from day one: `checkout` has no other way to reach `order`/inventory/payment, and each port already has two real implementations (the HTTP client below, and a fake used to unit-test the orchestration logic without booting three services).

```go
package app

type OrderCreator interface {
    CreateOrder(ctx context.Context, customer string, total money.Money) (string, error)
}
type StockReserver interface {
    Reserve(ctx context.Context, sku string, qty int) error
}
type PaymentCharger interface {
    Charge(ctx context.Context, customer string, amount money.Money) error
}
```

```go
// checkout.go
package app

type CheckoutService struct {
    orders    OrderCreator
    inventory StockReserver
    payments  PaymentCharger
}

func NewCheckoutService(orders OrderCreator, inventory StockReserver, payments PaymentCharger) *CheckoutService {
    return &CheckoutService{orders: orders, inventory: inventory, payments: payments}
}

func (s *CheckoutService) Checkout(ctx context.Context, cmd CheckoutCommand) error {
    if err := s.inventory.Reserve(ctx, cmd.SKU, cmd.Qty); err != nil {
        return fmt.Errorf("checkout: reserve stock: %w", err)
    }
    if err := s.payments.Charge(ctx, cmd.Customer, cmd.Total); err != nil {
        return fmt.Errorf("checkout: charge payment: %w", err) // caller compensates: release reservation
    }
    if _, err := s.orders.CreateOrder(ctx, cmd.Customer, cmd.Total); err != nil {
        return fmt.Errorf("checkout: create order: %w", err)
    }
    return nil
}
```

## `internal/checkout/adapters/secondary/orderclient/client.go` — crosses the real service boundary

An anticorruption layer per [ddd-architect §8](../../design/ddd-architect/SKILL.md#8-anticorruption-layers-acl): translates at the boundary, never lets `order`'s wire model leak into `checkout`.

```go
package orderclient

var _ app.OrderCreator = (*Client)(nil)

type Client struct {
    base *url.URL
    hc   *http.Client
}

func (c *Client) CreateOrder(ctx context.Context, customer string, total money.Money) (string, error) {
    resp, err := c.hc.Post(c.base.String()+"/orders", "application/json", encode(customer, total))
    if err != nil {
        return "", fmt.Errorf("orderclient: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        return "", problem.DecodeAsError(resp.Body) // decodes the RFC 9457 body back into a Go error
    }
    return decodeID(resp.Body)
}
```

## `app/order/wire.go` — composition root

```go
package app

func Run(ctx context.Context, cfg Config) error {
    db := mustConnect(cfg.DB)
    repo := &postgres.OrderRepo{DB: db}
    svc := orderapp.NewPlaceOrderService(repo)
    handler := &orderhttp.Handler{Svc: svc}
    return startServer(ctx, cfg.Addr, handler)
}
```

## `cmd/order/main.go` — thin entrypoint

```go
func main() {
    cfg := loadConfig()
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()
    if err := app.Run(ctx, cfg); err != nil {
        slog.Error("order service exited", "err", err)
        os.Exit(1)
    }
}
```

## An external program behind a port

For a tool you do not control (a CLI run as a subprocess, a vendor SDK) that runs batches of work
and writes results somewhere. The failure mode this prevents: the caller learns what happened by
scanning the tool's output directory, re-derives events from files, and every vendor flag, file
name and identifier convention leaks into the rest of the code. Rules: §4.

```
internal/execution/            # the domain's view: request, event, outcome, ports. Imports no vendor.
internal/execution/vendorcli/  # adapter named by technology: owns flags, files, parsing
internal/execution/fake/       # in-process adapter, scripted; the second adapter
internal/execution/executiontest/  # one conformance suite, run against every adapter
```

### `internal/execution/execution.go` — the port and its types

```go
package execution

import (
    "context"
    "errors"
)

type (
    // Kind names an event; the payload is the Event field that matches it.
    Kind int

    // Request is what to run, in the domain's words. It never holds the tool's flags.
    Request struct {
        JobID    string
        Tasks    []string
        Attempts int
    }

    // Outcome is the facts of one finished trial. The verdict is derived from facts elsewhere.
    Outcome struct {
        TaskID  string
        Attempt int
        Status  string // a small closed set of values; a defined type in real code
        Reward  *float64
        Error   string
    }

    // CancelInfo is why a job was stopped. It is an error so it can be a context cause.
    CancelInfo struct {
        By     string
        Reason string
    }

    Event struct {
        Kind    Kind
        TaskID  string
        Attempt int
        Outcome *Outcome    // KindTrialFinished
        Cancel  *CancelInfo // KindCanceled
        Err     error       // KindFailed
    }

    Summary struct{ Finished, Passed int }

    // Executor runs a request until every trial has finished, ctx is canceled or it fails,
    // reporting through sink. It stops its children when ctx is canceled.
    Executor interface {
        Execute(ctx context.Context, req Request, sink EventSink) (Summary, error)
    }

    // EventSink receives a job's events in order. An error from Emit fails the job.
    EventSink interface {
        Emit(ctx context.Context, ev Event) error
    }
)

const (
    KindJobStarted Kind = iota + 1
    KindTrialStarted
    KindTrialFinished
    KindCanceled
    KindFailed
)

// ErrCanceled wraps the CancelInfo cause when Execute returns because ctx was canceled.
var ErrCanceled = errors.New("execution canceled")

func (c CancelInfo) Error() string { return c.By + ": " + c.Reason }
```

### `internal/execution/vendorcli/executor.go` — the adapter is the anti-corruption layer

```go
package vendorcli

// Executor runs the vendor program. The vendor's flags, directory layout and result files are
// known here and nowhere else; vendorTrial never leaves the package.
type Executor struct {
    bin     string
    jobsDir string
}

func (e *Executor) Execute(ctx context.Context, req execution.Request, sink execution.EventSink) (execution.Summary, error) {
    cmd := exec.CommandContext(ctx, e.bin, buildArgs(req, e.jobsDir)...) // the only place flags exist
    // start the process; watch jobsDir/<JobID> and emit KindTrialStarted / KindTrialFinished as the
    // vendor's files appear (toOutcome is the one mapping); on exit, one reconciliation scan, so
    // a trial that finished between two ticks is not lost.
    // On ctx cancel: kill the process group, then
    //   return Summary{}, fmt.Errorf("%w: %w", execution.ErrCanceled, context.Cause(ctx))
    // after emitting KindCanceled with the CancelInfo cause.
    return summary, nil
}
```

### `internal/execution/executiontest/conformance.go` — one suite, every adapter

```go
// Conformance asserts what any Executor must do. Each adapter's test file calls it with its own
// constructor; the scripted fake and the vendor adapter (against a stand-in program) both pass.
// Script, PassAll, SlowAll, request and recordingSink are small helpers, elided here.
func Conformance(t *testing.T, newExecutor func(t *testing.T, script Script) execution.Executor) {
    t.Helper()

    t.Run("a job starts before any trial", func(t *testing.T) {
        t.Parallel()
        sink := &recordingSink{}
        _, err := newExecutor(t, PassAll(2)).Execute(t.Context(), request(2), sink)
        require.NoError(t, err)
        require.Equal(t, execution.KindJobStarted, sink.events[0].Kind)
    })

    t.Run("a cancel carries who and why", func(t *testing.T) {
        t.Parallel()
        ctx, cancel := context.WithCancelCause(t.Context())
        sink := &recordingSink{onEmit: func(ev execution.Event) {
            if ev.Kind == execution.KindTrialStarted {
                cancel(execution.CancelInfo{By: "test", Reason: "stop"})
            }
        }}
        _, err := newExecutor(t, SlowAll(3)).Execute(ctx, request(3), sink)
        require.ErrorIs(t, err, execution.ErrCanceled)
        require.Equal(t, execution.KindCanceled, sink.last().Kind)
        require.Equal(t, "test", sink.last().Cancel.By)
    })
}
```

Notes:

- The fake is a plain struct that plays a `Script` (pass, fail, error, slow, cancel on the Nth);
  it is what worker and orchestration tests use, so they never need the vendor program.
- Keep the vendor's real output as golden fixtures under the adapter's `testdata/` and assert
  `toOutcome` against them: a vendor format change then fails one test, in one package.
- Enforce the direction with a test that fails if anything outside the adapter imports it.
- Decide *where a verdict is made* separately: store facts from the tool, derive the verdict at
  read. Introduce a grader port only when a second grader exists.

## Database access examples (§12)

### sqlc (default)

```sql
-- query.sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
```

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "query.sql"
    schema: "migrations"
    gen:
      go:
        package: "db"
        out: "db"
```

```go
func (r *UserRepo) GetByID(ctx context.Context, id int64) (db.User, error) {
    return r.queries.GetUserByID(ctx, id)
}
```

Layout:

```
internal/userrepo/
├── repo.go
├── sqlc.yaml
├── query.sql
└── db/              # generated: models.go, query.sql.go, db.go — never hand-edit
```

### sqlx + `//go:embed` (alternative)

```go
//go:embed queries/get_user_by_id.sql
var getUserByIDSQL string

func (r *UserRepo) GetByID(ctx context.Context, id int64) (User, error) {
    var u User
    err := r.db.GetContext(ctx, &u, getUserByIDSQL, id)
    return u, err
}
```

Layout:

```
internal/userrepo/
├── repo.go
└── queries/
    ├── get_user_by_id.sql
    ├── insert_user.sql
    └── list_users.sql
```
