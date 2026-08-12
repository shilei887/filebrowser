package fbhttp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type VideoPathCache struct {
	cacheRoot string
	fsRoot    string
}

func NewVideoPathCache(cacheRoot, fsRoot string) *VideoPathCache {
	return &VideoPathCache{
		cacheRoot: cacheRoot,
		fsRoot:    fsRoot,
	}
}

func (v *VideoPathCache) targetPath(videoPath string, previewSize PreviewSize) string {
	relPath := strings.TrimPrefix(videoPath, v.fsRoot)
	relDir := filepath.Dir(relPath)
	baseName := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))
	
	sizeSuffix := "thumb"
	if previewSize == PreviewSizeBig {
		sizeSuffix = "big"
	}
	
	cacheDir := filepath.Join(v.cacheRoot, relDir)
	return filepath.Join(cacheDir, baseName+"-"+sizeSuffix+".jpg")
}

func (v *VideoPathCache) Store(_ context.Context, videoPath string, value []byte, previewSize PreviewSize) error {
	targetPath := v.targetPath(videoPath, previewSize)
	targetDir := filepath.Dir(targetPath)
	
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(targetPath, value, 0644)
}

func (v *VideoPathCache) Load(_ context.Context, videoPath string, previewSize PreviewSize) ([]byte, bool, error) {
	targetPath := v.targetPath(videoPath, previewSize)
	
	data, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	
	return data, true, nil
}

func (v *VideoPathCache) Delete(_ context.Context, videoPath string, previewSize PreviewSize) error {
	targetPath := v.targetPath(videoPath, previewSize)
	
	err := os.Remove(targetPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}