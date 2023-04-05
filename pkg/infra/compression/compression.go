package compression

import (
	"bytes"
	"compress/gzip"
	"io"
)

func (cp *CompressionProvider) Compress(fileIn []byte) ([]byte, error) {

	file, err := bytesToFile(fileIn)
	if err != nil {
		return nil, err
	}

	// Create a buffer to hold the compressed data
	var buf bytes.Buffer

	// Create a gzip writer with the buffer as its destination
	gzipWriter := gzip.NewWriter(&buf)

	// Copy the input file data to the gzip writer, which compresses it and writes to the buffer
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		return nil, err
	}

	// Close the gzip writer to flush any remaining data
	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (cp *CompressionProvider) Decompress(fileIn []byte) ([]byte, error) {
	// Create a bytes buffer with the input file data
	buf := bytes.NewBuffer(fileIn)

	// Create a gzip reader with the buffer as its source
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	// Create a buffer to hold the decompressed data
	var result bytes.Buffer

	// Copy the gzip reader data to the result buffer, which decompresses it
	_, err = io.Copy(&result, gzipReader)
	if err != nil {
		return nil, err
	}

	// Close the gzip reader
	err = gzipReader.Close()
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
