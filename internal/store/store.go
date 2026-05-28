package store

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/RED-Collective/red-engine/internal/render"
)

type Article struct {
	Path     string
	Title    string
	Body     template.HTML
	Hash     string // SHA-256 Hash for UI Display
	Verified bool   // Ed25519 Verification Status
	Author   string // Name of the verified signer
}

type Section struct {
	Name     string
	Articles []*Article
	Sub      map[string]*Section
}

type Store struct {
	dataDir string
	nav     map[string]*Section
	mu      sync.RWMutex
}

func New(dataDir string) *Store {
	return &Store{
		dataDir: dataDir,
		nav:     make(map[string]*Section),
	}
}

func (s *Store) DataDir() string {
	return s.dataDir
}

// --- Cryptographic Structs matching the Obsidian Plugin ---
type Contributor struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type ManifestEntry struct {
	FileHash  string `json:"file_hash"`
	Hash      string `json:"hash"` // Fallback support
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type Manifest struct {
	Files map[string]ManifestEntry `json:"files"`
}

// Helper to unmarshal flexible manifest format
func parseManifestJSON(data []byte) map[string]ManifestEntry {
	result := make(map[string]ManifestEntry)

	// Try wrapped format first
	var wrapped Manifest
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Files) > 0 {
		return wrapped.Files
	}

	// Try flat format: {filepath: entry, ...}
	if err := json.Unmarshal(data, &result); err == nil {
		return result
	}

	return result
}

func (s *Store) Reload() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 1. Load the Trusted Public Keys from contributors.json
	trustedKeys := make(map[string]string)
	if trustData, err := os.ReadFile("contributors.json"); err == nil {
		var contributors []Contributor
		if err := json.Unmarshal(trustData, &contributors); err == nil {
			for _, c := range contributors {
				trustedKeys[strings.ToLower(c.PublicKey)] = c.Name
			}
		}
	} else {
		log.Println("⚠️  Warning: contributors.json not found. Verification checks disabled.")
	}

	// 2. Pre-load all signatures from any manifest.json in the data directory
	// Map by filepath for easier lookup
	allSignatures := make(map[string]ManifestEntry)
	filepath.WalkDir(s.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Base(path) != "manifest.json" {
			return nil
		}

		manifestData, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Warning: cannot read manifest %s: %v", path, err)
			return nil
		}

		manifest := parseManifestJSON(manifestData)
		if len(manifest) == 0 {
			return nil
		}

		// Get manifest directory relative to dataDir
		manifestDir := filepath.Dir(path)
		relManifestDir, err := filepath.Rel(s.dataDir, manifestDir)
		if err != nil {
			relManifestDir = "."
		}
		relManifestDir = filepath.ToSlash(relManifestDir)

		for key, entry := range manifest {
			key = filepath.ToSlash(key)
			var fullKey string
			if relManifestDir == "." {
				fullKey = key
			} else if strings.HasPrefix(key, relManifestDir+"/") || key == relManifestDir {
				fullKey = key
			} else {
				fullKey = filepath.ToSlash(filepath.Join(relManifestDir, key))
			}
			allSignatures[fullKey] = entry
		}
		return nil
	})

	log.Printf("[DEBUG] Loaded %d signature entries", len(allSignatures))
	for k := range allSignatures {
		log.Printf("[DEBUG]   %s", k)
	}

	newNav := make(map[string]*Section)

	err := filepath.WalkDir(s.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// 3. Calculate SHA-256
		hashBytes := sha256.Sum256(content)
		fileHash := hex.EncodeToString(hashBytes[:])

		res, err := render.Markdown(string(content))
		if err != nil {
			return nil
		}

		// 4. Ed25519 Cryptographic Verification
		isVerified := false
		authorName := "Unverified / Unknown Origin"

		// Get relative path in the format used in manifest
		rel, _ := filepath.Rel(s.dataDir, path)
		relativePath := strings.TrimPrefix(filepath.ToSlash(rel), "/")
		log.Printf("[DEBUG] Checking file: %s -> relativePath=%s", path, relativePath)

		// Try to find manifest entry by filepath
		if entry, exists := allSignatures[relativePath]; exists {
			// Use file_hash if available, otherwise hash
			entryHash := entry.FileHash
			if entryHash == "" {
				entryHash = entry.Hash
			}

			// Does the hash match?
			if entryHash == fileHash {
				// Does the signature belong to a trusted public key?
				if trustedAuthor, isTrusted := trustedKeys[strings.ToLower(entry.PublicKey)]; isTrusted {
					pubBytes, err1 := hex.DecodeString(entry.PublicKey)
					sigBytes, err2 := hex.DecodeString(entry.Signature)

					if err1 == nil && err2 == nil && len(pubBytes) == ed25519.PublicKeySize {
						// Check 1: Did the plugin sign the raw Markdown content?
						if ed25519.Verify(pubBytes, content, sigBytes) {
							isVerified = true
							authorName = trustedAuthor
							// Check 2: Did the plugin sign the Hex string of the SHA256 hash? (Very common)
						} else if ed25519.Verify(pubBytes, []byte(fileHash), sigBytes) {
							isVerified = true
							authorName = trustedAuthor
							// Check 3: Did the plugin sign the raw SHA256 bytes?
						} else if ed25519.Verify(pubBytes, hashBytes[:], sigBytes) {
							isVerified = true
							authorName = trustedAuthor
						} else {
							log.Printf("[DEBUG] Signature verification failed for %s (Tried Content, Hex Hash, and Byte Hash)", relativePath)
						}
					}
				} else {
					log.Printf("[DEBUG] Public key not trusted for %s: %s", relativePath, entry.PublicKey)
				}
			} else {
				log.Printf("[DEBUG] Hash mismatch for %s: stored=%s, actual=%s", relativePath, entryHash, fileHash)
			}
		}

		// 5. Build Article Structure
		parts := strings.Split(filepath.ToSlash(rel), "/")

		title := strings.TrimSuffix(parts[len(parts)-1], ".md")
		title = strings.ReplaceAll(title, "-", " ")
		title = strings.Title(title)

		art := &Article{
			Path:     "/" + filepath.ToSlash(rel),
			Title:    title,
			Body:     template.HTML(res.HTMLContent),
			Hash:     fileHash,
			Verified: isVerified,
			Author:   authorName,
		}

		// Tree Building
		if len(parts) == 1 {
			if newNav["root"] == nil {
				newNav["root"] = &Section{Name: "root"}
			}
			newNav["root"].Articles = append(newNav["root"].Articles, art)
		} else {
			secName := parts[0]
			if newNav[secName] == nil {
				newNav[secName] = &Section{Name: secName, Sub: make(map[string]*Section)}
			}
			sec := newNav[secName]
			if len(parts) == 2 {
				sec.Articles = append(sec.Articles, art)
			} else {
				subName := parts[1]
				if sec.Sub[subName] == nil {
					sec.Sub[subName] = &Section{Name: subName}
				}
				sec.Sub[subName].Articles = append(sec.Sub[subName].Articles, art)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	s.nav = newNav
	return nil
}

func (s *Store) Nav() map[string]*Section {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nav
}

func (s *Store) Get(path string) *Article {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 1 {
		if sec, ok := s.nav["root"]; ok {
			for _, a := range sec.Articles {
				if a.Path == "/"+path {
					return a
				}
			}
		}
	} else if len(parts) == 2 {
		if sec, ok := s.nav[parts[0]]; ok {
			for _, a := range sec.Articles {
				if a.Path == "/"+path {
					return a
				}
			}
		}
	} else if len(parts) == 3 {
		if sec, ok := s.nav[parts[0]]; ok {
			if sub, ok := sec.Sub[parts[1]]; ok {
				for _, a := range sub.Articles {
					if a.Path == "/"+path {
						return a
					}
				}
			}
		}
	}
	return nil
}

func (s *Store) Root() map[string]*Section {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make(map[string]*Section, len(s.nav))
	for k, v := range s.nav {
		copy[k] = v
	}
	return copy
}
