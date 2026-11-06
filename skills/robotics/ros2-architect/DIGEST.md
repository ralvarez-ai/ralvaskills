# ros2-architect (digest)

ROS2 standards across Jazzy, Kilted and Lyrical: colcon and ament layout, lifecycle nodes, explicit QoS, services and actions, declared parameters, Python launch; Pixi + RoboStack dev environment.

## Always-on rules
- Jazzy (LTS) is the default for new production systems; Lyrical (LTS) once the ecosystem catches up; Kilted only as a bridge, migrate before Dec 2026.
- Pixi owns the dev environment; commit `pixi.lock`; no `sudo apt install` on contributor machines.
- `colcon build --symlink-install` in dev; drop it for release builds.
- `src/` holds packages; `_msgs` packages hold only interfaces; `_bringup` holds launch and config, no source.
- One package per responsibility; one node per job; no business logic in `main` or launch files.
- Lifecycle nodes for anything holding an open connection; standard `Node` otherwise; single-threaded executor by default.
- Explicit QoS on every publisher and subscriber; publisher and subscriber QoS must match.
- Sensor streams: best-effort + KEEP_LAST(N). Low-rate state: reliable + KEEP_LAST(10).
- Services for operations under ~100 ms; actions for navigation-, manipulation- and planning-like work.
- Declare every parameter with a `ParameterDescriptor`; validate on change; YAML config lives in `_bringup`.
- Python launch files (`*.launch.py`); launch arguments for environment differences.
- `package.xml` `format="3"`; declare every dependency; `ament_lint_*` in every package's tests.
- Log through the ROS2 logger, never `printf` or `print`.

## Hard constraints
- Never call a service synchronously inside a topic callback; use `call_async` or `async_send_request`.
- Reliable + KEEP_ALL on a high-rate topic is a memory leak; never do it.
- Never mix ROS distros in one workspace.
- Never put code at the workspace root.
- Rolling only for contributor development against ROS2 itself.

Invoke the full skill when: scaffolding or reviewing a ROS2 workspace, choosing a distro, writing QoS, lifecycle, service or action code, launch files, or launch_testing.

Full guide: invoke the ros2-architect skill.
