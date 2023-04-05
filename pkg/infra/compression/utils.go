package compression

import "net/http"

func GetFileType(file []byte) string {
	fileType := http.DetectContentType(file)

	return fileType
}
