package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RED-Collective/red-engine/internal/config"
	"github.com/RED-Collective/red-engine/internal/fetch"
	"github.com/RED-Collective/red-engine/internal/router"
	"github.com/RED-Collective/red-engine/internal/store"
)

func main() {
	cfgPath := flag.String("config", "config.json", "path to config file")
	pull := flag.Bool("pull", false, "fetch knowledge base before starting")
	flag.Parse()

	cfg := config.Default()
	if _, err := os.Stat(*cfgPath); err == nil {
		loaded, err := config.Load(*cfgPath)
		if err != nil {
			log.Fatalf("config: %v", err)
		}
		cfg = loaded
	}

	if cfg.AdminToken == "" || cfg.AdminToken == "secret123" {
		log.Println("=================================================================")
		log.Println("⚠️  SECURITY WARNING: Using default or missing Admin Token!    ⚠️")
		log.Println("⚠️  Anyone on the internet can overwrite your markdown files!  ⚠️")
		log.Println("=================================================================")
	}

	if *pull && cfg.SourceURL != "" {
		if err := fetch.Pull(cfg.SourceURL, cfg.SourceType, cfg.DataDir); err != nil {
			log.Fatalf("fetch: %v", err)
		}
	}

	// 2. Startup Sync
	if len(cfg.StartupSync) > 0 {
		if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
			log.Fatalf("CRITICAL: Failed to create data directory: %v", err)
		}

		client := &http.Client{Timeout: 15 * time.Second}

		for _, sync := range cfg.StartupSync {
			// FIX: Directly join DataDir and Filename. NO remote folder!
			destinationPath := filepath.Join(cfg.DataDir, sync.Filename)

			if err := executeSync(client, sync.URL, destinationPath); err != nil {
				log.Printf("Startup Sync Error (%s): %v", sync.Filename, err)
			} else {
				log.Printf("Startup Sync: Successfully downloaded %s", sync.Filename)
			}
		}
	}

	// 3. Initialize Memory Store
	s := store.New(cfg.DataDir)
	if err := s.Reload(); err != nil {
		log.Fatalf("store: %v", err)
	}

	// 4. Start HTTP Server
	h := router.New(s, &cfg, *cfgPath)
	log.Printf("RED listening on %s", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, h))
}

func executeSync(client *http.Client, targetURL, destPath string) error {
	lowerURL := strings.ToLower(targetURL)

	// --- NATIVE GIT SUPPORT ---
	if strings.HasSuffix(lowerURL, ".git") {
		return fetch.Pull(targetURL, "git", destPath)
	}
	// --------------------------

	if strings.HasSuffix(lowerURL, ".tar.gz") || strings.HasSuffix(lowerURL, ".zip") {
		srcType := "tar.gz"
		if strings.HasSuffix(lowerURL, ".zip") {
			srcType = "zip"
		}
		// FIX: Use destPath directly. NO hardcoded "data" folder!
		return fetch.Pull(targetURL, srcType, destPath)
	}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "RED-Engine-Startup-Sync/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return os.ErrPermission
	}

	// FIX: Use destPath directly. NO hardcoded "data" folder!
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
