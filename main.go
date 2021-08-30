package main

import (
	"github.com/rajatparida86/location-history/internal/api"
	"github.com/rajatparida86/location-history/internal/pkg/config"
	"github.com/rajatparida86/location-history/internal/pkg/location"
)

func main() {
	conf := config.SetUpConfiguration()
	store := location.NewInMemoryStore(conf)
	a := api.New(store, conf)
	a.Run()
}
