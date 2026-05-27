package router

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/RED-Collective/red-engine/internal/config"
	"github.com/RED-Collective/red-engine/internal/store"
)

//go:embed templates static
var files embed.FS

type handler struct {
	store     *store.Store
	tmpl      *template.Template
	adminTmpl *template.Template
	cfg       *config.Config
	cfgPath   string
}

func New(s *store.Store, cfg *config.Config, cfgPath string) http.Handler {
	tmpl := template.Must(template.ParseFS(files, "templates/base.html"))
	adminTmpl := template.Must(template.ParseFS(files, "templates/admin.html"))

	staticFS, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}

	h := &handler{
		store:     s,
		tmpl:      tmpl,
		adminTmpl: adminTmpl,
		cfg:       cfg,
		cfgPath:   cfgPath,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	mux.HandleFunc("/", h.serve)
	mux.HandleFunc("/-/reload", h.reload)
	mux.HandleFunc("/-/import", h.importRemote)
	mux.HandleFunc("/-/health", h.health)
	mux.HandleFunc("/-/manifest", h.manifest)
	mux.HandleFunc("/-/source/", h.source)
	mux.HandleFunc("/-/download/", h.download)
	// NEW: Secure Admin UI
	mux.HandleFunc("/-/admin", h.adminUI)

	return mux
}

func (h *handler) adminUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	h.adminTmpl.Execute(w, nil)
}

func (h *handler) reload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.store.Reload(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
