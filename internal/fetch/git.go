package fetch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// pullGit handles native cloning and aggressive hard-resetting of git repositories
func pullGit(url, destDir string) error {
	gitDir := filepath.Join(destDir, ".git")

	// 1. Check if the directory is already a cloned repository
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		log.Printf("📥 Native Git: Cloning fresh repository into %s...", destDir)

		// Ensure parent directories exist
		if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
			return err
		}

		cmd := exec.Command("git", "clone", url, destDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git clone failed: %v", err)
		}
		return nil
	}

	// 2. If it already exists, aggressively sync it to mirror the remote state
	log.Printf("🔄 Native Git: Syncing existing repository at %s...", destDir)

	// Fetch latest changes from remote without merging
	fetchCmd := exec.Command("git", "fetch", "origin")
	fetchCmd.Dir = destDir
	fetchCmd.Stdout = os.Stdout
	fetchCmd.Stderr = os.Stderr
	if err := fetchCmd.Run(); err != nil {
		return fmt.Errorf("git fetch failed: %v", err)
	}

	// Hard reset to the exact state of what was just fetched (handles reverts natively)
	resetCmd := exec.Command("git", "reset", "--hard", "FETCH_HEAD")
	resetCmd.Dir = destDir
	resetCmd.Stdout = os.Stdout
	resetCmd.Stderr = os.Stderr
	if err := resetCmd.Run(); err != nil {
		return fmt.Errorf("git reset failed: %v", err)
	}

	// Clean out any untracked files or directories (handles file deletions natively)
	cleanCmd := exec.Command("git", "clean", "-fd")
	cleanCmd.Dir = destDir
	cleanCmd.Stdout = os.Stdout
	cleanCmd.Stderr = os.Stderr
	if err := cleanCmd.Run(); err != nil {
		log.Printf("⚠️ Native Git: Warning during git clean: %v", err)
	}

	return nil
}
