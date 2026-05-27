package router

import (
	"encoding/json"
	"net/http"
)

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"site":   h.siteName,
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
