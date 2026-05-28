package router

import (
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
