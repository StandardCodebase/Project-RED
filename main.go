package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type PostMetadata struct {
	Title         string   `yaml:"title"`
	Authors       []string `yaml:"authors"`      // list of primary authors
	Contributors  []string `yaml:"contributors"` // optional additional contributors
	CreatedAt     string   `yaml:"created_at"`
	UpdatedAt     string   `yaml:"updated_at"`
	LastEditor    string   `yaml:"last_editor"`
	DiscussionHub string   `yaml:"discussion_hub"`
}

type PageData struct {
	PostMetadata // embedded – provides Title, Author, Contributors, etc.
	NodeName     string
	ContentHash  string
	ContentPath  string
	HTMLContent  template.HTML
}

// Configurable state
var (
	NodeName string
	DataDir  string
	Port     string

	indexTemplate *template.Template
)

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

// 1. Fetch from URL and automatically assign directory based on the URL path
func handleImportURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	targetURL := r.FormValue("url")
	if targetURL == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Automatically build directory structure: /data/domain.com/path/to/
	dirPath := filepath.Join(DataDir, parsed.Host, filepath.Dir(parsed.Path))
	fileName := filepath.Base(parsed.Path)
	if !strings.HasSuffix(fileName, ".md") {
		fileName += ".md" // Ensure it's saved as markdown
	}

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		http.Error(w, "Failed to create directories", http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(targetURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch content from URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	outPath := filepath.Join(dirPath, fileName)
	outFile, err := os.Create(outPath)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	io.Copy(outFile, resp.Body)
	w.Write([]byte("Successfully imported guide to: " + filepath.Join(parsed.Host, parsed.Path)))
}

// 2. Dynamically scan directories to populate the side panel
func handleManifest(w http.ResponseWriter, r *http.Request) {
	guides := make(map[string]map[string]string)
	filepath.Walk(DataDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			relPath, _ := filepath.Rel(DataDir, path)
			relPath = strings.TrimSuffix(relPath, ".md")

			// Use the filename as the default title for the side panel
			guides[relPath] = map[string]string{"title": filepath.Base(relPath)}
		}
		return nil
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guides)
}

func main() {
	var err error
	indexTemplate, err = template.ParseFiles("templates/layout.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","node":"` + NodeName + `"}`))
	})

	http.HandleFunc("/guides/", handleRenderGuide)
	http.HandleFunc("/download/", handleDownloadGuide)
	http.HandleFunc("/import", handleImportURL)
	http.HandleFunc("/manifest", handleManifest)

	srv := &http.Server{
		Addr:         ":" + Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Project R.E.D. Engine [%s] initiating on port :%s...\n", NodeName, Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Engine panic: %v", err)
	}
}

func secureJoin(baseDir, requestedPath string) (string, error) {
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(filepath.Join(baseDir, requestedPath))
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(targetAbs, baseAbs+string(os.PathSeparator)) && targetAbs != baseAbs {
		return "", os.ErrPermission
	}
	return targetAbs, nil
}

func readFileWithContext(ctx context.Context, path string) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 1)

	go func() {
		data, err := os.ReadFile(path)
		ch <- result{data, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-ch:
		return r.data, r.err
	}
}

func handleRenderGuide(w http.ResponseWriter, r *http.Request) {
	httpError := func(status int) {
		http.Error(w, http.StatusText(status), status)
	}

	requestedPath := strings.TrimPrefix(r.URL.Path, "/guides/")
	if requestedPath == "" {
		httpError(http.StatusNotFound)
		return
	}

	cleanedPath := filepath.Clean(requestedPath)
	if cleanedPath == "." || strings.HasPrefix(cleanedPath, "..") || strings.Contains(cleanedPath, string(os.PathSeparator)+"..") {
		httpError(http.StatusForbidden)
		return
	}

	targetFile, err := secureJoin(DataDir, cleanedPath+".md")
	if err != nil {
		httpError(http.StatusForbidden)
		return
	}

	fileBytes, err := readFileWithContext(r.Context(), targetFile)
	if err != nil {
		if os.IsNotExist(err) {
			httpError(http.StatusNotFound)
		} else {
			httpError(http.StatusInternalServerError)
		}
		return
	}

	hasher := sha256.New()
	hasher.Write(fileBytes)
	hashString := hex.EncodeToString(hasher.Sum(nil))

	var meta PostMetadata
	markdownRaw, err := frontmatter.Parse(strings.NewReader(string(fileBytes)), &meta)
	if err != nil {
		httpError(http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	if err := goldmark.Convert(markdownRaw, &buf); err != nil {
		httpError(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-RED-Content-Hash", hashString)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	data := PageData{
		PostMetadata: meta,
		NodeName:     NodeName,
		ContentHash:  hashString,
		ContentPath:  cleanedPath,
		HTMLContent:  template.HTML(buf.String()),
	}

	if err := indexTemplate.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handleDownloadGuide(w http.ResponseWriter, r *http.Request) {
	httpError := func(status int) {
		http.Error(w, http.StatusText(status), status)
	}

	requestedPath := strings.TrimPrefix(r.URL.Path, "/download/")
	if requestedPath == "" {
		httpError(http.StatusNotFound)
		return
	}

	cleanedPath := filepath.Clean(requestedPath)
	if cleanedPath == "." || strings.HasPrefix(cleanedPath, "..") || strings.Contains(cleanedPath, string(os.PathSeparator)+"..") {
		httpError(http.StatusForbidden)
		return
	}

	targetFile, err := secureJoin(DataDir, cleanedPath+".md")
	if err != nil {
		httpError(http.StatusForbidden)
		return
	}

	fileBytes, err := readFileWithContext(r.Context(), targetFile)
	if err != nil {
		if os.IsNotExist(err) {
			httpError(http.StatusNotFound)
		} else {
			httpError(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filepath.Base(targetFile)}))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(fileBytes)
}
