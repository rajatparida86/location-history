package location

import (
	"fmt"
	"github.com/rajatparida86/location-history/internal/pkg/config"
	"sync"
	"time"
)

type InMemoryStore struct {
	history map[string][]*Location
	lock    sync.Mutex
}

func NewInMemoryStore(conf *config.Configuration) *InMemoryStore {
	store := &InMemoryStore{history: make(map[string][]*Location)}

	// Expire location entries as per TTL
	go func() {
		for current := range time.Tick(1 * time.Minute) {
			store.lock.Lock()
			for orderId, history := range store.history {
				for i := len(history) - 1; i >= 0; i-- {
					if current.Unix()-history[i].createdAt > int64(conf.StoreTtl) {
						store.history[orderId] = history[i+1:]
						break
					}
				}
			}
			store.lock.Unlock()
		}
	}()
	return store
}

func (i *InMemoryStore) UpdateLocation(orderId string, location *Location) error {
	i.lock.Lock()
	locationHistory, ok := i.history[orderId]
	if !ok {
		locationHistory = make([]*Location, 0)
	}
	location.createdAt = time.Now().Unix()
	locationHistory = append(locationHistory, location)
	i.history[orderId] = locationHistory
	i.lock.Unlock()
	return nil
}

func (i *InMemoryStore) GetLocation(orderId string, depth *int) ([]*Location, error) {
	locationHistory, ok := i.history[orderId]
	if !ok {
		return nil, fmt.Errorf("order with id - %v not found", orderId)
	}
	if depth == nil {
		return locationHistory, nil
	}
	return locationHistory[len(locationHistory)-*depth:], nil
}

func (i *InMemoryStore) DeleteLocation(orderId string) error {
	i.lock.Lock()
	_, ok := i.history[orderId]
	if !ok {
		return fmt.Errorf("order not found")
	}
	delete(i.history, orderId)
	i.lock.Unlock()
	return nil
}
