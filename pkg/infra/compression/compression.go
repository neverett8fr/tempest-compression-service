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

func (cp *CompressionProvider) compressLZW(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}

	// Initialize the dictionary with all possible bytes
	dictionary := make(map[string]int)
	for i := 0; i < 256; i++ {
		dictionary[string([]byte{byte(i)})] = i
	}

	var prefix string
	result := make([]int, 0, len(data)/2)

	for _, c := range data {
		s := prefix + string(c)
		if _, ok := dictionary[s]; ok {
			prefix = s
		} else {
			result = append(result, dictionary[prefix])
			dictionary[s] = len(dictionary)
			prefix = string(c)

			// Reset dictionary if it becomes too large
			if len(dictionary) == 4096 {
				dictionary = make(map[string]int)
				for i := 0; i < 256; i++ {
					dictionary[string([]byte{byte(i)})] = i
				}
			}
		}
	}

	if len(prefix) > 0 {
		result = append(result, dictionary[prefix])
	}

	// Convert the result to bytes
	resultBytes := make([]byte, len(result)*2+4)
	binary.LittleEndian.PutUint32(resultBytes, 0x4C5A5700)
	for i, v := range result {
		resultBytes[i*2+4] = byte(v)
		resultBytes[i*2+5] = byte(v >> 8)
	}

	return resultBytes, nil
}

func (cp *CompressionProvider) decompressLZW(fileIn []byte) ([]byte, error) {
	if len(fileIn) == 0 {
		return nil, errors.New("input data is empty")
	}

	// Initialize the dictionary with the 256 ASCII characters.
	dictSize := int64(256)
	dict := make(map[int64][]byte)
	for i := 0; int64(i) < dictSize; i++ {
		dict[int64(i)] = []byte{byte(i)}
	}

	// Initialize the current code and previous code to invalid values.
	currentCode := int64(-1)
	previousCode := int64(-1)

	// Initialize the output buffer and output buffer index.
	output := make([]byte, 0, len(fileIn))
	outputIndex := 0

	// Read the compressed input.
	for len(fileIn) > 0 {
		var code int64
		err := binary.Read(bytes.NewReader(fileIn[:2]), binary.LittleEndian, &code)
		if err != nil {
			return nil, err
		}
		fileIn = fileIn[2:]
		currentCode = code // add this line to update currentCode

		entry, ok := dict[currentCode]
		if !ok {
			if currentCode == dictSize {
				entry = append(dict[previousCode], dict[previousCode][0])
			} else {
				return nil, errors.New("invalid compressed data")
			}
		}

		output = append(output, entry...)
		outputIndex += len(entry)

		if previousCode != -1 {
			dict[dictSize] = append(dict[previousCode], entry[0])
			dictSize++
		}

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
