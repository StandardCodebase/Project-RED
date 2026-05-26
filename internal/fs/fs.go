package fs

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

func SecureJoin(baseDir, requestedPath string) (string, error) {
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

func ReadFileWithContext(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	
	return os.ReadFile(path)
}
