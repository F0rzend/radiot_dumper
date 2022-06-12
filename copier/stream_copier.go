package copier

import (
	"bufio"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

const (
	fileHeaderSize = 262 // Maximal size of a file header. It's enough for detecting the mime type.
)

type StreamCopier struct {
	client *http.Client
	logger zerolog.Logger
}

type GetOutputFunc func(ext string) (io.WriteCloser, error)

func NewStreamCopier(client *http.Client, logger zerolog.Logger) *StreamCopier {
	if client == nil {
		client = http.DefaultClient
	}
	return &StreamCopier{
		client: client,
		logger: logger,
	}
}

var (
	ErrStreamClosed = errors.New("stream closed")
)

func (d *StreamCopier) CopyStream(url string, getOutput GetOutputFunc) error {
	log := d.logger.With().Str("request_id", uuid.New().String()).Logger()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNotFound {
		d.logger.Debug().Int("status_code", resp.StatusCode).Msg("got response")
	}

	if resp.StatusCode != http.StatusOK {
		return ErrStreamClosed
	}

	buf := bufio.NewReader(resp.Body)
	fileHeader, err := buf.Peek(fileHeaderSize)
	if err != nil && err != io.EOF {
		return err
	}
	mime := mimetype.Detect(fileHeader)
	fileExtension := mime.Extension()
	log.Debug().Str("extension", fileExtension).Msg("detected extension")

	output, err := getOutput(fileExtension)
	if err != nil {
		return err
	}
	defer func() {
		if err := output.Close(); err != nil {
			log.Error().Err(err).Msg("error closing output")
		}
	}()
	if file, ok := output.(*os.File); ok {
		log.Debug().Str("filename", file.Name()).Msg("output in file")
	}

	bytesCopied, err := io.Copy(output, buf)
	log.Debug().Int64("bytes_copied", bytesCopied).Msg("copied bytes")

	return err
}
