package interceptors

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func MetadataRequestPath() func(ctx context.Context, request *http.Request) metadata.MD {
	return func(ctx context.Context, request *http.Request) metadata.MD {
		return metadata.New(map[string]string{
			"path": request.URL.Path,
		})
	}
}
