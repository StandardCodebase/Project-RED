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

	"github.com/RED-Collective/red-engine/internal/models"
	"github.com/RED-Collective/red-engine/internal/render"
	"github.com/fsnotify/fsnotify"
)

type Store struct {
	dataDir string
	nav     map[string]*models.Section
	mu      sync.RWMutex
}

func New(dataDir string) *Store {
	return &Store{
		dataDir: dataDir,
		nav:     make(map[string]*models.Section),
	}
}

func (s *Store) DataDir() string {
	return s.dataDir
}

func parseManifestJSON(data []byte) map[string]models.ManifestEntry {
	result := make(map[string]models.ManifestEntry)
	var wrapped models.Manifest
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Files) > 0 {
		return wrapped.Files
	}
	if err := json.Unmarshal(data, &result); err == nil {
		return result
	}
	return result
}

func (s *Store) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	absDataDir, _ := filepath.Abs(s.dataDir)

	filepath.WalkDir(absDataDir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					log.Printf("🔄 File change detected: %s. Reloading store...", event.Name)
					s.Reload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("⚠️ Watcher error:", err)
			}
		}
	}()

	log.Printf("[DEBUG] File watcher started on %s", absDataDir)
	return nil
}

func (s *Store) Reload() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	trustedKeys := make(map[string]string)
	if trustData, err := os.ReadFile("contributors.json"); err == nil {
		var contributors []models.Contributor
		if err := json.Unmarshal(trustData, &contributors); err == nil {
			for _, c := range contributors {
				trustedKeys[strings.ToLower(c.PublicKey)] = c.Name
			}
		}
	} else {
		log.Println("⚠️  Warning: contributors.json not found. Verification checks disabled.")
	}

	allSignatures := make(map[string]models.ManifestEntry)
	filepath.WalkDir(s.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Base(path) != "manifest.json" {
			return nil
		}
		manifestData, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		manifest := parseManifestJSON(manifestData)
		if len(manifest) == 0 {
			return nil
		}
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

	newNav := make(map[string]*models.Section)

	err := filepath.WalkDir(s.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		hashBytes := sha256.Sum256(content)
		fileHash := hex.EncodeToString(hashBytes[:])

		res, err := render.Markdown(string(content))
		if err != nil {
			return nil
		}

		rel, _ := filepath.Rel(s.dataDir, path)
		relativePath := strings.TrimPrefix(filepath.ToSlash(rel), "/")

		cleanPath := strings.TrimSuffix(relativePath, ".md")

		isVerified := false
		authorName := "Unverified / Unknown Origin"
		verifyErr := "File signature not found in manifest"

		if entry, exists := allSignatures[relativePath]; exists {
			entryHash := entry.FileHash
			if entryHash == "" {
				entryHash = entry.Hash
			}

			if entryHash == fileHash {
				if trustedAuthor, isTrusted := trustedKeys[strings.ToLower(entry.PublicKey)]; isTrusted {
					pubBytes, err1 := hex.DecodeString(entry.PublicKey)
					sigBytes, err2 := hex.DecodeString(entry.Signature)

					if err1 == nil && err2 == nil && len(pubBytes) == ed25519.PublicKeySize {
						if ed25519.Verify(pubBytes, content, sigBytes) ||
							ed25519.Verify(pubBytes, []byte(fileHash), sigBytes) ||
							ed25519.Verify(pubBytes, hashBytes[:], sigBytes) {
							isVerified = true
							authorName = trustedAuthor
							verifyErr = ""
						} else {
							verifyErr = "Invalid Signature: Cryptographic verification failed"
						}
					} else {
						verifyErr = "Malformed Signature or Public Key data"
					}
				} else {
					verifyErr = "Untrusted Key: The public key is not mapped in contributors.json"
				}
			} else {
				verifyErr = "Hash Mismatch: File content was modified after signing"
			}
		}

		parts := strings.Split(filepath.ToSlash(cleanPath), "/")

		title := parts[len(parts)-1]
		title = strings.ReplaceAll(title, "-", " ")
		title = strings.Title(title)

		art := &models.Article{
			Path:              "/" + filepath.ToSlash(cleanPath),
			Title:             title,
			Body:              template.HTML(res.HTMLContent),
			Hash:              fileHash,
			Verified:          isVerified,
			Author:            authorName,
			VerificationError: verifyErr,
		}

		if len(parts) == 1 {
			if newNav["root"] == nil {
				newNav["root"] = &models.Section{Name: "root"}
			}
			newNav["root"].Articles = append(newNav["root"].Articles, art)
		} else {
			secName := parts[0]
			if newNav[secName] == nil {
				newNav[secName] = &models.Section{Name: secName, Sub: make(map[string]*models.Section)}
			}
			sec := newNav[secName]
			if len(parts) == 2 {
				sec.Articles = append(sec.Articles, art)
			} else {
				subName := parts[1]
				if sec.Sub[subName] == nil {
					sec.Sub[subName] = &models.Section{Name: subName}
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

func (s *Store) Nav() map[string]*models.Section {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nav
}

func (s *Store) Get(path string) *models.Article {
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

func (s *Store) Root() map[string]*models.Section {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copy := make(map[string]*models.Section, len(s.nav))
	for k, v := range s.nav {
		copy[k] = v
	}
	return copy
}
