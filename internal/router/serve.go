package router

import (
	"html/template"
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/RED-Collective/red-engine/internal/store"
)

type crumb struct {
	Label string
	Path  string
}

type pageData struct {
	Site     string
	Nav      map[string]*store.Section
	Body     template.HTML
	Title    string
	Path     string
	TopCat   string
	Crumb    []crumb
	Verified bool
	Author   string
	Hash     string
}

func (h *handler) serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	topCat := ""
	if parts[0] != "" {
		topCat = parts[0]
	}

	d := pageData{
		Site:   h.cfg.SiteName,
		Nav:    h.store.Root(),
		Path:   path,
		TopCat: topCat,
	}

	switch {
	case path == "/":
		d.Body = template.HTML(`<div class="article"><h1>` + h.cfg.SiteName + `</h1><p>The free practical knowledge base. Choose a topic from the sidebar.</p></div>`)

	case len(parts) == 1 && topCat != "":
		sec, ok := h.store.Root()[topCat]
		if !ok {
			http.NotFound(w, r)
			return
		}
		d.Title = capitalize(topCat)
		d.Crumb = []crumb{{Label: capitalize(topCat), Path: "/" + topCat}}
		d.Body = template.HTML(sectionHTML(sec))


	default:
    art := h.store.Get(path)
    if art == nil {
        http.NotFound(w, r)
        return
    }
    d.Title = capitalize(parts[len(parts)-1])
    d.Crumb = buildCrumbs(parts)
    d.Body = art.Body
    d.Verified = art.Verified   // from store.Article
    d.Author = art.Author
    d.Hash = art.Hash

	if err := h.tmpl.ExecuteTemplate(w, "base.html", d); err != nil {
		http.Error(w, "template error: "+err.Error(), 500)
		return
	}
}

func sectionHTML(sec *store.Section) string {

	var b strings.Builder
	b.Grow(1024)
	b.WriteString(`<div class="section-index"><h1>`)
	b.WriteString(capitalize(sec.Name))
	b.WriteString(`</h1><ul>`)

	for _, a := range sec.Articles {
		b.WriteString(`<li><a href="`)
		b.WriteString(a.Path)
		b.WriteString(`">`)
		b.WriteString(a.Title)
		b.WriteString(`</a></li>`)
	}
	b.WriteString(`</ul>`)

	for _, sub := range sec.Sub {
		b.WriteString(`<h2>`)
		b.WriteString(capitalize(sub.Name))
		b.WriteString(`</h2><ul>`)
		for _, a := range sub.Articles {
			b.WriteString(`<li><a href="`)
			b.WriteString(a.Path)
			b.WriteString(`">`)
			b.WriteString(a.Title)
			b.WriteString(`</a></li>`)
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
		crumbs = append(crumbs, crumb{Label: capitalize(p), Path: path})
	}
	return crumbs
}

func capitalize(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
