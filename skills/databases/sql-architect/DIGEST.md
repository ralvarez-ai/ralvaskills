# sql-architect digest

SQL standards for schemas, queries, migrations and database layers. PostgreSQL 18 primary; MySQL 9 and SQLite 3.53 noted.

## Always-on rules

- PK: `id UUID PRIMARY KEY DEFAULT uuidv7()` on every table. Natural identifiers stay UNIQUE.
- Tables snake_case plural; columns snake_case singular; FKs `<reftable_singular>_id`.
- Declare FKs in the DB, `ON DELETE RESTRICT` by default. NOT NULL by default. Push invariants into CHECK constraints.
- Soft delete: `deleted_at timestamptz NULL`, filtered by partial indexes. Hard delete only for GDPR/compliance or append-only event tables.
- Index every FK column. Add other indexes only when EXPLAIN ANALYZE proves the need.
- Explicit column lists, never `SELECT *`. Use RETURNING on writes. Use ON CONFLICT for upserts.
- Paginate by cursor `(created_at, id)`, not OFFSET.
- Prevent N+1: one JOIN or a batched `WHERE id IN (...)`.
- Keep transactions short. Roll back explicitly on error. Retry 40001 under SERIALIZABLE.
- Migrations: forward-only, idempotent where possible, one change each, named `NNNNNN_snake_case_description.sql`.
- Large tables: add nullable, backfill in batches, then add NOT NULL or the index.
- JSONB only for genuinely schemaless payloads. Structured data gets columns.
- Read EXPLAIN ANALYZE before optimising.

## Hard constraints

- Parameter binding always. Never concatenate user input into SQL.
- App runtime role: no CREATE, DROP or ALTER. Migrations run as a separate elevated role.
- No secrets in connection strings stored in the repo.

Invoke the full skill when: designing schemas, writing or reviewing queries, writing migrations, indexing, choosing transaction isolation, or auditing a database layer. Per-engine detail (MySQL, SQLite, PostgreSQL types) is in the skill's ENGINES.md and POSTGRES.md.

Full guide: invoke the sql-architect skill.
