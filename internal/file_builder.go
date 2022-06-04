package internal

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultDateFormat = "02_01_2006"
)

type DatedFileBuilder struct {
	opts DatedFileOptions
}

func NewDatedFileBuilder(opts DatedFileOptions) *DatedFileBuilder {
	if opts.DateFormat == "" {
		opts.DateFormat = defaultDateFormat
	}

	return &DatedFileBuilder{
		opts: opts,
	}
}

type DatedFileOptions struct {
	Prefix     string
	DateFormat string
	Extension  string
}

func GetFileName(opts DatedFileOptions) string {
	return fmt.Sprintf(
		"%s%s.%s",
		opts.Prefix,
		time.Now().Format(opts.DateFormat),
		opts.Extension,
	)
}

func (f *DatedFileBuilder) CreateFile() (*os.File, error) {
	filename := GetFileName(f.opts)
	return os.Create(filename)
}
