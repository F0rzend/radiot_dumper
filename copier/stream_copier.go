package copier

import (
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

type StreamCopier struct {
	client *http.Client
}

func NewStreamCopierService(client *http.Client) *StreamCopier {
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

func (d *StreamCopier) CopyStream(url string, fileBuilder FileBuilder) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer logClosing(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return ErrStreamClosed
	}

	mime, err := mimetype.DetectReader(resp.Body)
	if err != nil {
		return err
	}
	fileExtension := mime.Extension()

	output, err := fileBuilder.CreateFile(fileExtension)
	if err != nil {
		return err
	}
	defer logClosing(output)

	if _, err := io.Copy(output, resp.Body); err != nil {
		return err
	}
	log.Println("copied", url)

	return nil
}

func (d *StreamCopier) ListenAndCopy(
	url string,
	fileBuilder FileBuilder,
	timeout time.Duration,
) error {
	for {
		if err := d.CopyStream(url, fileBuilder); err != nil && err != ErrStreamClosed {
			return err
		}
		time.Sleep(timeout)
	}
}
