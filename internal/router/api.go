package router

import (
	"encoding/json"
	"net/http"
)

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
