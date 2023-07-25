package storage

import (
	"fmt"
	"io"
	"os"
)

type FileStorage struct {
	basePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{basePath}
}

func (fs *FileStorage) SaveFile(id string, src io.Reader) error {
	dst, err := os.Create(fmt.Sprintf("%s/%s.pdf", fs.basePath, id))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	return nil
}
