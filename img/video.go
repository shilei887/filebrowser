package img

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type VideoService struct {
	ffmpegPath string
}

func NewVideoService() *VideoService {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return &VideoService{ffmpegPath: ""}
	}
	return &VideoService{ffmpegPath: ffmpegPath}
}

func (s *VideoService) IsAvailable() bool {
	return s.ffmpegPath != ""
}

func (s *VideoService) GenerateThumbnail(ctx context.Context, videoPath string, width, height int, timeOffset string) ([]byte, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("ffmpeg is not available")
	}

	args := []string{
		"-ss", timeOffset,
		"-i", videoPath,
		"-vframes", "1",
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-q:v", "2",
		"-y",
		"-f", "image2pipe",
		"-",
	}

	cmd := exec.CommandContext(ctx, s.ffmpegPath, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return output, nil
}

func (s *VideoService) GenerateThumbnailFromReader(ctx context.Context, reader *bytes.Reader, width, height int, timeOffset string) ([]byte, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("ffmpeg is not available")
	}

	args := []string{
		"-ss", timeOffset,
		"-i", "pipe:0",
		"-vframes", "1",
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-q:v", "2",
		"-y",
		"-f", "image2pipe",
		"-",
	}

	cmd := exec.CommandContext(ctx, s.ffmpegPath, args...)
	cmd.Stdin = reader

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg error: %w, stderr: %s", err, stderr.String())
	}

	return output, nil
}

var videoExtensions = map[string]bool{
	".mp4":  true,
	".avi":  true,
	".mkv":  true,
	".mov":  true,
	".webm": true,
	".flv":  true,
	".wmv":  true,
	".m4v":  true,
	".3gp":  true,
	".ts":   true,
	".mts":  true,
	".m2ts": true,
	".vob":  true,
	".ogv":  true,
	".rm":   true,
	".rmvb": true,
	".asf":  true,
	".f4v":  true,
	".swf":  true,
}

func IsVideoExtension(ext string) bool {
	ext = strings.ToLower(ext)
	return videoExtensions[ext]
}