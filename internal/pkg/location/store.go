package location

import "context"

type Store interface {
	UpdateLocation(orderId string, location *Location) error
	GetLocation(ctx context.Context, orderId string, depth *int) ([]*Location, error)
	DeleteLocation(orderId string) error
}
