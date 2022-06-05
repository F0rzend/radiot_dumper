package copier

import (
	"os"
	"time"
)

type FileBuilder interface {
	CreateFile(ext string) (*os.File, error)
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
	return d.copier.ListenAndCopy(path, d.fileBuilder, d.timeout)
}
