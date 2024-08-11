package asset

import (
	"io"
	"net/http"
)

// Compression is the type of compression used for an asset.
type Compression string

const (
	NotCompressed     Compression = ""
	GzipCompression   Compression = "gzip"
	ZstdCompression   Compression = "zstd"
	BrotliCompression Compression = "brotli"
)

// CompressedAsset is an open reader that can read possibly compressed data.
type CompressedAsset[Reader io.Reader] struct {
	MIMEType    string
	Compression Compression

	r Reader
}

// NewCompressedAssetReader creates a new compressed asset reader.
func NewCompressedAssetReader[Reader io.Reader](r Reader, compression Compression, mimeType string) CompressedAsset[Reader] {
	return CompressedAsset[Reader]{
		Compression: compression,
		MIMEType:    mimeType,
		r:           r,
	}
}

// Close closes the underlying reader.
// If Reader is not derived from io.Closer, then this method does nothing.
func (c CompressedAsset[Reader]) Close() error {
	if closer, ok := any(c.r).(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// CompressedReader returns the underlying compressed reader.
// This reader reads possibly compressed data.
func (c CompressedAsset[Reader]) CompressedReader() io.Reader {
	return c.r
}

// WriteToResponse writes the compressed data to an HTTP response.
// This sets the appropriate headers and writes the compressed data.
func (c CompressedAsset[Reader]) WriteToResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Encoding", string(c.Compression))
	w.Header().Set("Content-Type", c.MIMEType)
	_, err := io.Copy(w, c.r)
	return err
}
