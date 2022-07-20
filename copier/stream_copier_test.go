package copier

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
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

func TestStreamCopier(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		handler http.HandlerFunc
		err     error
		output  []byte
	}{
		{
			name:    "success",
			handler: handlerSuccess,
			output:  []byte("Hello World!"),
			err:     nil,
		},
		{
			name:    "not found",
			handler: handlerNotFound,
			output:  []byte{},
			err:     ErrStreamClosed,
		},
		{
			name:    "with interrupt",
			handler: getHandlerWithInterrupt(),
			output:  []byte("12"),
			err:     nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tc.handler)
			defer server.Close()

			copier := NewStreamCopier(http.DefaultClient, testLogger)
			output := new(closableBuffer)

			err := copier.CopyStream(server.URL, func(_ string) (io.WriteCloser, error) {
				return output, nil
			})
			assert.Equal(t, tc.err, err)
		})
	}
}

func handlerSuccess(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello World!"))
	if err != nil {
		log.Println(err)
	}
}

func handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func getHandlerWithInterrupt() http.HandlerFunc {
	responses := []struct {
		status int
		body   []byte
	}{
		{
			status: http.StatusOK,
			body:   []byte("1"),
		},
		{
			status: http.StatusNotFound,
		},
		{
			status: http.StatusOK,
			body:   []byte("2"),
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.NotFound(w, r)
			return
		}

		for _, response := range responses {
			w.WriteHeader(response.status)
			_, err := w.Write(response.body)
			log.Println(err)
			flusher.Flush()
		}
	}
}

func TestDetectExtension(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		header   string
		body     string
		expected string
	}{
		{
			name:     "by header",
			header:   "audio/mpeg",
			body:     "",
			expected: ".mp3",
		},
		{
			name:     "by body",
			header:   "",
			body:     "Hello, World!",
			expected: ".txt",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r := &http.Response{
				Header: http.Header{
					"Content-Type": []string{tc.header},
				},
				Body: &closableBuffer{
					Buffer: *bytes.NewBuffer(
						[]byte(tc.body),
					),
				},
			}

			actual, err := DetectExtension(r)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
