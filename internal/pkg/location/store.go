package location

type Store interface {
	UpdateLocation(orderId string, location *Location) error
	GetLocation(orderId string, depth *int) ([]*Location, error)
	DeleteLocation(orderId string) error
}
