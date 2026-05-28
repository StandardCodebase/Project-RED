package router

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *handler) webhookSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// We launch this in a goroutine so we can immediately return a 200 OK
	// to the Git server, preventing the webhook from timing out.
	go func() {
		log.Println("🔄 Webhook received! Triggering global knowledge base sync...")

		// TODO: In Phase 2, we will loop through h.cfg.StartupSync
		// and run the native git pull commands here.

		// Hot-reload the memory map after sync
		if err := h.store.Reload(); err != nil {
			log.Printf("⚠️ Webhook sync completed, but memory reload failed: %v", err)
		} else {
			log.Println("✅ Webhook sync complete. Memory index updated.")
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sync process initiated"))
}

// SearchItem represents a single searchable entry in the client-side index
type SearchItem struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"site":   h.cfg.SiteName,
	})
}

func (h *handler) manifest(w http.ResponseWriter, r *http.Request) {
	m := map[string]map[string]string{}

	for _, sec := range h.store.Root() {
		for _, a := range sec.Articles {
			m[a.Path] = map[string]string{"title": a.Title}
		}
		for _, sub := range sec.Sub {
			for _, a := range sub.Articles {
				m[a.Path] = map[string]string{"title": a.Title}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// searchIndex provides a flat list of articles for real-time frontend filtering
func (h *handler) searchIndex(w http.ResponseWriter, r *http.Request) {
	var index []SearchItem

	// Traverse the memory tree and build a flat array
	for _, sec := range h.store.Root() {
		for _, a := range sec.Articles {
			index = append(index, SearchItem{Title: a.Title, Path: a.Path})
		}
		for _, sub := range sec.Sub {
			for _, a := range sub.Articles {
				index = append(index, SearchItem{Title: a.Title, Path: a.Path})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(index)
}
