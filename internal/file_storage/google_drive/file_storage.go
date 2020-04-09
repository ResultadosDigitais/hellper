package googledrive

import (
	"encoding/json"
	"net/http"

	filestorage "hellper/internal/file_storage"
	"hellper/internal/log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

// Retrieve a token, saves the token, then returns the generated client.
func (s *storage) getGClient(ctx context.Context, gConfig *oauth2.Config) *http.Client {
	b := []byte(config.Env.GoogleDriveToken)

	tok := &oauth2.Token{}
	err := json.Unmarshal(b, tok)
	if err != nil {
		s.logger.Error(
			ctx,
			"googleDrive/google_drive.getGClient Unmarshal error",
			log.NewValue("error", err),
		)
	}

	return gConfig.Client(ctx, tok)
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

	driveCredentialBytes := []byte(config.Env.GoogleDriveCredentials)

	gConfig, err := google.ConfigFromJSON(driveCredentialBytes, drive.DriveScope)
	if err != nil {
		s.logger.Error(
			ctx,
			"googleDrive/google_drive.CreatePostMortemDocument ConfigFromJSON error",
			log.NewValue("postMortemName", postMortemName),
			log.NewValue("error", err),
		)

		return "", err
	}

	gClient := s.getGClient(ctx, gConfig)

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
