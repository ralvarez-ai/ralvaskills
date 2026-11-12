# cli-tool-architect (digest)

Cross-language CLI standards: subcommand structure, flag > env > config > default precedence, TOML in XDG, stdout for data and stderr for logs, `--output json|yaml`, exit codes, NO_COLOR, completions.

## Always-on rules
- `tool <noun> <verb>` or `tool <verb>`; one verb per subcommand; `tool` alone prints help.
- Positional args only for the required noun, and no more than two in most cases; everything else is a flag.
- Short flags only for the top ~5; boolean flags are pure switches (`--no-foo` to negate).
- Precedence is fixed: flag > env var > config file > built-in default. Env prefix is the upper-snake tool name.
- Provide `tool config show` with the source of each resolved value.
- Config: TOML at `$XDG_CONFIG_HOME/<tool>/config.toml`; discovery is `--config`, then `TOOL_CONFIG`, then `./<tool>.toml`, then the XDG path. One canonical filename.
- `tool init` refuses to overwrite without `--force`.
- Data to stdout; logs, progress and errors to stderr.
- Every list or get command supports `--output json|yaml` (`-o`); JSON schema is stable; JSONL for lists.
- Exit codes: 0 success, 1 failure, 2 misuse, 130 SIGINT; document every custom code.
- Color: auto-detect TTY; respect `NO_COLOR` and `FORCE_COLOR`; `--color=auto|always|never`.
- Progress only for work over ~2 seconds; suppressed when not a TTY; Ctrl-C exits 130.
- Errors say what failed, why, and what to try next.
- Ship completions for bash, zsh, fish and pwsh.
- `tool --version` and `tool version` print the same thing; semver.

## Hard constraints
- Never run an implicit default action.
- Never print an error to stdout and exit 0.
- Never mix data and logs across streams.
- Never reuse an exit code across categories within one tool.
- Never auto-upgrade; `tool upgrade` is opt-in.
- No stack traces in user-facing errors.

Invoke the full skill when: designing or reviewing a CLI, or when you need the Go (cobra, pflag, viper) or Python (typer) recipes, exit-code, output-format, error-message or distribution details.

Full guide: invoke the cli-tool-architect skill.
