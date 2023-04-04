package compression

import "context"

type CompressionProvider struct {
}

func InitialiseCompressionProvider(ctx context.Context) (CompressionProvider, error) {

	return CompressionProvider{}, nil
}
