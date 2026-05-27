package router

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"unicode"
    "unicode/utf8"

	"github.com/RED-Collective/red-engine/internal/render"
	"github.com/RED-Collective/red-engine/internal/store"
)

//go:embed templates static
var files embed.FS

type handler struct {
	store    *store.Store
	tmpl     *template.Template
	siteName string
}

func New(s *store.Store, siteName string) http.Handler {
	tmpl := template.Must(template.ParseFS(files, "templates/base.html"))

	staticFS, err := fs.Sub(files, "static")
	if err != nil {
	    panic(err)
	}

	h := &handler{store: s, tmpl: tmpl, siteName: siteName}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	mux.HandleFunc("/", h.serve)
	mux.HandleFunc("/-/reload", h.reload)
	return mux
}

type crumb struct {
	Label string
	Path  string
}

type pageData struct {
	Site   string
	Nav    map[string]*store.Section
	Body   template.HTML
	Title  string
	Path   string
	TopCat string
	Crumb  []crumb
}

func (h *handler) serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	topCat := ""
	if parts[0] != "" {
		topCat = parts[0]
	}

	d := pageData{
		Site:   h.siteName,
		Nav:    h.store.Root(),
		Path:   path,
		TopCat: topCat,
	}

	switch {
	case path == "/":
		d.Body = template.HTML(`<div class="article"><h1>` + h.siteName + `</h1><p>The free practical knowledge base. Choose a topic from the sidebar.</p></div>`)

	case len(parts) == 1 && topCat != "":
		sec, ok := h.store.Root()[topCat]
		if !ok {
			http.NotFound(w, r)
			return
		}
		d.Title = cap(topCat)
		d.Crumb = []crumb{{Label: cap(topCat), Path: "/" + topCat}}
		d.Body = template.HTML(sectionHTML(sec))

	default:
		raw, ok := h.store.Resolve(path)
		if !ok {
			http.NotFound(w, r)
			return
		}
		out, err := render.Markdown(raw)
		if err != nil {
			http.Error(w, "render error", 500)
			return
		}
		d.Title = cap(parts[len(parts)-1])
		d.Crumb = buildCrumbs(parts)
		d.Body = template.HTML(`<div class="article">` + out + `</div>`)
	}

	if err := h.tmpl.Execute(w, d); err != nil {
    	http.Error(w, "template error", 500)
    	return
	}
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
	w.WriteHeader(204)
}

func sectionHTML(sec *store.Section) string {
	var b strings.Builder
	b.WriteString(`<div class="section-index"><h1>` + cap(sec.Name) + `</h1>`)
	// open/close ul once around all articles
	b.WriteString(`<ul>`)
	for _, a := range sec.Articles {
	    b.WriteString(`<li><a href="` + a.Path + `">` + a.Title + `</a></li>`)
	}
	b.WriteString(`</ul>`)
	for _, sub := range sec.Sub {
		b.WriteString(`<h2>` + cap(sub.Name) + `</h2><ul>`)
		for _, a := range sub.Articles {
			b.WriteString(`<li><a href="` + a.Path + `">` + a.Title + `</a></li>`)
		}
		b.WriteString(`</ul>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func buildCrumbs(parts []string) []crumb {
	crumbs := make([]crumb, 0, len(parts))
	path := ""
	for _, p := range parts {
		path += "/" + p
		crumbs = append(crumbs, crumb{Label: cap(p), Path: path})
	}
	return crumbs
}

func cap(s string) string {
    s = strings.ReplaceAll(s, "-", " ")
    s = strings.ReplaceAll(s, "_", " ")
    if s == "" {
        return s
    }
    r, size := utf8.DecodeRuneInString(s)
    return string(unicode.ToUpper(r)) + s[size:]
}
