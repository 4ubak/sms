package httpapi

import (
	"github.com/xakpro/sms/internal/domain/entities"
	"net/http"
)

func (a *API) hSend(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.SendReqSt{}

	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	if reqObj.Sync {
		sendErr := a.cr.Send(reqObj)
		if sendErr != nil {
			a.uHandleError(sendErr, w)
			return
		}
	} else {
		go func() { _ = a.cr.Send(reqObj) }()
	}

	w.WriteHeader(200)
}

func (a *API) hCall(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.SendReqSt{}

	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	code, sendErr := a.cr.Call(reqObj)
	if sendErr != nil {
		a.uHandleError(sendErr, w)
		return
	}

	a.uRespondJSON(w, map[string]string{
		"code": code,
	})
}

func (a *API) hBcast(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.SendReqSt{}

	if !a.uParseRequestJSON(w, r, reqObj) {
		return
	}

	if reqObj.Sync {
		sendErr := a.cr.Bcast(reqObj)
		if sendErr != nil {
			a.uHandleError(sendErr, w)
			return
		}
	} else {
		go func() { _ = a.cr.Bcast(reqObj) }()
	}

	w.WriteHeader(200)
}

func (a *API) hGetBalance(w http.ResponseWriter, r *http.Request) {
	balance := a.cr.GetBalance()

	a.uRespondJSON(w, map[string]float64{
		"balance": balance,
	})
}

func (a *API) hCronCheckBalance(w http.ResponseWriter, r *http.Request) {
	a.cr.CheckBalance()

	w.WriteHeader(200)
}
