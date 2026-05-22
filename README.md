
# Project R.E.D

## Sovereign Knowledge Node Engine

Project R.E.D. rejects both centralized database monopolies and overly complex distributed consensus protocols. Instead, it systematically decouples the **Independent Data Layer** from the **Social Curation Layer**.

The engine operates as a stateless, high-performance Go runtime that compiles raw Markdown files into visually polarized technical templates, dynamically injecting cryptographic integrity signatures on every request loop.

## 1. The Philosophy: Red Sovereignty vs. Blue Illusion

The internet is broken not from a lack of information, but from a fundamental crisis of architecture. Well-meaning idealists continually attempt to fix this by proposing centralized chokepoints like B.L.U.E. (Broad Learning Universal Education)—a theoretical centralized archive designed to hold all human "how-to" knowledge in a single, strictly moderated, unified database.

While noble in its intent to provide universal free education, the BLUE System suffers from structural flaws that guarantee its immediate compromise or destruction:

- **The Single Point of Failure:** By attempting to gather all technical knowledge onto a single domain, the Blue System creates an existential target for corporate lawfare and state suppression.
- **The Conflict of Interest:** The moment a platform's revenue is derived from a percentage of material sales via affiliate links, its neutrality is permanently compromised.
- **The Jury-Style Delusion:** Relying on an uncompensated "jury" of randomly selected expert verifiers to audit complex technical workflows is an economic fantasy.
- **The Dependency Hell of Linear Hierarchies:** Forcing all human knowledge into a rigid, linear hierarchy where Level 4 strictly depends on Level 1 creates a catastrophic maintenance bottleneck.

### The R.E.D Correction: Rejecting the Centralized Silo

Project R.E.D. explicitly rejects the concept of a master website. More importantly, **we reject the reliance on centralized corporate forums** as permanent curation layers. 

Knowledge cannot be trusted to a boardroom or a single database. It must be scattered so thoroughly across the crust of the earth that it cannot be stamped out. R.E.D. replaces the fragile, centralized tower with an indestructible, decentralized mesh of sovereign nodes and cryptographic trust networks.

## 2. Architecture: The Dual-Tier Deployment

Project R.E.D operates via a highly secure, dual-tier Docker network topology. It separates public "Clearnet" knowledge from restricted "Darknet" knowledge using strict Docker network isolation.

```text
GLOBAL INTERNET & TOR NETWORK
       │                   │
  [Port 80/443]       [Tor Network Bridge]
       │                   │
+------▼-------+    +------▼------------------+
| Caddy Proxy  |    | Tor Sidecar Container   |
+------┬-------+    +------┬---------┬--------+
       │                   │         │
[Clearnet-Tier Network]    │    [Onion-Tier Network (Internal)]
       │                   │         │
+------▼-------+           │  +------▼-------+
|  Light Node  |           │  |  Dark Node   |
| (Go Engine)  |           │  | (Go Engine)  |
+------┬-------+           │  +------┬-------+
       │                   │         │
       +----------+        │         |
                  │        │         │
               +--▼--------▼---------▼--+
               |   Host Volume (/data)  |
               | - /public              |
               | - /restricted          |
               +------------------------+
````

### Key Security Mechanisms:

- **Stateless Web Routing:** No database layer (SQL/NoSQL) exists to be corrupted, breached, or subpoenaed.
    
- **Cryptographic Validation:** The engine computes a standard SHA-256 hash over the raw file bytes on demand, injecting it into response headers (`X-RED-Content-Hash`).
    
- **Air-Gapped Vault:** The `Dark Node` lives on a strictly internal Docker network. It has zero outbound internet access. The `Tor Sidecar` acts as the exclusive cryptographic bridge, publishing the hidden service descriptor to the global Tor network while routing incoming requests securely to the Vault.
    

## 3. Technical Stack Blueprint

- **Language & Runtime:** Go (Golang) — Compiles down to a single memory-safe static binary with zero virtual machine overhead.
    
- **Markdown Parser:** `goldmark` — Highly efficient, fully CommonMark compliant.
    
- **Container Environment:** Multi-stage minimal Docker pipeline (`alpine`) for complete process isolation and instant replication.
    

## 4. Local Quickstart (From Scratch)

### Prerequisites

- Go 1.26+ installed locally (for native execution)
    
- Docker Desktop installed and running
    

### Clone & Initialize Project Workspace

Bash

```
# Set up directory structural blocks
mkdir -p templates static data/public data/restricted

# Initialize module tracking and download dependencies
go mod init red-engine
go get [github.com/yuin/goldmark](https://github.com/yuin/goldmark)
go get [github.com/adrg/frontmatter](https://github.com/adrg/frontmatter)
```

Execute the engine natively for local testing:

Bash

```
go run main.go
```

Navigate to `http://localhost:8080/guides/...` on your browser to test the local deployment.

## 5. Production Deployment Pipeline

Project R.E.D is designed to be deployed using `docker-compose`, which spins up both the Clearnet Gateway and the Onion Vault simultaneously.

### 1. Structure Your Data

Because the nodes operate on different access tiers, you must divide your Markdown files into the respective subdirectories on your host machine. The containers will only look in their assigned folders:

- **Clearnet Node:** Place files in `./data/public/`
    
- **Dark Node:** Place files in `./data/restricted/`
    

### 2. Launch the Matrix

Bring up the entire stack in detached mode:

Bash

```
docker-compose up -d --build
```

### 3. Retrieve Your Onion Address

The Tor sidecar will automatically generate cryptographic keys and a `.onion` address upon its first boot. To find out where your Vault is being hosted, run:

Bash

```
cat ./tor_keys/hostname
```

_(Note: It may take 5-10 minutes for a newly generated V3 hidden service to propagate fully across the global Tor directory.)_

## 6. Simulating Knowledge Base Structures

Rather than enforcing a brittle, top-down linear hierarchy that breaks under dependency hell, R.E.D. leverages your computer's native filesystem directory tracking. Group, version, and fork your files within folders inside your local storage volume dynamically:

Plaintext

```
data/
└── restricted/
    └── solar-array-build/
        ├── 00-index.md
        ├── 01-wiring-schematics-v1.0.md
        └── 02-inverter-fusing.md
```

The node pathing mechanics resolve these automatically into accessible downstream routes:

- `http://<your-onion>.onion/guides/solar-array-build/00-index`
    

## 7. Curation Philosophy: The Sovereign Web-of-Trust

1. **The Software Only Knows State:** The Go runtime handles zero moderation logic. It treats strings objectively.
    
2. **Cryptographic Peer Review:** Instead of relying on a centralized platform's voting buttons, users and curators sign content hashes with their independent cryptographic keys (PGP/Nostr).
    
3. **End-User Agency:** You, the reader, choose which curators to trust. Your local index client aggregates links signed by your trusted network.
    

The creator of the Blue System lamented that he lacked the millions of dollars needed for a startup, resigning himself to the belief that his vision would never be a reality. He failed because he thought like a corporate founder seeking capital.

We do not need capital. We do not need permission. We do not need a master website.

It's about time we stop waiting for someone or something to fix our problems for us. No hero is coming.

**Claim Your Agency. Run a Node.**
