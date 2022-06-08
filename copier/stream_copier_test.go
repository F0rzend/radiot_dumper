package copier

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testLogger = zerolog.New(nil)

type closableBuffer struct {
	bytes.Buffer
}

func (c *closableBuffer) Close() error {
	return nil
}

func TestStreamCopier_CopyStream_Success(t *testing.T) {
	t.Parallel()

	serverOutput := []byte("Hello World!")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(serverOutput)
		assert.NoError(t, err)
	}))
	defer server.Close()

	copier := NewStreamCopier(http.DefaultClient, testLogger)

	output := new(closableBuffer)

	err := copier.CopyStream(server.URL, func(_ string) (io.WriteCloser, error) {
		return output, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, serverOutput, output.Bytes())
}

func TestStreamCopier_CopyStream_WithStreamClosed(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	copier := NewStreamCopier(http.DefaultClient, testLogger)

	err := copier.CopyStream(server.URL, func(_ string) (io.WriteCloser, error) {
		return new(closableBuffer), nil
	})

	assert.ErrorIs(t, err, ErrStreamClosed)
}
