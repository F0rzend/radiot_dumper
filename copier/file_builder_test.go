package copier

import (
	"github.com/gabriel-vasile/mimetype"
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

func TestFileDetection(t *testing.T) {
	t.Skip("Library doesn't support this")

	t.Parallel()

	// Radio-T stream header
	header := []byte{
		112, 44, 209, 245, 107, 234, 157, 24, 68, 65, 203, 234, 91, 217, 92, 73, 96, 16, 165, 1, 242, 4, 238, 116, 56,
		128, 52, 53, 117, 94, 102, 36, 253, 83, 85, 108, 139, 149, 47, 151, 242, 21, 53, 27, 182, 95, 85, 145, 197, 130,
		94, 11, 37, 107, 43, 248, 53, 209, 117, 91, 174, 48, 44, 49, 35, 126, 230, 33, 171, 150, 173, 81, 214, 149, 99,
		220, 174, 89, 239, 211, 127, 243, 91, 252, 165, 55, 41, 36, 88, 82, 99, 123, 41, 172, 169, 106, 210, 218, 167,
		238, 234, 95, 175, 91, 120, 229, 188, 113, 199, 120, 210, 217, 164, 169, 115, 243, 207, 154, 255, 254, 111, 255,
		247, 87, 255, 255, 255, 227, 7, 0, 24, 118, 82, 29, 123, 108, 221, 142, 183, 44, 159, 75, 88, 210, 185, 137, 3,
		24, 180, 233, 151, 146, 152, 130, 193, 143, 132, 32, 4, 160, 228, 48, 236, 120, 48, 207, 65, 21, 145, 25, 144,
		252, 2, 22, 168, 204, 4, 1, 166, 182, 54, 148, 185, 223, 116, 209, 12, 235, 119, 95, 140, 61, 92, 42, 55, 156,
		64, 183, 165, 56, 159, 73, 167, 190, 14, 95, 206, 66, 211, 90, 237, 134, 71, 46, 92, 174, 10, 87, 202, 53, 105,
		74, 212, 13, 131, 206, 207, 38, 139, 42, 198, 146, 142, 159, 249, 255, 251, 144, 196, 205, 128, 38, 37, 169, 73,
		249, 205, 0, 67, 51, 44, 234, 119, 55, 128, 0, 190, 78, 119, 255, 185, 198,
	}

	headerMime := mimetype.Detect(header)
	fileExtension := headerMime.Extension()

	if fileExtension == "" {
		t.Errorf("File extension not detected")
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
