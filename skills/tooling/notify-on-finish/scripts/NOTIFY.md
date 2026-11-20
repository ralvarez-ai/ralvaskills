# notify-on-finish - backends, state and troubleshooting

Companion to [SKILL.md](../SKILL.md). Everything here is what the script already does; read it when a notification or a sound does not arrive.

## Notification backends

Tried in order, first success wins, all best-effort:

| # | Backend | Platform | Notes |
|---|---|---|---|
| 1 | `notify-send` | Linux, BSD | Standard freedesktop client; needs a daemon (`org.freedesktop.Notifications`) |
| 2 | `gdbus call … Notify` | Linux | Same D-Bus call without `libnotify`; survives a missing `notify-send` |
| 3 | `osascript -e 'display notification'` | macOS | Notification Center; app name shows as Script Editor |
| - | `dunstify`, `makoctl`, `swaync-client` | Linux (tiling WMs) | Daemon-specific clients; use them when the generic path is rejected |

`notify-send` urgency maps to what each daemon understands: `normal` on success, `critical` on failure (sticky in most daemons, so a failure is not lost while you are away).

## Sound backends

Tried in order for `ok` → `complete`, for failure → `bell`:

| # | Backend | Played |
|---|---|---|
| 1 | `canberra-gtk-play -i <name>` | freedesktop *sound event* name - theme-aware, preferred |
| 2 | `pw-play` / `paplay` on `/usr/share/sounds/freedesktop/stereo/<name>.oga` | PipeWire / PulseAudio, needs `sound-theme-freedesktop` |
| 3 | `afplay /System/Library/Sounds/{Glass,Basso}.aiff` | macOS |
| 4 | terminal bell | last resort |

`--quiet` skips sound only. `--silent` skips sound **and** notification (log only). `--no-notify` skips the desktop notification only - useful when you want the chime but not the popup.

## State layout

```
$XDG_STATE_HOME/notify-on-finish/     # default: ~/.local/state/notify-on-finish
  <id>.meta      # id, label, cmd, state, code, start, end, duration, log, pid
  <id>.log       # stdout+stderr of the command
  <id>.launch    # anything the detached wrapper printed to stderr (should stay empty)
```

`<id>` is `<UTC timestamp>-<launcher pid>`. Files older than 30 days are pruned on the next launch; `--clean` drops every finished run now and keeps the running ones.

Never edit a `.meta` by hand while a run is in flight - the worker rewrites it atomically on completion.

## Exit codes of the wrapper itself

| Code | Meaning |
|---|---|
| command's own | detached: reported in `--status`; `--wait`: the wrapper exits with it |
| 2 | misuse - no command, `--label` without a value, unknown flag |

## Troubleshooting

**No notification, no error.** No daemon on the session bus. `notify-on-finish --doctor` shows `daemon: no reply`. Start one (`dunst`, `mako`, `swaync`, GNOME/KDE shell), or run `--silent`. Over SSH or in CI there is no session bus - expected, not an error, and the exit code is still correct.

**Notification appears but silent.** DND / focus mode, or `canberra-gtk-play` cannot reach the audio server. Test the sound alone: `canberra-gtk-play -i complete`. In a container, `/run/user/$UID/pulse` or `PIPEWIRE_RUNTIME_DIR` is usually missing.

**`canberra-gtk-play: command not found`.** Install `libcanberra-gtk3` (or `libcanberra-gtk2`); the `pw-play` path still works if the freedesktop theme is present.

**No `complete.oga`.** `sound-theme-freedesktop` is not installed; only the canberra event name path will work.

**Run shows `running` forever.** The command is still alive, or it inherited a tty it should not have. Find it with `pgrep -af notify-on-finish`, read `<id>.log`, and kill the pid from the meta. The wrapper itself has no timeout by design - a build that legitimately takes two hours must not be killed at 30 minutes.

**macOS: nothing fires.** Terminal notifications need permission for the calling app (System Settings → Notifications → Script Editor / your terminal). `afplay` path does not need it.

**Windows.** Native path is `osascript`'s absence + no `notify-send`; wrap in WSL and use `--silent`, or install `wsl-notify-send` and symlink it as `notify-send` so backend #1 picks it up.

## Quick self-test

```bash
notify-on-finish --doctor
notify-on-finish --label "ok path"   -- sh -c 'sleep 2; echo fine'
notify-on-finish --label "fail path" -- sh -c 'sleep 2; echo boom >&2; exit 3'
notify-on-finish --status
```
