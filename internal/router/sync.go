package router

import (
	"encoding/json"
	"net/http"

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
		SaveToStartup bool   `json:"saveToStartup"` // NEW: Frontend toggle
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// ... (Keep the existing SSRF, Filename sanitization, and Download code exactly as it is down to `io.Copy(outFile, resp.Body)`) ...

	// ... [Replace the bottom of the function with this:] ...

	// Reload the in-memory store index
	if err := h.store.Reload(); err != nil {
		http.Error(w, "File saved but failed to reload store", http.StatusInternalServerError)
		return
	}

	// NEW: Append to config.json if requested
	if req.SaveToStartup {
		h.cfg.StartupSync = append(h.cfg.StartupSync, config.RemoteSync{
			URL:      req.URL,
			Filename: safeName,
		})
		if err := h.cfg.Save(h.cfgPath); err != nil {
			http.Error(w, "Synced successfully, but failed to write to config.json", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully synced to remote/" + safeName))
}
