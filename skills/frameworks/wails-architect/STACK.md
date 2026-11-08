# Stack Versions

| Dependency | Pinned version | Purpose |
|---|---|---|
| go | 1.27.x | Language and `wails3` toolchain |
| github.com/wailsapp/wails/v3 | 3.1.0 | Desktop/mobile shell (desktop stable, mobile experimental) |
| wails3 CLI | ships with the pinned Go toolchain | `dev`, `build`, `task`, `tool package`, `generate appimage` |
| gtk4 | 4.22 | Linux window toolkit (needed to build the shell) |
| webkitgtk | 6.0 (2.54) | Linux webview |
| JDK | 21 (full, with `javac`) | Android Gradle build |
| Android SDK | targetSdk 36, minSdk 21 | Android build |
| Android App Bundle | upload-key signed | Google Play internal testing |

## Notes

- **Inherits the `go-architect` stack** for module layout, logging and tooling. `STACK.md` there pins shared Go libraries.
- **The SPA is not part of this stack.** Pin React/Vite/pnpm in [react-architect/STACK.md](../react-architect/STACK.md); the shell only embeds the build.
- **Platform coverage in the reference implementation:** Linux and Android were run. Windows and macOS window options compile but were not exercised. iOS is out of scope.
- **AppImage is large** (~118 MB) because it bundles GTK/WebKit; `.deb`, `.rpm` and `.pkg.tar.zst` are ~7 MB.
- **No keychain API** in Wails 3.1.0; persist sessions yourself.

_Last reviewed: 2026-11-08_
_Skill version at last review: 1.0.0_
