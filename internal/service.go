package internal

import (
	"os"
	"time"
)

type FileBuilder interface {
	CreateFile() (*os.File, error)
}

type DumperService struct {
	copier      *StreamCopier
	fileBuilder FileBuilder
	timeout     time.Duration
}

func NewDumberService(
	copier *StreamCopier,
	fileBuilder FileBuilder,
	timeout time.Duration,
) *DumperService {
	return &DumperService{
		copier:      copier,
		fileBuilder: fileBuilder,
		timeout:     timeout,
	}
}

func (d *DumperService) ListenAndCopy(path string) error {
	file, err := d.fileBuilder.CreateFile()
	if err != nil {
		return err
	}

	defer logClosing(file)

	return d.copier.ListenAndCopy(path, file, d.timeout)
}
