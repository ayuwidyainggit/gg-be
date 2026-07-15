package utils

import (
	"context"
	"github.com/google/uuid"
	"scyllax-pjp/data/response"
)

func ResponseInterceptor(ctx context.Context, resp *response.Response) {
	traceId := uuid.Must(uuid.NewRandom())
	resp.TraceID = traceId.String()
}
