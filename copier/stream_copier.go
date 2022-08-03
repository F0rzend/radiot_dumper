package copier

import (
	"bytes"
	"context"
	"errors"
	"github.com/rs/zerolog"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

type StreamCopier struct {
	client *http.Client
}

type GetOutputFunc func(ext string) (io.WriteCloser, error)

func NewStreamCopier(client *http.Client) *StreamCopier {
	if client == nil {
		client = http.DefaultClient
	}
	return &StreamCopier{
		client: client,
	}
}

var (
	ErrStreamClosed = errors.New("stream closed")
)

func (d *StreamCopier) CopyStream(ctx context.Context, url string, getOutput GetOutputFunc) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("error closing response body")
		}
	}()
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNotFound {
		zerolog.Ctx(ctx).Debug().Int("status_code", resp.StatusCode).Msg("got response")
	}

	if resp.StatusCode != http.StatusOK {
		return ErrStreamClosed
	}

	zerolog.Ctx(ctx).Info().Msg("recording started")

	fileExtension, body, err := DetectExtension(ctx, resp)
	if err != nil {
		return err
	}
	zerolog.Ctx(ctx).Debug().Str("extension", fileExtension).Msg("detected extension")

	output, err := getOutput(fileExtension)
	if err != nil {
		return err
	}
	defer func() {
		if err := output.Close(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("error closing output")
		}
	}()
	if file, ok := output.(*os.File); ok {
		zerolog.Ctx(ctx).Debug().Str("filename", file.Name()).Msg("output in file")
	}

	bytesCopied, err := io.Copy(output, body)
	zerolog.Ctx(ctx).Debug().Int64("bytes_copied", bytesCopied).Msg("copied bytes")

	zerolog.Ctx(ctx).Info().Msg("recording finished")
	return err
}

// DetectExtension returns response extension by first looking into response headers.
// As a fallback, it looks into response body and returns the extension and a new
// body containing the original content.
func DetectExtension(ctx context.Context, r *http.Response) (string, io.Reader, error) {
	contentType := r.Header.Get("Content-Type")
	ext, err := extensionFromContentType(ctx, contentType)
	if err == nil {
		return ext, r.Body, err
	}
	zerolog.Ctx(ctx).Warn().Err(err).Msg("extension detection from bytes")
	return extensionFromBody(r.Body)
}

func extensionFromContentType(ctx context.Context, contentType string) (string, error) {
	weightedExtensions := []string{".mp3"}

	acceptableExtensions, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return "", err
	}
	if len(acceptableExtensions) == 0 {
		return ".unknown", nil
	}
	zerolog.Ctx(ctx).
		Debug().
		Interface("extensions", acceptableExtensions).
		Msg("extensions from content type")

	for _, preferred := range weightedExtensions {
		for _, option := range acceptableExtensions {
			if preferred == option {
				return option, nil
			}
		}

	}

	return acceptableExtensions[0], nil
}

// extensionFromBody returns the extension of the file contained by body and a
// new body containing the original input file.
func extensionFromBody(body io.Reader) (ext string, newBody io.Reader, err error) {
	// header will store the bytes mimetype uses for detection.
	header := bytes.NewBuffer(nil)

	// After DetectReader, the data read from input is copied into header.
	mtype, err := mimetype.DetectReader(io.TeeReader(body, header))

	// Concatenate back the header to the rest of the file.
	// newBody now contains the complete, original data.
	newBody = io.MultiReader(header, body)

	return mtype.Extension(), newBody, err
}
