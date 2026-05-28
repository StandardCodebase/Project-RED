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

	// --- NEW: Security Token Warning ---
	if cfg.AdminToken == "" || cfg.AdminToken == "secret123" {
		log.Println("=================================================================")
		log.Println("⚠️  SECURITY WARNING: Using default or missing Admin Token!    ⚠️")
		log.Println("⚠️  Anyone on the internet can overwrite your markdown files!  ⚠️")
		log.Println("")
		log.Println("CHANGE YOUR TOKEN BY RUNNING:")
		log.Println("Linux/macOS: ./manage-token.sh")
		log.Println("Windows:     .\\manage-token.ps1")
		log.Println("=================================================================")
	}
	// -----------------------------------

	// 1. Core Knowledge Base Pulling
	if *pull && cfg.SourceURL != "" {
		log.Printf("pulling core knowledge base from %s", cfg.SourceURL)
		if err := fetch.Pull(cfg.SourceURL, cfg.SourceType, cfg.DataDir); err != nil {
			log.Fatalf("fetch: %v", err)
		}
		log.Println("fetch complete")
	}

	// 2. Startup Sync (Ported from Legacy Gateway)
	if len(cfg.StartupSync) > 0 {
		remoteDir := filepath.Join(cfg.DataDir, "remote")

		// Check for permission errors right here before proceeding
		if err := os.MkdirAll(remoteDir, 0755); err != nil {
			log.Fatalf("CRITICAL: Failed to create remote directory. Check volume permissions: %v", err)
		}

		client := &http.Client{Timeout: 15 * time.Second}

		for _, sync := range cfg.StartupSync {
			log.Printf("Startup Sync: Fetching %s...", sync.Filename)
			if err := executeSync(client, sync.URL, filepath.Join(remoteDir, sync.Filename)); err != nil {
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

	// 4. Start HTTP Server with the Refactored Router
	h := router.New(s, &cfg, *cfgPath)
	log.Printf("RED listening on %s", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, h))
}

func executeSync(client *http.Client, targetURL, destSubPath string) error {
	// Reconstruct target file paths relative to data root directory
	// Note: main.go has context of cfg.DataDir

	// Let's resolve the path correctly depending on the initialization parameters
	lowerURL := strings.ToLower(targetURL)

	if strings.HasSuffix(lowerURL, ".tar.gz") || strings.HasSuffix(lowerURL, ".zip") {
		srcType := "tar.gz"
		if strings.HasSuffix(lowerURL, ".zip") {
			srcType = "zip"
		}
		// Pull the dynamic folder contents using the internal archive worker
		return fetch.Pull(targetURL, srcType, filepath.Join("data", destSubPath))
	}

	// Otherwise, proceed with single file retrieval flow
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

	fullFilePath := filepath.Join("data", destSubPath)
	if err := os.MkdirAll(filepath.Dir(fullFilePath), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(fullFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
