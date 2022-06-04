package internal

import (
	"errors"
	"io"
	"net/http"
	"time"
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

func (d *StreamCopier) CopyStream(url string, output io.Writer) error {
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

	if _, err := io.Copy(output, resp.Body); err != nil {
		return err
	}

	return nil
}

func (d *StreamCopier) ListenAndCopy(
	url string,
	output io.Writer,
	timeout time.Duration,
) error {
	for {
		if err := d.CopyStream(url, output); err != nil && err != ErrStreamClosed {
			return err
		}

		time.Sleep(timeout)
	}
}
