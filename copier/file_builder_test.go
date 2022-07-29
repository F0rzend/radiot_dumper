package copier

import (
	"github.com/stretchr/testify/assert"
	"io/fs"
	"mime"
	"testing"
	"testing/fstest"
	"time"
)

func TestDatedFileBuilder_getFileName(t *testing.T) {
	t.Parallel()

	layout := "02_01_2006"
	today := time.Now().Format(layout)

	outputFS := fstest.MapFS{
		"file":                     {},
		"prefix_" + today + ".mp3": {},
	}

	testCases := []struct {
		name string

		outputDirectory string
		outputFS        fs.FS
		prefix          string
		dateFormat      string
		extension       string

		expected string
	}{
		{
			name: "empty filename",

			outputDirectory: "",
			outputFS:        outputFS,
			prefix:          "",
			dateFormat:      "",
			extension:       "",

			expected: today,
		},
		{
			name: "second file",

			outputDirectory: "",
			outputFS:        outputFS,
			prefix:          "",
			dateFormat:      "file",
			extension:       "",

			expected: "file.1",
		},
		{
			name: "in directory",

			outputFS:        outputFS,
			outputDirectory: "output",
			prefix:          "",
			dateFormat:      layout,
			extension:       "",

			expected: "output/" + today,
		},
		{
			name: "directory name with slash",

			outputFS:        outputFS,
			outputDirectory: "output/",
			prefix:          "",
			dateFormat:      layout,
			extension:       "",

			expected: "output/" + today,
		},
		{
			name: "with prefix",

			outputFS:        outputFS,
			outputDirectory: "",
			prefix:          "prefix_",
			dateFormat:      layout,
			extension:       "",

			expected: "prefix_" + today,
		},
		{
			name: "with extension",

			outputFS:        outputFS,
			outputDirectory: "",
			prefix:          "",
			dateFormat:      layout,
			extension:       "mp3",

			expected: today + ".mp3",
		},
		{
			name: "extension with dot",

			outputFS:        outputFS,
			outputDirectory: "",
			prefix:          "",
			dateFormat:      layout,
			extension:       ".mp3",

			expected: today + ".mp3",
		},
		{
			name: "full data",

			outputFS:        outputFS,
			outputDirectory: "output",
			prefix:          "prefix_",
			dateFormat:      layout,
			extension:       ".mp3",

			expected: "output/prefix_" + today + ".1.mp3",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			builder := NewDatedFileBuilder(
				tc.outputDirectory,
				tc.outputFS,
				tc.prefix,
				tc.dateFormat,
			)
			actual := builder.GetFileName(tc.extension)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMimeTypeDetection(t *testing.T) {
	t.Parallel()

	const (
		contentType     = "audio/mpeg"
		expectExtension = ".mp3"
	)

	extension, err := mime.ExtensionsByType(contentType)

	assert.NoError(t, err)
	assert.Contains(t, extension, expectExtension)
}
