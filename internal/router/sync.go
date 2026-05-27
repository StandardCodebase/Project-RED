package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RED-Collective/red-engine/internal/config"
	"github.com/RED-Collective/red-engine/internal/fetch"
)

func (h *handler) importRemote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedToken := r.Header.Get("X-Admin-Token")
	if h.cfg.AdminToken != "" && providedToken != h.cfg.AdminToken {
		http.Error(w, "Unauthorized: Invalid Admin Token", http.StatusUnauthorized)
		return
	}

	var req struct {
		URL           string `json:"url"`
		Filename      string `json:"filename"` // Used as the clean directory or path name
		SaveToStartup bool   `json:"saveToStartup"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// 1. SSRF Protection
	parsedURL, err := url.Parse(req.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		http.Error(w, "Invalid URL scheme", http.StatusBadRequest)
		return
	}
	host := strings.ToLower(parsedURL.Host)
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") || strings.HasPrefix(host, "10.") || strings.HasPrefix(host, "192.168.") {
		http.Error(w, "Local network imports are strictly forbidden", http.StatusForbidden)
		return
	}

	// 2. Directory & Path Sanitization
	targetSubPath := filepath.Clean(req.Filename)
	if targetSubPath == "." || targetSubPath == "" || strings.HasPrefix(targetSubPath, "..") || filepath.IsAbs(targetSubPath) {
		targetSubPath = filepath.Join("remote", "sync-"+time.Now().Format("20060102150405"))
	}

	// Formulate the destination path inside the engine's store directory
	destinationDir := filepath.Join(h.store.DataDir(), targetSubPath)

	// 3. Determine if it is a single file or a compressed multi-file bundle
	lowerURL := strings.ToLower(req.URL)
	if strings.HasSuffix(lowerURL, ".tar.gz") || strings.HasSuffix(lowerURL, ".zip") {
		srcType := "tar.gz"
		if strings.HasSuffix(lowerURL, ".zip") {
			srcType = "zip"
		}

		if err := fetch.Pull(req.URL, srcType, destinationDir); err != nil {
			http.Error(w, "Failed to pull and unpack remote directory: "+err.Error(), http.StatusBadGateway)
			return
		}
	} else {
		// It's a single raw markdown file! Save it cleanly to the specified directory path
		if err := os.MkdirAll(filepath.Dir(destinationDir), 0755); err != nil {
			http.Error(w, "Failed to create directory structure", http.StatusInternalServerError)
			return
		}

		// Ensure the file itself ends with .md
		if !strings.HasSuffix(strings.ToLower(destinationDir), ".md") {
			destinationDir += ".md"
			targetSubPath += ".md"
		}

		httpReq, err := http.NewRequest(http.MethodGet, req.URL, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		httpReq.Header.Set("User-Agent", "RED-Engine-Sync/1.0")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(httpReq)
		if err != nil {
			http.Error(w, "Failed to connect to remote server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Remote server returned non-OK status", http.StatusBadGateway)
			return
		}

		outFile, err := os.Create(destinationDir)
		if err != nil {
			http.Error(w, "Failed to create file on disk", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, resp.Body); err != nil {
			http.Error(w, "Failed to write content", http.StatusInternalServerError)
			return
		}
	}

	// 4. Hot-reload the memory map store index
	if err := h.store.Reload(); err != nil {
		http.Error(w, "Content updated but failed to update memory index", http.StatusInternalServerError)
		return
	}

	// 5. Persist to configuration state if requested
	if req.SaveToStartup {
		h.cfg.StartupSync = append(h.cfg.StartupSync, config.RemoteSync{
			URL:      req.URL,
			Filename: targetSubPath,
		})
		if err := h.cfg.Save(h.cfgPath); err != nil {
			http.Error(w, "Synced successfully, but config save failed", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synced to data/" + targetSubPath))
}

// --- SECURE DASHBOARD ENDPOINTS ---

func (h *handler) adminConfig(w http.ResponseWriter, r *http.Request) {
	providedToken := r.Header.Get("X-Admin-Token")
	if h.cfg.AdminToken != "" && providedToken != h.cfg.AdminToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.cfg.StartupSync)
}

func (h *handler) adminRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedToken := r.Header.Get("X-Admin-Token")
	if h.cfg.AdminToken != "" && providedToken != h.cfg.AdminToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// 1. Remove from the Config Array
	var newSync []config.RemoteSync
	for _, sync := range h.cfg.StartupSync {
		if sync.Filename != req.Filename {
			newSync = append(newSync, sync)
		}
	}
	h.cfg.StartupSync = newSync

	// 2. Save Config to Disk
	if err := h.cfg.Save(h.cfgPath); err != nil {
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	// 3. Delete the physical target directory or file safely
	safeName := filepath.Clean(req.Filename)
	if safeName != "." && safeName != "" && !strings.HasPrefix(safeName, "..") && !filepath.IsAbs(safeName) {
		fullRemovalPath := filepath.Join(h.store.DataDir(), safeName)
		os.RemoveAll(fullRemovalPath) // Cleans folders and files instantly
	}

	// 4. Hot-reload the engine so it vanishes from the UI mapping
	h.store.Reload()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully removed " + safeName))
}
