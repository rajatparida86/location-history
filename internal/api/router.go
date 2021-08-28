package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rajatparida86/location-history/internal/pkg/location"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (a *Api) healthHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("health check")
	a.writeSimpleResponse(w, http.StatusOK)
}

func (a *Api) addLocation(w http.ResponseWriter, r *http.Request) {
	var loc *location.Location
	id := mux.Vars(r)["orderId"]

	err := json.NewDecoder(r.Body).Decode(&loc)
	if err != nil {
		a.writeFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	err = a.store.UpdateLocation(id, loc)
	if err != nil {
		a.writeFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	a.writeSimpleResponse(w, http.StatusOK)
}

func (a *Api) getLocation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["orderId"]
	var d *int
	if r.FormValue("max") != "" {
		depth, _ := strconv.Atoi(r.FormValue("max"))
		d = new(int)
		*d = depth
	}

	history, err := a.store.GetLocation(id, d)
	if err != nil {
		a.writeFailedResponse(w, http.StatusNotFound, err)
		return
	}

	resp := struct {
		OrderId string               `json:"order_id"`
		History []*location.Location `json:"history"`
	}{
		id,
		history,
	}

	a.writeResponse(w, http.StatusOK, resp)
}
func (a *Api) deleteLocationHistory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["orderId"]

	err := a.store.DeleteLocation(id)
	if err != nil {
		a.writeFailedResponse(w, http.StatusNotFound, err)
		return
	}
	a.writeSimpleResponse(w, http.StatusOK)
}
