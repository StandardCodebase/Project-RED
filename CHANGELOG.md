## [Unreleased] - 2026-05-29

### Added
- **Native Git Engine:** Integrated `go-git/go-git/v5` for true delta-pulling and cloning, eliminating the engine's reliance on host OS shell commands and bypassing container permission traps.
- **Granular Memory Hot-Reloading:** Added `Store.UpdateFiles(changedPaths []string)` to surgically patch the active memory map. The engine now drops and re-renders only the specific files modified in a commit or local save, eliminating CPU spikes and full-site downtime during syncs.
- **Smart Webhook Routing:** Added intelligent JSON payload parsing to the `/-/webhook/sync` endpoint. Webhooks now extract the origin URL and only trigger delta-pulls for matching repositories.
- **Container-Safe Local Polling:** Added `radovskyb/watcher` to bypass Docker/Podman hypervisor limitations where `inotify` events fail to cross into the container.

### Changed
- **Replaced `fsnotify`:** Local file watching is now handled by a 2-second interval background poller, which directly feeds into the new granular memory hot-reloading module.
- **Silent Background Poller:** Refactored the 1-minute brute force loop in `cmd/red/main.go`. It now uses `fetch.PullDelta` to silently check for remote Git changes without downloading entire repository archives.
- **Installation Scripts:** Updated `install-red-engine.sh` and `install-red-engine.ps1` to automatically assign global read/write permissions (`chmod 777` and `icacls Everyone`) to the `data/` volume. 
- **Docker Dependencies:** Updated `Dockerfile` to install `ca-certificates`, `git`, and `openssh` directly into the Alpine container for native Git support.

### Fixed
- **Podman Permission Trap:** Prevented the restricted `reduser` (UID 1000) from being locked out of the `data/` directory when the host machine auto-creates missing volume mounts as `root`.
- **Mutex Panic in Store:** Fixed a fatal runtime concurrency bug in `store.go` where a deferred `mu.Unlock()` would cause a panic if security definitions (`manifest.json` or `contributors.json`) were modified.
- **Webhook Global Loop Bug:** Fixed an issue where a single webhook ping would force the engine to blindly re-download every tracked repository in the configuration list.
- **ZIP Archive Loop Bug:** Changed default URLs in `config.json` from `/archive/HEAD.zip` to `.git` to prevent the background sync from repeatedly destroying and recreating directories every 60 seconds.