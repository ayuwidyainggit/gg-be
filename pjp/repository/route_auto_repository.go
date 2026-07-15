package repository

import (
	"context"
	"scyllax-pjp/model"
)

type RouteAutoRepository interface {
	Insert(ctx context.Context, route []model.Route)
}
