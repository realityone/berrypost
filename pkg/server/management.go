package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Management struct {
}

func (m Management) intro(w http.ResponseWriter, req *http.Request) {
	introSchema := struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}{
		Name:    "berrypost-management-server",
		Version: "0.0.1",
	}
	json.NewEncoder(w).Encode(introSchema)
}

func (m Management) SetupRoute(in *mux.Route) {
	in.Path("/api/_intro").Methods("GET").HandlerFunc(m.intro)
}
