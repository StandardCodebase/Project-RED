# Changelog

All notable changes to the `red-engine` project are documented in this file. This version marks a significant architectural shift from the initial `project-red` codebase to the current modular `red-engine` framework.

## [2.0.0] - Architecture Overhaul (Current Release)

This release represents a complete refactoring of the internal structure to improve modularity, maintainability, and scalability.

### Changed (Architectural Improvements)
- **Project Structure**: Relocated the entry point from the project root to `cmd/red/main.go` to follow standard Go project layouts.
- **Routing Engine**: Completely replaced the basic `internal/handler` structure with a comprehensive `internal/router` package. Routing is now decoupled into specific responsibilities:
    - `api.go`: Handles API-specific endpoints.
    - `raw.go`: Manages raw file/data access.
    - `serve.go`: Dedicated to serving HTTP responses.
    - `sync.go`: Manages synchronization logic.
    - `router.go`: Centralized router initialization.
- **Rendering**: Enhanced the `internal/render` package to better support the new routing structure and template rendering requirements.

### Added
- **`internal/fetch`**: New module dedicated to fetching remote or local data, enabling the synchronization features.
- **`internal/store`**: New persistent storage interface, replacing ad-hoc file handling with a defined store implementation.
- **Infrastructure & DevOps**:
    - **Caddy Integration**: Added `caddy_routing/Caddyfile` for robust web server configuration and reverse proxying.
    - **Deployment Scripts**: Added `install-red-engine.sh` and `install-red-engine.ps1` for automated setup on Linux and Windows environments.
    - **Token Management**: Added `manage-token.sh` and `manage-token.ps1` to handle authentication tokens securely.
    - **CI/CD**: Introduced GitHub Actions workflows (`.github/workflows/discord.yml`) for automated notifications.
    - **Docker Support**: Updated `Dockerfile` and `docker-compose.yml` for containerized development and production deployment.

### Removed
- **Legacy Packages**: Removed `internal/fs` and `internal/handler` as their functionality has been superseded by the `internal/router` and `internal/store` modules.

---

*The baseline architecture providing the foundation for the current rewrite.*

### Core Components
- Root-level `main.go` entry point.
- Simple `internal/config` and `internal/fs` modules.
- Basic handler architecture (`internal/handler`).
- Static asset serving via standard directory structure.
