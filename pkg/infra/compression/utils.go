package compression

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func readBody(resp http.Response) (*MLResponse, error) {

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body, err %v", err)
	}

	applicationResponse := MLResponse{}
	err = json.Unmarshal(body, &applicationResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling body, err %v", err)
	}

	return &applicationResponse, nil
}

type mapFuncHelps struct {
	handyFunc func([]byte) ([]byte, error)
}

func emptyHelperFunc(in []byte) ([]byte, error) {
	return in, nil
}

func (cp *CompressionProvider) DecompressAuto(data []byte) ([]byte, error) {
	mp := map[string]mapFuncHelps{
		"none": {
			handyFunc: emptyHelperFunc,
		},
		".rle": {
			handyFunc: cp.decompressRLE,
		},
		".gzip": {
			handyFunc: cp.decompressGZip,
		},
		".lzw": {
			handyFunc: cp.decompressLZW,
		},
	}

	comp, err := cp.DetectCompressionType(data)
	if err != nil {
		return nil, fmt.Errorf("error detecting compression, err %v", err)
	}
	if comp == "" {
		return data, nil
	}

	// lzw doesnt currently work! so this is a workaround
	if comp == ".lzw" {
		return mp[".gzip"].handyFunc(data)
	}

	return mp[comp].handyFunc(data)
}

func (cp *CompressionProvider) CompressAuto(data []byte, method string) ([]byte, error) {
	mp := map[string]mapFuncHelps{
		"none": {
			handyFunc: emptyHelperFunc,
		},
		".rle": {
			handyFunc: cp.compressRLE,
		},
		".gzip": {
			handyFunc: cp.compressGZip,
		},
		".lzw": {
			handyFunc: cp.compressLZW,
		},
	}

	if method == "" {
		return data, nil
	}

	// lzw doesnt currently work! so this is a workaround
	if method == ".lzw" {
		return mp[".gzip"].handyFunc(data)
	}

	return mp[method].handyFunc(data)
}

func (cp *CompressionProvider) DetectCompressionType(data []byte) (string, error) {
	if len(data) < 3 {
		return "", fmt.Errorf("input data is too short")
	}

	// Check for RLE compression
	if data[0] == 0x01 && data[1] == 0x00 {
		return ".rle", nil
	}

	// Check for gzip compression
	if data[0] == 0x1f && data[1] == 0x8b {
		return ".gzip", nil
	}

	// Check for LZW compression
	if len(data) >= 4 && binary.LittleEndian.Uint32(data[0:4]) == 0x4C5A5700 {
		return ".lzw", nil
	}

	return "", nil
}

func GetFileType(file []byte) string {
	fileType := http.DetectContentType(file)

	return fileType
}

func bytesToFile(data []byte) (*os.File, error) {
	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		return nil, err
	}
	defer tmpfile.Close()

	// Write the []byte data to the temporary file
	_, err = tmpfile.Write(data)
	if err != nil {
		return nil, err
	}

	// Sync and close the temporary file to ensure the data is flushed to disk
	err = tmpfile.Sync()
	if err != nil {
		return nil, err
	}
	err = tmpfile.Close()
	if err != nil {
		return nil, err
	}

	// Open the temporary file and obtain an *os.File handle
	file, err := os.Open(tmpfile.Name())
	if err != nil {
		return nil, err
	}

	return file, nil
}
