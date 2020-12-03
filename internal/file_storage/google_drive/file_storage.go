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
			log.Trace(),
			log.Action("d.Files.Copy"),
			log.Reason(err.Error()),
			log.NewValue("fileID", fileID),
			log.NewValue("title", title),
		)
		return nil, err
	}
	return r, nil
}

// CreatePostMortemDocument creates a document on Google Drive from the PostMortem template.
func (s *storage) CreatePostMortemDocument(ctx context.Context, postMortemName string) (string, error) {
	s.logger.Info(
		ctx,
		log.Trace(),
		log.NewValue("postMortemName", postMortemName),
	)

	driveTokenBytes := []byte(config.Env.GoogleDriveToken)

	gClient, err := googleauth.Struct.GetGClient(ctx, s.logger, driveTokenBytes, drive.DriveScope)
	if err != nil {
		s.logger.Error(
			ctx,
			log.Trace(),
			log.Action("googleauth.Struct.GetGClient"),
			log.Reason(err.Error()),
			log.NewValue("postMortemName", postMortemName),
		)

		return "", err
	}

	driveService, err := drive.New(gClient)
	if err != nil {
		s.logger.Error(
			ctx,
			log.Trace(),
			log.Action("drive.New"),
			log.Reason(err.Error()),
			log.NewValue("postMortemName", postMortemName),
		)

		return "", err
	}

	file, err := s.copyFile(ctx, driveService, config.Env.GoogleDriveFileID, postMortemName)
	if err != nil {
		s.logger.Error(
			ctx,
			log.Trace(),
			log.Action("s.copyFile"),
			log.Reason(err.Error()),
			log.NewValue("postMortemName", postMortemName),
		)

		return "", err
	}

	if file != nil {
		s.logger.Info(
			ctx,
			log.Trace(),
			log.NewValue("postMortemName", postMortemName),
			log.NewValue("file", "https://docs.google.com/document/d/"+file.Id+"/edit"),
		)
	}

	return "https://docs.google.com/document/d/" + file.Id + "/edit", nil
}
