package compression

import (
	"io/ioutil"
	"net/http"
	"os"
)

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
