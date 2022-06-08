package copier

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultDateLayout = "02_01_2006"
)

type DatedFileBuilder struct {
	opts            DatedFileOptions
	outputDirectory string
}

func NewDatedFileBuilder(
	outputDirectory string,
	outputFS fs.FS,
	prefix string,
	dateFormat string,
) *DatedFileBuilder {
	if dateFormat == "" {
		dateFormat = defaultDateLayout
	}

	return &DatedFileBuilder{
		opts: DatedFileOptions{
			outputFS:   outputFS,
			prefix:     prefix,
			dateFormat: dateFormat,
		},
		outputDirectory: outputDirectory,
	}
}

type DatedFileOptions struct {
	outputFS   fs.FS
	prefix     string
	dateFormat string
}

func (f *DatedFileBuilder) buildFileName(index int, ext string) string {
	filePostfix := ""
	if index > 0 {
		filePostfix = fmt.Sprintf(".%d", index)
	}

	filename := f.opts.prefix + time.Now().Format(f.opts.dateFormat) + filePostfix + ext
	_, err := fs.Stat(f.opts.outputFS, filename)
	if os.IsNotExist(err) {
		return filename
	}
	return f.buildFileName(index+1, ext)
}

func (f *DatedFileBuilder) GetFileName(ext string) string {
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return filepath.Join(f.outputDirectory, f.buildFileName(0, ext))
}

func (f *DatedFileBuilder) GetOutput(ext string) (io.WriteCloser, error) {
	return os.Create(f.GetFileName(ext))
}
