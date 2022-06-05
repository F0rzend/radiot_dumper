package copier

import (
	"fmt"
	"os"
	"time"
)

type DatedFileBuilder struct {
	opts DatedFileOptions
}

func NewDatedFileBuilder(opts DatedFileOptions) *DatedFileBuilder {
	return &DatedFileBuilder{
		opts: opts,
	}
}

type DatedFileOptions struct {
	OutputDirectory string
	Prefix          string
	DateFormat      string
}

func buildFileName(opts DatedFileOptions, index int, ext string) string {
	filePostfix := ""
	if index > 0 {
		filePostfix = fmt.Sprintf(".%d", index)
	}

	return opts.OutputDirectory + opts.Prefix + time.Now().Format(opts.DateFormat) + filePostfix + ext
}

func GetFileName(opts DatedFileOptions, ext string) string {
	for i := 0; ; i++ {
		filename := buildFileName(opts, i, ext)

		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return filename
		}
	}
}

func (f *DatedFileBuilder) CreateFile(ext string) (*os.File, error) {
	filename := GetFileName(f.opts, ext)
	return os.Create(filename)
}
