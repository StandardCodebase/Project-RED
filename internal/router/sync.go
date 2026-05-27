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
)

func (h *handler) importRemote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Enforce the Admin Token
	providedToken := r.Header.Get("X-Admin-Token")
	if h.cfg.AdminToken != "" && providedToken != h.cfg.AdminToken {
		http.Error(w, "Unauthorized: Invalid Admin Token", http.StatusUnauthorized)
		return
	}

	var req struct {
		URL           string `json:"url"`
		Filename      string `json:"filename"`
		SaveToStartup bool   `json:"saveToStartup"` // Frontend toggle
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// 2. SSRF Protection
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

	// 3. Filename sanitization
	safeName := filepath.Clean(req.Filename)
	if safeName == "." || safeName == "" || strings.Contains(safeName, "/") || strings.Contains(safeName, "\\") {
		safeName = "community-sync-" + time.Now().Format("20060102150405") + ".md"
	}
	if !strings.HasSuffix(safeName, ".md") {
		safeName += ".md"
	}

	// 4. Download the file
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

	// 5. Save to disk
	remoteDir := filepath.Join(h.store.DataDir(), "remote")
	if err := os.MkdirAll(remoteDir, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	outFile, err := os.Create(filepath.Join(remoteDir, safeName))
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		http.Error(w, "Failed to write file contents", http.StatusInternalServerError)
		return
	}

	// 6. Reload the in-memory store index
	if err := h.store.Reload(); err != nil {
		http.Error(w, "File saved but failed to reload store", http.StatusInternalServerError)
		return
	}

	// 7. Append to config.json if requested
	if req.SaveToStartup {
		h.cfg.StartupSync = append(h.cfg.StartupSync, config.RemoteSync{
			URL:      req.URL,
			Filename: safeName,
		})
		// Save the struct to disk
		if err := h.cfg.Save(h.cfgPath); err != nil {
			http.Error(w, "Synced successfully, but failed to write to config.json", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synced to remote/" + safeName))
}
