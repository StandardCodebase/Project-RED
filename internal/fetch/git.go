package fetch

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func pullGit(url, destDir string) ([]string, error) {
	repo, err := git.PlainOpen(destDir)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			log.Printf("📥 Native go-git: Repository missing or corrupted. Rebuilding %s...", destDir)

			// FIX: Destroy any leftover garbage folders from old zip syncs so git can clone cleanly
			os.RemoveAll(destDir)

			if err := os.MkdirAll(destDir, 0755); err != nil {
				return nil, err
			}
			_, err = git.PlainClone(destDir, false, &git.CloneOptions{
				URL:      url,
				Progress: os.Stdout,
			})
			if err != nil {
				return nil, fmt.Errorf("go-git clone failed: %v", err)
			}
			return nil, nil // Tells the router to do a full memory reload
		}
		return nil, fmt.Errorf("failed to check existing repository: %v", err)
	}

	log.Printf("🔄 Native go-git: Checking for delta updates at %s...", destDir)
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get git worktree: %v", err)
	}

	var oldHash plumbing.Hash
	if head, err := repo.Head(); err == nil {
		oldHash = head.Hash()
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Force:      true,
		Progress:   os.Stdout,
	})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Printf("✅ Sync skipped: %s is already up to date.", destDir)
			return []string{}, nil
		}
		return nil, fmt.Errorf("go-git delta pull failed: %v", err)
	}

	var changedFiles []string
	if head, err := repo.Head(); err == nil {
		newHash := head.Hash()
		if oldHash != plumbing.ZeroHash && oldHash != newHash {
			oldCommit, err1 := repo.CommitObject(oldHash)
			newCommit, err2 := repo.CommitObject(newHash)
			if err1 == nil && err2 == nil {
				patch, err3 := oldCommit.Patch(newCommit)
				if err3 == nil {
					for _, fileStat := range patch.Stats() {
						fullPath := filepath.Join(destDir, fileStat.Name)
						changedFiles = append(changedFiles, fullPath)
					}
				}
			}
		}
	}

	log.Printf("✅ Native go-git: Applied delta updates to %s (%d files changed)", destDir, len(changedFiles))
	return changedFiles, nil
}
