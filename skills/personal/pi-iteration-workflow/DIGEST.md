# Pi Iteration Workflow (digest)

Edit-test loop for `platform/robot` (voldemorbot): edit on Windows, deploy to a Pi over SSH, test on hardware.

## Always-on rules
- Prefer `rpi-5-direct` (`192.168.250.2`, direct Ethernet) over WiFi aliases or the Cloudflare tunnel.
- Account is `ralvarezdev`, not `pi`. Identity file `D:\Secrets\ssh\raspberrypi\id_ed25519_ralvarezdev`.
- mDNS (`.local`) flakes. Use the raw IP with explicit `-i`.
- Loop: edit; commit and push `master`; `git pull` on the Pi.
- `pixi.toml`/`pixi.lock` changed: `~/.pixi/bin/pixi install -e dev`.
- `ros2_ws/src` changed: `~/.pixi/bin/pixi run -e dev build-ws`. Editing source alone does nothing.
- Restart `voldemorbot-pi-zero.service` (or `-pi5`), then `systemctl status` and `journalctl --since` for node output.
- Lock `pixi.lock` on the Pi 5, never Windows.
- Run long jobs (`pixi install`, `colcon build`) in `screen`. Detach `Ctrl+A` `D`; reattach `screen -d -r work`.
- `/tmp` is wiped on reboot.
- PowerShell: no nested quotes plus `|` in one ssh string. Run `ssh -t`, then the command block inside the Pi shell.
- Pi Zero sudo cache lasts ~15 min after an interactive sudo.

## Hard constraints
- Pi Zero sudo needs a password: run it in the user's own terminal, never a non-interactive tool call.
- If sudo reports a password is required, hand the command back. Do not retry.

Invoke the full skill when: deploying robot code to a Pi, regenerating `pixi.lock`, setting up `screen`, restarting the robot service, or an SSH session misbehaves.

Full guide: invoke the pi-iteration-workflow skill.
