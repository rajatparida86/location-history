package location

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type History struct {
	store map[string][]*Location
}

func SetUpHistoryStore() *History {
	return &History{store: make(map[string][]*Location)}
}

type Location struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

func (h *History) UpdateLocation(orderId string, location *Location) error {
	locationHistory, ok := h.store[orderId]
	if !ok {
		locationHistory = make([]*Location, 0)
	}
	locationHistory = append(locationHistory, location)
	h.store[orderId] = locationHistory
	return nil
}

func (h *History) GetLocation(orderId string, depth *int) ([]*Location, error) {
	locationHistory, ok := h.store[orderId]
	if !ok {
		for k, v := range h.store {
			log.Info("key %s", k)
			log.Info("val %s", v)
		}
		return nil, fmt.Errorf("order with id - %v not found", orderId)

	}
	if depth == nil {
		return locationHistory, nil
	}
	return locationHistory[len(locationHistory)-*depth:], nil
}

func (h *History) DeleteLocation(orderId string) error {
	_, ok := h.store[orderId]
	if !ok {
		return fmt.Errorf("order not found")
	}
	delete(h.store, orderId)
	return nil
}
