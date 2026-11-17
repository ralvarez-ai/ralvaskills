# notify-on-finish (digest)

Purpose: run a command detached, then play a sound and raise a native notification when it ends — without blocking the turn or losing the exit code.

## Always-on rules

- Launch background work with `scripts/notify-on-finish --label "<intent>" -- <command>`, never with `&`, `nohup`, `disown`, or a hand-rolled `cmd && notify-send`.
- Always put `--` before the command; the label is the notification title, written for a human (`"nightly test suite"`, not raw argv).
- Use it for work expected to exceed ~60 s; shorter commands just run normally.
- After launching, do not `sleep`-poll. Report the id and check later with `notify-on-finish --status` (TSV) or `--status --json`.
- `--wait` only when the next step depends on the result and a human needs the notification; it blocks and returns the command's exit code.
- `STATE` is `running`, `done`, `failed` (1–127) or `killed` (≥ 128). On failure read the tail of the `LOG` column, not the whole file.
- Report outcome, duration, the one line that matters, and the log path.
- `--quiet` = no sound; `--silent` = log only; `--no-notify` = no popup. Use `--silent` on SSH, CI and headless machines.
- Run `--doctor` when notifications do not arrive; never retry one — say what is missing and continue.

## Hard constraints

- Never wrap a command that needs a tty or reads stdin — the wrapper sets stdin to `/dev/null`.
- Never let a missing notification daemon fail the task: notification and sound are best-effort, the exit code is the command's.
- Never take the exit code from `notify-send` or from the launcher; detached runs report it through `--status`.
- No daemon, no polling loop, no hand-edited `.meta` while a run is in flight.

Invoke the full skill when: a command will take more than a minute, the user asks to leave something running and be told when it finishes, or a notification/sound does not arrive.

Full guide: invoke the notify-on-finish skill.
