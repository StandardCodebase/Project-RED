package store

import (
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func (s *Store) Watcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	absDataDir, _ := filepath.Abs(s.dataDir)

	// Watch the root and all subdirectories
	filepath.WalkDir(absDataDir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Trigger a reload if a file is written, created, removed, or renamed
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Printf("🔄 File change detected: %s. Reloading store...", event.Name)
					s.Reload()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("⚠️ Watcher error:", err)
			}
		}
	}()

	log.Printf("[DEBUG] File watcher started on %s", absDataDir)
	return nil
}
