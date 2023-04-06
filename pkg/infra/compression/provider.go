package compression

import (
	"context"
	"fmt"
	"tempest-compression-service/pkg/config"
)

type CompressionProvider struct {
	MLPath string
	UseML  bool
}

func InitialiseCompressionProvider(ctx context.Context, conf config.ML) (CompressionProvider, error) {
	pathToService := fmt.Sprintf("%s:%v/%s", conf.Host, conf.Port, "compression/decide")
	if conf.Port == 0 {
		pathToService = fmt.Sprintf("%s/%s", conf.Host, "compression/decide")
	}

	return CompressionProvider{
		MLPath: pathToService,
		UseML:  conf.UseML,
	}, nil
}
