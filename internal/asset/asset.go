package asset

import (
	"io"
	"net/http"
	"strconv"
)

type (
	Reader     = Asset[io.Reader]
	ReadCloser = Asset[io.ReadCloser]
)

// Asset is an open reader that can read possibly compressed data.
type Asset[Reader io.Reader] struct {
	ContentType   string
	ContentLength int64
	r             Reader
}

// NewAssetReader creates a new compressed asset reader.
func NewAssetReader[Reader io.Reader](r Reader, contentType string, contentLength int64) Asset[Reader] {
	return Asset[Reader]{
		ContentType:   contentType,
		ContentLength: contentLength,
		r:             r,
	}
}

// Close closes the underlying reader.
// If Reader is not derived from io.Closer, then this method does nothing.
func (c Asset[Reader]) Close() error {
	if closer, ok := any(c.r).(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Reader returns the underlying reader.
func (c Asset[Reader]) Reader() io.Reader {
	return c.r
}

// WriteToResponse writes the compressed data to an HTTP response.
// This sets the appropriate headers and writes the compressed data.
func (c Asset[Reader]) WriteToResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", c.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(c.ContentLength, 10))
	_, err := io.Copy(w, c.r)
	return err
}
