# protobuf-architect digest

Protocol Buffers (proto3) schema standards: package naming, protovalidate, buf toolchain, field-number reservation, well-known types. Language-agnostic.

## Always-on rules

- Package equals directory path: `<org>.<product>.<resource>.<version>`. Per-language `option *_package` in every file, set up front.
- One `service` per file.
- Fields snake_case. Messages PascalCase, singular. Enums PascalCase; values UPPER_SNAKE_CASE prefixed with the enum name.
- Every enum has `*_UNSPECIFIED = 0`.
- Use `optional` for presence, not primitive wrappers. Avoid `map<K, V>` when ordering or evolution matters.
- Field numbers 1-15 for hot fields (1 byte on the wire).
- Never reuse a field number: `reserved` deleted numbers and names. Never change a field's type.
- Version lives in the package path (`v1`), never in message names. A new major is a new `vN` directory.
- Additive changes stay in the version. Breaking changes go to a new major. Deprecate before remove.
- Validate with protovalidate (CEL), server-side at the boundary. Skip protoc-gen-validate.
- Generate with `buf generate`; plugin versions pinned in buf.gen.yaml. Never hand-edit generated code.
- Run buf lint and buf breaking on every PR.
- Use well-known types (Timestamp, Duration, FieldMask, google.type Date and Money). Timestamp over int64 unix time.
- `Struct` signals bad design. Share `types/v1/` only when the model is provably the same. Never put services in the shared package.

## Hard constraints

- Never reuse a field number. Never change a field's type.
- Never pin plugins to `:latest`.
- Never hand-edit generated code.

Invoke the full skill when: designing or evolving .proto files, versioning a package, configuring buf lint, breaking or generate, or choosing validation rules. Language-specific generated-code idioms belong to the language skill.

Full guide: invoke the protobuf-architect skill.
