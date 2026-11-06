# grafana-architect digest

Operator side of observability: dashboards and alerts as code (Grizzly), one service per folder, one question per panel, actionable alerts.

## Always-on rules
- Dashboards live in git and are applied with Grizzly (`grr`): staging on every push, prod on tagged releases.
- Build in the UI only to explore; export and commit once a dashboard earns a home.
- One folder per service; folder permissions match service ownership; fleet dashboards go in top-level `Platform` (SRE).
- Every panel answers one question; title is a question or noun phrase; y-axis unit is mandatory.
- Keep time range consistent across a dashboard unless comparison is its purpose.
- Bound legends: templatize queries (`sum by (route)`); high-cardinality data becomes heatmap or top-N.
- Variables come from label queries (`label_values(...)`); standard variables: `service`, `environment`, `region` if multi-region.
- Don't split a dashboard on a variable (`Orders Prod` and `Orders Staging` are one dashboard).
- Every service has exactly one SLO dashboard: compliance, error budget, burn rate, top contributors.
- Every alert is actionable; a "look in the morning" alert is a dashboard.
- Severity is two-tier: `severity=page` and `severity=ticket`, no middle ground.
- SLO alerts use multi-window, multi-burn-rate (fast and slow together).
- Group alerts by `team` and `service` labels so silences target one service.
- Soft cap ~15 panels per dashboard; split beyond that.
- Data sources are provisioned via Grizzly; one per signal per environment; service-account keys, rotated quarterly.

## Hard constraints
- Never edit production dashboards in the UI; Grizzly refuses drift, so pull the change into JSON.
- `for:` is mandatory on every alert (minimum `2m` on noisy signals, `5m`+ for slow burns).
- `runbook_url` annotation is mandatory on every alert.
- No CPU > 80% alerts, no per-instance alerts on horizontally scaled services, no static thresholds on traffic-growing metrics.
- "Disk full" alerts need `for: 15m`.
- Anonymous access off in every environment, staging included; SSO for humans.
- No personal folders for production.
- Dashboards use read-only data sources; write access goes through audited service accounts.

Invoke the full skill when: designing dashboards or folder layout, writing alert rules or Grizzly YAML, building SLO burn-rate alerts, or auditing Grafana permissions or drift.

Full guide: invoke the grafana-architect skill.
