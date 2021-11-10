package main

import (
	"context"
	"github.com/rajatparida86/location-history/internal/api"
	"github.com/rajatparida86/location-history/internal/pkg/config"
	"github.com/rajatparida86/location-history/internal/pkg/location"
	"github.com/rajatparida86/location-history/internal/pkg/observabilitySDK"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup opentelemetry
	err := observabilitySDK.InitOtel(context.Background())
	if err != nil {
		log.Infof("failed to init observability SDK. A call to tracer will result in NOOP tracer: %s", err)
	}
	defer observabilitySDK.Shutdown()

	/*otelLauncher := launcher.ConfigureOpentelemetry()
	defer otelLauncher.Shutdown()*/

	conf := config.SetUpConfiguration()
	store := location.NewInMemoryStore(conf)
	a := api.New(store, conf)
	a.Run()
}
