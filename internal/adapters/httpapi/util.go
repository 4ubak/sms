package httpapi

import (
	"encoding/json"
	"github.com/xakpro/sms/internal/errs"
	"log"
	"net/http"
)

func (a *API) uSetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (a *API) uRespondJSON(w http.ResponseWriter, obj interface{}) {
	a.uSetContentTypeJSON(w)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Panicln("Fail to encode json obj", err)
	}
}

func (a *API) uHandleError(err error, w http.ResponseWriter) {
	if err != nil {
		a.uSetContentTypeJSON(w)
		w.WriteHeader(http.StatusOK)
		switch cErr := err.(type) {
		case *errs.Err:
			a.uRespondJSON(w, ErrRepSt{
				ErrorCode: cErr.Error(),
			})
		default:
			a.uRespondJSON(w, ErrRepSt{
				ErrorCode: errs.ServerNA.Error(),
			})
		}
	}
}

func (a *API) uParseRequestJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		a.uHandleError(err, w)
		return false
	}
	return true
}
