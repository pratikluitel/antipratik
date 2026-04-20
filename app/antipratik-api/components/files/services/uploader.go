package services

import (
	"context"

	fileslogic "github.com/pratikluitel/antipratik/components/files/logic"
)

// UploaderService exposes file upload capabilities to other components.
// Inject this interface rather than importing files/logic directly.
type UploaderService interface {
	UploadPhotoFiles(ctx context.Context, postID string, files []fileslogic.FileInput) ([]fileslogic.PhotoImageResult, error)
	UploadMusicFiles(ctx context.Context, postID string, audioFile *fileslogic.FileInput, albumArtFile *fileslogic.FileInput) (fileslogic.MusicFilesResult, error)
	UploadThumbnail(ctx context.Context, postID string, suffix string, file fileslogic.FileInput) (fileslogic.ThumbnailResult, error)
}

type uploaderService struct {
	logic fileslogic.UploadLogic
}

// NewUploaderService returns an UploaderService backed by the given UploadLogic.
func NewUploaderService(l fileslogic.UploadLogic) UploaderService {
	return &uploaderService{logic: l}
}

func (s *uploaderService) UploadPhotoFiles(ctx context.Context, postID string, files []fileslogic.FileInput) ([]fileslogic.PhotoImageResult, error) {
	return s.logic.UploadPhotoFiles(ctx, postID, files)
}

func (s *uploaderService) UploadMusicFiles(ctx context.Context, postID string, audioFile *fileslogic.FileInput, albumArtFile *fileslogic.FileInput) (fileslogic.MusicFilesResult, error) {
	return s.logic.UploadMusicFiles(ctx, postID, audioFile, albumArtFile)
}

func (s *uploaderService) UploadThumbnail(ctx context.Context, postID string, suffix string, file fileslogic.FileInput) (fileslogic.ThumbnailResult, error) {
	return s.logic.UploadThumbnail(ctx, postID, suffix, file)
}
