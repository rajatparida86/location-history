package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rajatparida86/location-history/internal/pkg/config"
	"github.com/rajatparida86/location-history/internal/pkg/location"
	log "github.com/sirupsen/logrus"

	middleware "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"net/http"
)

type Api struct {
	store location.Store
	conf  *config.Configuration
}

func New(store location.Store, conf *config.Configuration) *Api {
	return &Api{
		store,
		conf,
	}
}

func (a *Api) Run() {
	r := mux.NewRouter()
	// Auto instrumentation for Gorilla MUX
	r.Use(middleware.Middleware("location-history"))
	r.HandleFunc("/health", a.healthHandler).Methods("GET")
	r.HandleFunc("/location/{orderId}/now", a.addLocation).Methods("POST")
	r.HandleFunc("/location/{orderId}", a.getLocation).Methods("GET")
	r.HandleFunc("/location/{orderId}", a.getLocation).Methods("GET").Queries("max", "{max:[0-9]+}")
	r.HandleFunc("/location/{orderId}", a.deleteLocationHistory).Methods("DELETE")

	serverAddress := fmt.Sprintf(":%s", a.conf.Port)
	log.Infof("server started on address %s", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}

func (a *Api) writeResponse(w http.ResponseWriter, status int, data interface{}) {
	resp, err := json.Marshal(data)
	if err != nil {
		a.writeSimpleResponse(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)
}

func (a *Api) writeSimpleResponse(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	w.Write([]byte(http.StatusText(code)))
}

type failedResponse struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func newFailedResponse(status int, detail string) *failedResponse {
	return &failedResponse{
		Status: status,
		Detail: detail,
	}
}

func (a *Api) writeFailedResponse(w http.ResponseWriter, code int, err error) {
	fr := newFailedResponse(code, err.Error())
	a.writeResponse(w, code, fr)
}
