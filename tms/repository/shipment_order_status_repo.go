package repository

import (
	"context"
	"scyllax-tms/model"
)

type ShipmentOrderStatusRepo interface {
	CreateOrUpdate(ctx context.Context, data model.ShipmentOrderStatus) error
	//Update(ctx context.Context, data model.ShipmentOrderStatus) error
}
