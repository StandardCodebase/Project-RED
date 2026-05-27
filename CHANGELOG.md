# Changelog

All notable changes to Project R.E.D. will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0-alpha] - 2026-05-25

### 🚀 Initial Alpha Release
Project R.E.D. (Sovereign Knowledge Node Engine) enters its first public alpha. The core engine is now fully functional as a zero-database, stateless Go runtime that compiles cryptographic-signed Markdown into a polarized UI over dual-tier networks (Clearnet + Tor).

### Added
- **Remote Markdown Ingestion:** Added a `/import` POST endpoint. The engine can now fetch raw `.md` files directly from remote URLs.
- **Dynamic File Structuring:** The Go backend now automatically parses imported URLs and replicates their domain/path structure locally inside the `data/` volume.
- **Stateless Manifest Generation:** Added a `/manifest` endpoint. The Go runtime dynamically walks the local file tree and serves a JSON map of all available knowledge files on demand.
- **Image and Diagrams Support:** Added a simple function to render mermaid diagrams and display images on the website. I am still working on making the download functionality for entire guides.
- **Hierarchical "QTreeView" UI:** Engineered a dynamic JavaScript frontend component that groups flat directory paths from the manifest into collapsible, nested folder accordions.
- **Off-Canvas Navigation:** Implemented a pure CSS/JS sliding side panel for cross-guide navigation, eliminating the need to return to the root index.
- **Dual-Tier Docker Matrix:** Published the `docker-compose.yml` for instantly deploying the air-gapped Onion Vault alongside the Clearnet Gateway.
- **Automated Caddy setup:** Now the container will generate the necessary files every time a new node goes online.
- **Cryptographic Hashing:** The engine now automatically calculates and serves a strict `X-RED-Content-Hash` (SHA-256) on every request loop to prevent file tampering.
- **Polarized Theme Engine:** Added native Light/Dark CSS variables scoped to the `:root` level, completely independent of external CSS frameworks.

### Changed
- Replaced the hardcoded `/guides` HTML template with the dynamic JSON-driven sidebar.
- Updated `layout.html` to utilize inline, cache-busting CSS delivery to prevent browser caching during rapid UI prototyping.
- Restructured `style.css` to remove fixed square bounds on header buttons, allowing for clean text-based icon scaling.

### Security
- **Directory Traversal Protection:** Hardened `filepath.Clean` and `filepath.Abs` boundaries in `main.go` to prevent malicious actors from requesting files outside the designated `RED_DATA_DIR`.

---
*Note to Node Operators: This is an alpha release. Ensure you are pulling the latest `main.go` binary rebuilds via Docker when updating.*
