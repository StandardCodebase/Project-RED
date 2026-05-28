package router

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/RED-Collective/red-engine/internal/fetch"
)

func (h *handler) webhookSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("🔄 Webhook received from GitHub. Syncing tracked repositories...")

	// We launch this in a goroutine so we can immediately return a 200 OK
	// to GitHub, preventing the webhook from timing out.
	go func() {
		successCount := 0

		// 1. Loop through all tracked repositories in the configuration
		for _, sync := range h.cfg.StartupSync {

			// Determine the source type dynamically just like the importer does
			lowerURL := strings.ToLower(sync.URL)
			srcType := "raw"
			if strings.HasSuffix(lowerURL, ".git") {
				srcType = "git"
			} else if strings.HasSuffix(lowerURL, ".tar.gz") {
				srcType = "tar.gz"
			} else if strings.HasSuffix(lowerURL, ".zip") {
				srcType = "zip"
			}

			// Destination directly inside the data directory
			destDir := filepath.Join(h.store.DataDir(), sync.Filename)

			log.Printf("📥 Webhook triggering network pull for: %s", sync.Filename)

			// Execute the actual download / git mirror operation
			if err := fetch.Pull(sync.URL, srcType, destDir); err != nil {
				log.Printf("⚠️ Failed to sync %s: %v", sync.Filename, err)
			} else {
				successCount++
			}
		}

		if successCount > 0 {
			// 2. Hot-reload the memory map AFTER the files are successfully updated on disk
			if err := h.store.Reload(); err != nil {
				log.Printf("⚠️ Webhook sync completed, but memory index reload failed: %v", err)
			} else {
				log.Println("✅ Webhook sync complete. Memory index updated.")
			}
		} else {
			log.Println("⚠️ Webhook finished, but no tracked repositories were successfully synced.")
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sync process initiated"))
}
