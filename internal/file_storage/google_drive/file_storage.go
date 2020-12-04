package googledrive

import (
	filestorage "hellper/internal/file_storage"
	googleauth "hellper/internal/google_auth"
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
			"googleDrive/google_drive.copyFile Copy ERROR",
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

	logWriter := s.logger.With(
		log.NewValue("postMortemName", postMortemName),
	)

	logWriter.Debug(
		ctx,
		"googleDrive/google_drive.CreatePostMortemDocument",
	)

	driveTokenBytes := []byte(config.Env.GoogleDriveToken)

	gClient, err := googleauth.Struct.GetGClient(ctx, s.logger, driveTokenBytes, drive.DriveScope)
	if err != nil {
		logWriter.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument GetGClient ERROR",
			log.NewValue("error", err),
		)

		return "", err
	}

	driveService, err := drive.New(gClient)
	if err != nil {
		logWriter.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument driveService ERROR",
			log.NewValue("error", err),
		)

		return "", err
	}

	file, err := s.copyFile(ctx, driveService, config.Env.GoogleDriveFileID, postMortemName)
	if err != nil {
		logWriter.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument copyFile ERROR",
			log.NewValue("error", err),
		)

		return "", err
	}

	if file != nil {
		logWriter.Debug(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument file",
			log.NewValue("file", "https://docs.google.com/document/d/"+file.Id+"/edit"),
		)
	}

	return "https://docs.google.com/document/d/" + file.Id + "/edit", nil
}
