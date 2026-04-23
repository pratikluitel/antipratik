package services

import (
	"context"

	"github.com/pratikluitel/antipratik/components/files"
)

type uploaderService struct {
	logic files.UploadLogic
}

// NewUploaderService returns an UploaderService backed by the given UploadLogic.
func NewUploaderService(l files.UploadLogic) files.UploaderService {
	return &uploaderService{logic: l}
}

func (s *uploaderService) UploadPhotoFiles(ctx context.Context, postID string, files []files.FileInput) ([]files.PhotoImageResult, error) {
	return s.logic.UploadPhotoFiles(ctx, postID, files)
}

func (s *uploaderService) UploadMusicFiles(ctx context.Context, postID string, audioFile *files.FileInput, albumArtFile *files.FileInput) (files.MusicFilesResult, error) {
	return s.logic.UploadMusicFiles(ctx, postID, audioFile, albumArtFile)
}

func (s *uploaderService) UploadThumbnail(ctx context.Context, postID string, suffix string, file files.FileInput) (files.ThumbnailResult, error) {
	return s.logic.UploadThumbnail(ctx, postID, suffix, file)
}
