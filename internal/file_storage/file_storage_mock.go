package filestorage

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type FileStorageMock struct {
	mock.Mock
}

func NewFileStorageMock() *FileStorageMock {
	return new(FileStorageMock)
}

func (mock *FileStorageMock) CreatePostMortemDocument(ctx context.Context, postMortemName string) string {
	args := mock.Called(ctx, postMortemName)
	return args.Get(0).(string)
}
