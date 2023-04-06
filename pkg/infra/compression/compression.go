package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func (cp *CompressionProvider) compressRLE(fileIn []byte) ([]byte, error) {
	if len(fileIn) == 0 {
		return nil, errors.New("input data is empty")
	}

	count := 1
	lastByte := fileIn[0]
	compressedData := make([]byte, 0)

	for _, b := range fileIn[1:] {
		if b == lastByte && count < 255 {
			count++
		} else {
			compressedData = append(compressedData, byte(count), lastByte)
			count = 1
			lastByte = b
		}
	}

	compressedData = append(compressedData, byte(count), lastByte)

	return compressedData, nil
}

func (cp *CompressionProvider) decompressRLE(fileIn []byte) ([]byte, error) {
	if len(fileIn) == 0 {
		return nil, errors.New("input data is empty")
	}

	decompressedData := make([]byte, 0)

	for i := 0; i < len(fileIn); i += 2 {
		count := int(fileIn[i])
		byteVal := fileIn[i+1]

		for j := 0; j < count; j++ {
			decompressedData = append(decompressedData, byteVal)
		}
	}

	return decompressedData, nil
}

func (cp *CompressionProvider) compressLZW(fileIn []byte) ([]byte, error) {
	if len(fileIn) == 0 {
		return nil, errors.New("input data is empty")
	}

	// Initialize the dictionary with all possible bytes
	dictionary := make(map[string]int)
	for i := 0; i < 256; i++ {
		dictionary[string([]byte{byte(i)})] = i
	}

	result := make([]int, 0)
	prefix := ""

	for _, c := range fileIn {
		s := prefix + string(c)
		if _, ok := dictionary[s]; ok {
			prefix = s
		} else {
			result = append(result, dictionary[prefix])
			dictionary[s] = len(dictionary)
			prefix = string(c)
		}
	}

	if len(prefix) > 0 {
		result = append(result, dictionary[prefix])
	}

	// Convert the result to bytes
	output := make([]byte, len(result)*2)
	for i, v := range result {
		binary.BigEndian.PutUint16(output[i*2:], uint16(v))
	}

	return output, nil
}

func (cp *CompressionProvider) decompressLZW(fileIn []byte) ([]byte, error) {
	if len(fileIn) == 0 {
		return nil, errors.New("input data is empty")
	}

	// Initialize the dictionary with the 256 ASCII characters.
	dictSize := 256
	dict := make(map[int][]byte)
	for i := 0; i < dictSize; i++ {
		dict[i] = []byte{byte(i)}
	}

	// Initialize the current code and previous code to invalid values.
	currentCode := -1
	previousCode := -1

	// Initialize the output buffer and output buffer index.
	output := make([]byte, 0, len(fileIn))
	outputIndex := 0

	// Iterate over the compressed input.
	for i := 0; i < len(fileIn); i += 2 {
		// Read the next code.
		currentCode = int(binary.BigEndian.Uint16(fileIn[i : i+2]))

		// Handle the first code.
		if previousCode == -1 {
			output = append(output, dict[currentCode]...)
			outputIndex += len(dict[currentCode])
			previousCode = currentCode
			continue
		}

		// Handle the rest of the codes.
		entry, ok := dict[currentCode]
		if ok {
			output = append(output, entry...)
			outputIndex += len(entry)

			// Add the new entry to the dictionary.
			dict[dictSize] = append(dict[previousCode], entry[0])
			dictSize++
		} else {
			// Handle the case where the code is not in the dictionary.
			entry = append(dict[previousCode], dict[previousCode][0])
			output = append(output, entry...)
			outputIndex += len(entry)

			// Add the new entry to the dictionary.
			dict[dictSize] = entry
			dictSize++
		}

		// Update the previous code.
		previousCode = currentCode
	}

	return output, nil
}

func (cp *CompressionProvider) compressGZip(fileIn []byte) ([]byte, error) {
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

func (cp *CompressionProvider) decompressGZip(fileIn []byte) ([]byte, error) {
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

func (cp *CompressionProvider) Compress(fileIn []byte) ([]byte, error) {

	// call service
	if cp.UseML {
		resp, err := cp.CallML(fileIn)
		if err != nil {
			return nil, fmt.Errorf("error calling ML service, err %v", err)
		}

		return cp.CompressAuto(fileIn, resp.Data.MethodExt)
	}

	return cp.CompressAuto(fileIn, ".gzip")
}

func (cp *CompressionProvider) Decompress(fileIn []byte) ([]byte, error) {
	return cp.DecompressAuto(fileIn)
}
