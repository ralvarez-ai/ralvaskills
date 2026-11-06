# Python Architecture Standards (digest)

Python 3.14 enterprise standards: typing, immutable dataclasses, Protocol DI, asyncio discipline, pytest 9, psycopg with .sql files.

## Always-on rules
- Built-in generics (`list[str]`, `X | None`); read annotations via `annotationlib.get_annotations()`.
- `typing.NewType` for distinct domain IDs; `Literal` for narrow function-local flags; `TypedDict` for known-shape mappings.
- `Enum` (with `__str__` overridden) for closed domain state sets; `StrEnum` / `IntEnum` only when string or int interop is required.
- `typing.override` on every overriding method.
- Generators for large, unbounded, or cursor-backed sequences (`itertools` for lazy composition); `list[T]` for small bounded results.
- `@dataclass(slots=True, frozen=True)` for DTOs and value objects; `order=True` for comparable value objects.
- Pydantic only at application boundaries (API, DB rows, config); plain dataclasses for domain logic.
- Typed models over `dict[str, Any]` and raw tuples for stable shapes.
- `typing.Protocol` over deep ABC inheritance; define protocols where consumed.
- Dependencies passed into `__init__`; no globals; `contextvars` only for unavoidable request state.
- Sync I/O or CPU work via `asyncio.to_thread()`; `asyncio.TaskGroup` over bare `gather`.
- `concurrent.interpreters` for CPU-bound parallelism.
- `pathlib.Path` for all paths, never `os.path` strings; `importlib.resources.files(__package__)` for shipped SQL and templates.
- Imports in three groups (stdlib, third-party, local); minimal `__init__.py` with `__all__`.
- Exceptions: one base exception per module; chain with `raise ... from err`; handle ExceptionGroup with `except*`; no `return`/`break`/`continue` in `finally` (PEP 765).
- Data access (recommended pattern, not mandatory): psycopg 3, raw SQL in `.sql` files loaded at import; alembic for migrations. No ORM by default; SQLAlchemy needs an ADR.
- Tooling: uv (`uv sync --frozen` in CI), ruff, `mypy --strict`, pytest 9 + pytest-asyncio.
- CI gates are mandatory and run separately: `ruff check`, `ruff format --check`, `mypy`. Treat warnings as errors.

## Hard constraints
- Never legacy `typing.List` / `typing.Optional`.
- Never bare `except:`; never bare `except Exception` around TaskGroup failures.
- Never a mutable literal as a parameter or bare dataclass field default.
- Never block the event loop.
- Never instantiate external clients inside a class.
- Never `__file__`-relative paths for shipped assets.
- Never the `unittest` module.
- Never concatenate user input into SQL; bind parameters.
- Never let `--fix` silently remove commented-out code, unused imports, or unused locals (ERA, F401, F841).
- Never a blanket `ignore_missing_imports`; scope it to named untyped dependencies.

Invoke the full skill when: writing, reviewing, or scaffolding Python code, or setting up ruff and mypy config.

Full guide: invoke the python-architect skill.
