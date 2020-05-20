package googledrive

import (
	filestorage "hellper/internal/file_storage"
	googleapi "hellper/internal/google_api"
	"hellper/internal/log"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"

	"hellper/internal/config"
)

type storage struct {
	logger log.Logger
}

// NewFileStorage initialize the file storage service
func NewFileStorage(logger log.Logger) filestorage.Driver {
	return &storage{
		logger: logger,
	}
}

func (s *storage) copyFile(ctx context.Context, d *drive.Service, fileID string, title string) (*drive.File, error) {
	f := &drive.File{Name: title}
	r, err := d.Files.Copy(fileID, f).Do()
	if err != nil {
		s.logger.Error(
			ctx,
			"googleDrive/google_drive.copyFile Copy error",
			log.NewValue("fileID", fileID),
			log.NewValue("title", title),
			log.NewValue("error", err),
		)
		return nil, err
	}
	return r, nil
}

// CreatePostMortemDocument creates a document on Google Drive from the PostMortem template.
func (s *storage) CreatePostMortemDocument(ctx context.Context, postMortemName string) (string, error) {
	s.logger.Info(
		ctx,
		"googleDrive/google_drive.CreatePostMortemDocument",
		log.NewValue("postMortemName", postMortemName),
	)

	driveTokenBytes := []byte(config.Env.GoogleDriveToken)

	gClient, err := googleapi.GoogleAuthStruct.GetGClient(ctx, s.logger, driveTokenBytes, drive.DriveScope)
	if err != nil {
		s.logger.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument GetGClient error",
			log.NewValue("postMortemName", postMortemName),
			log.NewValue("error", err),
		)

		return "", err
	}

	driveService, err := drive.New(gClient)
	if err != nil {
		s.logger.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument driveService error",
			log.NewValue("postMortemName", postMortemName),
			log.NewValue("error", err),
		)

		return "", err
	}

	file, err := s.copyFile(ctx, driveService, config.Env.GoogleDriveFileId, postMortemName)
	if err != nil {
		s.logger.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument copyFile error",
			log.NewValue("postMortemName", postMortemName),
			log.NewValue("error", err),
		)

		return "", err
	}

	if file != nil {
		s.logger.Info(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument file",
			log.NewValue("postMortemName", postMortemName),
			log.NewValue("file", "https://docs.google.com/document/d/"+file.Id+"/edit"),
		)
	}

	return "https://docs.google.com/document/d/" + file.Id + "/edit", nil
}
