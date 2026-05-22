package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type PostMetadata struct {
	Title          string `yaml:"title"`
	AuthorIdentity string `yaml:"author_identity"`
	CreatedAt      string `yaml:"created_at"`
	DiscussionHub  string `yaml:"discussion_hub"`
}

type PageData struct {
	PostMetadata
	NodeName    string
	ContentHash string
	HTMLContent template.HTML
}

// Configurable state variables mapped via environment
var NodeName string
var DataDir string
var Port string

func init() {
	NodeName = os.Getenv("RED_NODE_NAME")
	if NodeName == "" {
		NodeName = "Alpha-Centauri-01"
	}

	DataDir = os.Getenv("RED_DATA_DIR")
	if DataDir == "" {
		DataDir = "./data"
	}

	Port = os.Getenv("RED_PORT")
	if Port == "" {
		Port = "8080"
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/guides/", handleRenderGuide)

	log.Printf("Project R.E.D. Engine [%s] initiating on port :%s...\n", NodeName, Port)
	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("Engine panic: %v", err)
	}
}

func handleRenderGuide(w http.ResponseWriter, r *http.Request) {
	requestedPath := strings.TrimPrefix(r.URL.Path, "/guides/")
	cleanedPath := filepath.Clean(requestedPath)
	if strings.HasPrefix(cleanedPath, "..") || strings.HasPrefix(cleanedPath, "/") {
		http.Error(w, "Access Denied: Bounds violation.", http.StatusForbidden)
		return
	}

	targetFile := filepath.Join(DataDir, cleanedPath+".md")

	fileBytes, err := os.ReadFile(targetFile)
	if err != nil {
		http.Error(w, "Document not found inside local state volume.", http.StatusNotFound)
		return
	}

	hasher := sha256.New()
	hasher.Write(fileBytes)
	hashString := hex.EncodeToString(hasher.Sum(nil))

	var meta PostMetadata
	markdownRaw, err := frontmatter.Parse(strings.NewReader(string(fileBytes)), &meta)
	if err != nil {
		http.Error(w, "Malformed front-matter blueprint.", http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	if err := goldmark.Convert(markdownRaw, &buf); err != nil {
		http.Error(w, "Goldmark compilation loop exception.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-RED-Content-Hash", hashString)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		http.Error(w, "Template resolution compilation fault.", http.StatusInternalServerError)
		return
	}

	data := PageData{
		PostMetadata: meta,
		NodeName:     NodeName,
		ContentHash:  hashString,
		HTMLContent:  template.HTML(buf.String()),
	}

	tmpl.Execute(w, data)
}
