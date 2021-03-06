package httpapi

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (a *API) router() http.Handler {
	r := mux.NewRouter()

	r.PathPrefix("/send").HandlerFunc(a.hSend).Methods("POST")
	r.PathPrefix("/call").HandlerFunc(a.hCall).Methods("POST")
	r.PathPrefix("/bcast").HandlerFunc(a.hBcast).Methods("POST")
	r.PathPrefix("/balance").HandlerFunc(a.hGetBalance).Methods("GET")
	r.PathPrefix("/cron/check_balance").HandlerFunc(a.hCronCheckBalance).Methods("GET")

	return r
}
