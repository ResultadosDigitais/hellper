package filestorage

import (
	"context"
)

// Driver interface for File Storage
type Driver interface {
	CreatePostMortemDocument(ctx context.Context, postMortemName string) (string, error)
}
