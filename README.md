
# Project R.E.D.

## Sovereign Knowledge Node Engine

Project R.E.D. rejects both centralized database monopolies and overly complex distributed consensus protocols. Instead, it systematically decouples the **Independent Data Layer** from the **Social Curation Layer**.

The engine operates as a stateless, high‑performance Go runtime that compiles raw Markdown files into visually polarized technical templates, dynamically injecting cryptographic integrity signatures on every request loop.

## 1. The Philosophy: Red Sovereignty

The internet is broken not from a lack of information, but from a fundamental crisis of architecture. Well‑meaning idealists continually attempt to fix this by proposing centralized chokepoints—a theoretical centralized archive designed to hold all human “how‑to” knowledge in a single, strictly moderated, unified database.

While noble in its intent to provide universal free education, all centralized systems suffer from structural flaws that guarantee its immediate compromise or destruction:

* **The Single Point of Failure:** By attempting to gather all technical knowledge onto a single domain, the Blue System creates an existential target for corporate lawfare and state suppression.
* **The Conflict of Interest:** The moment a platform’s revenue is derived from a percentage of material sales via affiliate links, its neutrality is permanently compromised. Centralized platforms reliant on contextual advertising financially incentivize the prioritization of expensive corporate partner materials over raw, cost‑effective tutorials.
* **The Jury‑Style Delusion:** Relying on an uncompensated “jury” of randomly selected expert verifiers to audit complex technical workflows is an economic fantasy. High‑earning technical professionals will not donate hundreds of hours of uncompensated labor to audit technical execution guides when keeping that operational knowledge exclusive is their livelihood.
* **The Gamified Trust Trap:** Implementing video‑game‑inspired mechanics—such as author “Elo” rankings or weighting votes based on account age—creates a closed aristocracy. This permanently locks out new, highly‑qualified outside experts who refuse to play popularity contests, turning a critical repository into a digital playground.
* **The Dependency Hell of Linear Hierarchies:** Forcing all human knowledge into a rigid, linear hierarchy where Level 4 strictly depends on Level 1 creates a catastrophic maintenance bottleneck. If a foundational tool or component becomes obsolete, the minor change at the bottom breaks every single cascading guide built on top of it.

## 2. How R.E.D. Fixes That

Project R.E.D. systematically dismantles the architectural vulnerabilities inherent in centralized knowledge repositories.

* **Eradicating the Single Point of Failure:** Open platforms rely mainly on a master domain, creating a massive target for corporate lawfare and global de‑indexing. R.E.D. operates as a stateless containerised engine, eliminating the centralized attack vector entirely.
* **Eliminating Financial Conflicts:** R.E.D. requires zero centralized funding, ensuring information remains free from commercial manipulation, corporate de‑indexing, and review‑bombing botnets.  
  **Ethical Node Monetization (Optional):** Individual node operators may choose to earn back their costs by including affiliate links or donation addresses within the guides they host. The revenue goes directly to the operator who maintains the node – never to a central entity. Readers can freely choose a different node if they prefer an ad‑free experience. No special code is needed: operators simply create a `Requirements.md` file inside a guide’s folder listing necessary parts (with their affiliate links), and link to it from the guide’s main `index.md` or introductory file. This keeps monetization transparent, decentralized, and entirely opt‑in.
* **Outsourcing Content Curation & Bot Defense:** Centralized platforms collapse trying to independently build anti‑bot algorithms and content moderation teams. R.E.D. outsources the social layer to multi‑billion dollar ecosystems like Reddit or Discord. By relying on networks with existing phone‑verification and automated bot mitigation, R.E.D. bypasses the need for a non‑profit to design complex, native bot‑detection firewalls.
* **Bypassing Dependency Hell:** B.L.U.E. enforces a rigid, linear learning hierarchy where a single obsolete foundational guide can collapse the entire structure. R.E.D. prevents this by leveraging native filesystem directories and dynamic versioning, allowing knowledge to adapt organically without cascading failures.
* **Resolving the “Spin‑Off” Paradox:** B.L.U.E. mandates one definitive guide per topic, yet paradoxically suggests forking contested guides during internal disputes, guaranteeing a fractured, redundant database. R.E.D. removes moderation logic from the runtime entirely; the end‑user’s local client seamlessly curates the best version based on established network consensus.

## 3. Core Architectural Counter‑Measures

To address the recurrent theoretical objections regarding security, DDoS attacks, content manipulation, and illegal uploads, the R.E.D. engine uses hard‑coded cryptographic and infrastructure barriers.

### A. Cryptographic Integrity Verification (Ed25519 + SHA‑256)

R.E.D. completely eliminates the threat of unauthorised guide modification using standard SHA‑256 hashes and **Ed25519 signatures**.

- Every `.md` file can be accompanied by a `manifest.json` inside its vault folder.  
  Example entry:
  ```json
  {
    "File.md": {
      "file_hash": "20bcee7014ad5cff...",
      "public_key": "66f32e250b0bafb2...",
      "signature": "61dcac5bd139dca0..."
    }
  }
  ```
- The engine calculates the SHA‑256 hash of the raw file on every request and compares it with the `file_hash` in the manifest.
- If the hash matches, it verifies the Ed25519 signature using the provided public key.
- Finally, it checks whether the public key is **trusted** – listed in the root‑level `contributors.json` file with a human‑readable name.
- If all checks pass, the article shows a green **✅ Verified Contributor** badge with the contributor’s name. The full SHA‑256 checksum is displayed in the footer for manual verification.
- If the file is unsigned, the hash mismatches, or the signer is not trusted, a **⚠️ Unverified / Unknown Origin** badge is shown.

This mechanism requires **no database** – the verification is purely file‑based, stateless, and auditable.

### B. Stateless Immunity to DDoS and Subpoenas

Centralized sites are easy targets for DDoS attacks and corporate subpoenas because they rely on massive, active SQL/NoSQL databases.

- The R.E.D. engine is **stateless**. It possesses no database layer to exploit, breach, corrupt, or subpoena.
- Because it serves raw Markdown natively from file storage, it operates with minimal memory overhead and zero database lookup latency. Containers can be replicated instantly across alternative IP addresses.

### C. Instant Decentralized Mirroring (Import API)

A sovereign network is only as strong as its ability to replicate data rapidly.

- R.E.D. features a built‑in `/import` ingestion endpoint (accessible via the admin panel).
- Node operators can pass any raw Markdown URL, GitHub repository link, or archive (`.zip`, `.tar.gz`) to the engine. The Go runtime will instantly fetch the content, automatically reconstruct the original folder hierarchy on the local drive, and serve it **without requiring a server restart**.
- This allows networks of nodes to rapidly mirror critical data before a master source is taken offline.

### D. Admin Panel & Token Protection

- The admin panel is available at `/-/admin` and is protected by a random `adminToken` generated during installation.
- From the panel you can:
  - Import new knowledge bases (single files, GitHub repos, archives).
  - List all currently synced sources.
  - Remove a source and optionally delete its local files.
  - Persist import rules to `config.json` so they are re‑synced on container start.
- Use the provided `manage-token.sh` (Linux/macOS) or `manage-token.ps1` (Windows) scripts to regenerate the admin token at any time.

## 4. Architecture: Single‑Container Deployment (Clearnet Only)

Earlier versions experimented with dual‑tier (Clearnet + Tor) isolation, but this proved operationally insecure and unnecessarily complex. The current reference deployment uses a **single container** behind a Caddy reverse proxy.

```mermaid
graph TD
    Inet((Global Internet))
    Inet --> Caddy[Caddy Proxy<br/>Port 80/443]
    Caddy --> RED[RED Engine Container<br/>(Go + Goldmark)]
    RED --> Volume[(Host Volume /data)]
    Volume --> Markdown[Markdown Files<br/>+ manifest.json<br/>+ contributors.json]
```

All components run as standard Podman (or Docker) containers, orchestrated via `podman-compose` (or `docker compose`). The RED engine listens on port `8080` internally, while Caddy provides automatic HTTPS (or plain HTTP on port 80 for local testing).

## 5. Technical Stack Blueprint

- **Language & Runtime:** Go (Golang) — compiles down to a single memory‑safe static binary with zero virtual machine overhead.
- **Markdown Parser:** `goldmark` — highly efficient, fully CommonMark compliant.
- **HTML Sanitizer:** `bluemonday` – strict user‑content policy to prevent XSS.
- **Container Environment:** Multi‑stage minimal Docker/Podman pipeline (Alpine) for complete process isolation and instant replication.
- **Orchestration:** `podman-compose` or `docker compose`.

## 6. Installation & Deployment

### Prerequisites
- **Podman** (recommended) or **Docker** (with `docker compose` V2)
- **Git**
- **Bash** (Linux/macOS) or **PowerShell** (Windows)

### Automatic Installation

#### Linux / macOS
```bash
curl -sSL https://raw.githubusercontent.com/RED-Collective/red-engine/main/install-red-engine.sh | bash
```
Or clone and run manually:
```bash
git clone https://github.com/RED-Collective/red-engine.git
cd red-engine
./install-red-engine.sh
```

#### Windows (PowerShell as Administrator)
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force
iex (irm https://raw.githubusercontent.com/RED-Collective/red-engine/main/install-red-engine.ps1)
```

The installer will:
- Clone the repository (if not already inside it).
- Create the `./data` directory.
- Generate a `config.json` with a cryptographically random `adminToken` (displayed once).
- Start the container using `podman-compose up --build -d` (or fallback to `docker compose`).
- Show the node URL (`http://localhost`) and admin panel URL (`http://localhost/-/admin`).

### Manual Setup (Docker/Podman)

1. **Clone the repository**
   ```bash
   git clone https://github.com/RED-Collective/red-engine.git
   cd red-engine
   ```

2. **Prepare configuration**
   ```bash
   cp config.json.example config.json   # or let the installer create it
   # Edit config.json – set a strong adminToken, adjust dataDir if needed
   ```

3. **Build and start**
   ```bash
   podman-compose up --build -d
   # or: docker compose up --build -d
   ```

4. **Verify**
   - Node: `http://localhost`
   - Admin panel: `http://localhost/-/admin` (use the token from `config.json`)

### Managing the Admin Token
```bash
# Linux/macOS
./manage-token.sh

# Windows
.\manage-token.ps1
```
The script reads `config.json`, displays the current token, and optionally generates a new secure token. After changing the token, restart the container:
```bash
podman-compose restart red_engine
```

## 7. Simulating Knowledge Base Structures

Rather than enforcing a brittle, top‑down linear hierarchy that breaks under dependency hell, R.E.D. leverages your computer’s native filesystem directory tracking. Group, version, and fork your files within folders inside your local storage volume dynamically:

```plaintext
data/
└── remote/
    └── solar-array-build/
        ├── 00-index.md
        ├── 01-wiring-schematics-v1.0.md
        └── 02-inverter-fusing.md
```

The node pathing resolves these automatically into accessible routes:

* `http://your-node/remote/solar-array-build/00-index`

## 8. Cryptographic Verification in Depth

### `contributors.json` (root of the repository)
```json
[
  {
    "name": "StandardCodebase",
    "public_key": "66f32e250b0bafb2683ae987e65a390901f030ccbc7745435d99f35da1bfccf1",
    "contact-information": { ... }
  }
]
```
Only public keys listed here are considered **trusted**. The `name` field is displayed in the verified badge.

### `manifest.json` (inside any vault folder)
```json
{
  "relative/path/to/file.md": {
    "file_hash": "sha256-of-file-content",
    "public_key": "same-as-in-contributors.json",
    "signature": "ed25519-signature-of-the-file-content"
  }
}
```
- The signature must be computed **over the raw file content** (not over the hash).
- The engine automatically walks the entire `dataDir`, finds all `manifest.json` files, and matches each `.md` file to its entry.

### Verification Flow (per request)
1. Read `.md` file → compute SHA‑256 → compare with `file_hash` from manifest.
2. If hash matches, decode public key and signature (hex).
3. Verify using Go’s `ed25519.Verify(publicKey, fileContent, signature)`.
4. If valid, look up the public key in `contributors.json` → if found, set `Verified=true` and `Author=contributor.name`.
5. Otherwise, set `Verified=false` and `Author="Unverified / Unknown Origin"`.
6. Inject badges and SHA‑256 footer into the HTML template.

This design ensures that **readers always know the provenance** of every guide, and that **malicious alterations** (even a single character) break the hash, immediately flagging the content as untrusted.

## 9. Curation Philosophy: The Sovereign Web‑of‑Trust

1. **The Software Only Knows State:** The Go runtime handles zero moderation logic. It treats strings objectively.
2. **Cryptographic Peer Review:** Instead of relying on a centralized platform’s voting buttons or video‑game‑inspired Elo rankings, users and curators sign content hashes with their independent cryptographic keys (Ed25519).
3. **End‑User Agency:** You, the reader, choose which curators to trust. Your local index client aggregates links signed by your trusted network.

The creator of the Blue System lamented that he lacked the millions of dollars needed for a startup, resigning himself to the belief that his vision would never be a reality. He failed because he thought like a corporate founder seeking capital.

**We do not need capital. We do not need permission. We do not need a master website.**

It’s about time we stop waiting for someone or something to fix our problems for us. No hero is coming.

**Claim Your Agency. Run a Node.**
