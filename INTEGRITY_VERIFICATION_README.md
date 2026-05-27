# RED Engine Integrity Verification Implementation

## Overview
This document outlines the integrity verification system integrated into RED Engine that validates file signatures and displays verification status to users.

## How It Works

### 1. **Manifest System**
- External tools (like the Obsidian RED Signer plugin) create `manifest.json` in vault directories
- Format: `{ "filepath": { "file_hash": "hex", "public_key": "hex", "signature": "hex" } }`
- Each entry signs the file's SHA-256 hash using Ed25519 cryptography

### 2. **Contributor Registry**
- `contributors.json` at the repo root lists trusted public keys
- Format:
```json
[
  { "name": "StandardCodebase", "public_key": "66f32e250b0bafb2683ae987e65a390901f030ccbc7745435d99f35da1bfccf1" }
]
```

### 3. **Verification Flow** (in `internal/store/store.go`)
1. **Load trusted keys** from `contributors.json` (by public_key hex string)
2. **Scan data directories** for `manifest.json` files
3. **For each markdown file**:
   - Calculate SHA-256 hash of file content
   - Look up matching entry in manifest by filename
   - Compare stored hash with calculated hash
   - Verify signature using Ed25519 public key
   - Check if public key is in trusted contributors list
4. **Set Article status**:
   - `Verified=true` + `Author=name` if signature valid and key trusted
   - `Verified=false` + `Author="Unverified / Unknown Origin"` otherwise

### 4. **UI Display** (in `internal/router/templates/base.html`)
- **Security Badge**: Prominently shows verification status at top of each article
- **Hash Footer**: Displays full SHA-256 checksum for manual verification
- **Visual Indicators**:
  - ✅ **Verified Contributor**: Green badge with contributor name
  - ⚠️ **Unverified / Unknown Origin**: Red badge for unsigned/untrusted files

## Code Changes

### Modified Files

#### 1. `internal/store/store.go`
**Key improvements:**
- Updated manifest structure to support flat format: `{filepath: entry}`
- Added `parseManifestJSON()` helper to handle both wrapped and flat formats
- Changed lookup from hash-based to filepath-based
- Fixed field names: `file_hash` (not just `hash`)
- Proper Ed25519 verification using `ed25519.Verify()`

**Important structs:**
```go
type ManifestEntry struct {
    FileHash  string `json:"file_hash"`
    PublicKey string `json:"public_key"`
    Signature string `json:"signature"`
}
```

#### 2. `internal/router/templates/base.html`
**Visual improvements:**
- Added emoji indicators (✅ / ⚠️) to badges
- Better spacing and visual hierarchy
- Emphasized "Verified Contributor: <name>" text
- Marked hash as "File integrity hash"

#### 3. `internal/router/static/article.css`
**Enhanced styling:**
- Gradient backgrounds for badges
- Better box shadows and borders
- Improved hash footer layout
- Added `user-select: all` for easy copying
- Responsive design maintained

## Testing

### Manual Test
1. Ensure `contributors.json` exists with trusted keys
2. Create `manifest.json` in vault with signed entries
3. Start RED Engine: `./red -config config.json`
4. Browse to any article
5. Check for:
   - Green ✅ badge if file is verified
   - Red ⚠️ badge if unverified
   - SHA-256 hash displayed at bottom

### Current Test Setup
- Vault: `data/remote/API/`
- Test files: `File.md`, `Future proofing.md`
- Contributors: `contributors.json` (root level)
- Manifest: `data/remote/API/manifest.json`

## Security Considerations

### ✅ What's Protected
- Files are verified using **Ed25519 signatures**
- Only files signed by **trusted contributors** are marked as verified
- Hash tampering is detected (hash mismatch → unverified)
- Untrusted signers are rejected even with valid signatures

### ⚠️ Limitations
- Does NOT prevent serving unverified files (UI shows warning instead)
- Does NOT auto-reject unsigned files (design choice: inform, don't block)
- Verification happens at **serve time** (not on write)
- No automatic signature generation (external tool responsibility)

## Integration Points

### For External Tools (Signers)
1. Generate Ed25519 keypair
2. Register public key in `contributors.json`
3. Create/update `manifest.json` with:
   - File path as key
   - SHA-256 hash of file content
   - Signature of content (Ed25519)
   - Public key reference

### For RED Engine Admins
1. Ensure `contributors.json` exists and is accessible
2. Point `dataDir` to vault containing `manifest.json`
3. Restart RED Engine to reload verification data
4. Monitor UI badges for verification status

## Future Improvements
- [ ] Add API endpoint `/api/verify` to check file integrity programmatically
- [ ] Implement file access control (reject unverified files if policy set)
- [ ] Add signature timestamp for audit trails
- [ ] Support multiple signatures per file
- [ ] Add admin panel to manage contributors

## Files Modified
- `internal/store/store.go` - Core verification logic
- `internal/router/templates/base.html` - Visual display
- `internal/router/static/article.css` - Enhanced styling
