package compression

import (
	"fmt"
	"net/http"
)

type MLResponse struct {
	Errors []error            `json:"errors"`
	Data   expectedMLResponse `json:"data"`
}

type expectedMLResponse struct {
	MethodExt string `json:"method_ext"`
	MethodNum int    `json:"method_num"`
}

func fileToMLNum(in string) int {
	// https://github.com/neverett8fr/tempest-compression-decision-tree-proof-of-concept/blob/main/files_to_compress/file_type_mapping_num
	mp := map[string]int{
		"unknown":    0,
		"text/plain": 1,
		"image/png":  2,
		"image/jpeg": 3,
	}

	return mp[in]
}

// "convert" - convert either to compressed, or uncompressed
func (cp *CompressionProvider) CallML(data []byte) (*MLResponse, error) {

	fileLength := len(data)
	fileExt := fileToMLNum(GetFileType(data))

	// compression not yet implemented
	if fileExt == 0 {
		return &MLResponse{
			Data: expectedMLResponse{
				MethodExt: "none",
				MethodNum: 0,
			},
		}, nil
	}

	route := fmt.Sprintf("%s/%v/%v", cp.MLPath, fileExt, fileLength)
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request, err %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error calling ML service, err %v", err)
	}

	return readBody(*resp)
}
