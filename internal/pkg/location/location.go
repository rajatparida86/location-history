package location

import (
	"fmt"
	"github.com/rajatparida86/location-history/internal/pkg/config"
	log "github.com/sirupsen/logrus"
	"time"
)

type History struct {
	store map[string][]*Location
}

func SetUpHistoryStore(conf *config.Configuration) *History {
	h := &History{store: make(map[string][]*Location)}

	// Expire location entries as per TTL
	go func() {
		for current := range time.Tick(time.Second) {
			for orderId, history := range h.store {
				for i := len(history) - 1; i >= 0; i-- {
					if current.Unix()-history[i].createdAt > int64(conf.StoreTtl) {
						h.store[orderId] = history[i+1:]
						break
					}
				}
			}
		}
	}()
	return h
}

type Location struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
	createdAt int64
}

func (h *History) UpdateLocation(orderId string, location *Location) error {
	locationHistory, ok := h.store[orderId]
	if !ok {
		locationHistory = make([]*Location, 0)
	}
	location.createdAt = time.Now().Unix()
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
