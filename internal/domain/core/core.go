package core

import (
	"bytes"
	"encoding/json"
	"github.com/xakpro/sms/internal/constants"
	"github.com/xakpro/sms/internal/domain/entities"
	"github.com/xakpro/sms/internal/errs"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const balanceCacheTimeout = 20 * time.Minute

var (
	smscHttpClient          = &http.Client{Timeout: 15 * time.Second}
	balanceNotifyHttpClient = &http.Client{Timeout: 15 * time.Second}
)

func (c *Core) Send(pars *entities.SendReqSt) error {
	var err error

	err = c.validateValues(pars, false)
	if err != nil {
		return err
	}

	urlValues := url.Values{
		"login":   {c.smscUsername},
		"psw":     {c.smscPassword},
		"phones":  {pars.Phones},
		"mes":     {pars.Message},
		"charset": {"utf-8"},
		"fmt":     {"3"},
	}

	if c.smscSender != "" {
		urlValues.Add("sender", c.smscSender)
	}

	urlString := constants.SMSCUrlPrefix + "send.php?" + urlValues.Encode()

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return errs.ServerNA
	}

	rep, err := smscHttpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err)
		return errs.ServerNA
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return errs.ServerNA
	}

	resultObj := &entities.SendRepSt{}

	if err = json.NewDecoder(rep.Body).Decode(resultObj); err != nil {
		c.lg.Errorw("Fail to parse http-body", err)
		return errs.ServerNA
	}

	if (resultObj.ErrorCode != 0) || (resultObj.Error != "") {
		c.lg.Infow("User phone", "phone", pars.Phones)
		if resultObj.ErrorCode != 8 && resultObj.ErrorCode != 7 && resultObj.ErrorCode != 6 {
			c.lg.Errorw("Bad response smsc.kz", nil, "error_code", resultObj.ErrorCode, "error", resultObj.Error)
		}
		return errs.ServerNA
	}

	return nil
}

func (c *Core) Call(pars *entities.SendReqSt) (string, error) {
	var err error

	err = c.validateValues(pars, true)
	if err != nil {
		return "", err
	}

	urlValues := url.Values{
		"login":  {c.smscUsername},
		"psw":    {c.smscPassword},
		"phones": {pars.Phones},
		"mes":    {"code"},
		"call":   {"1"},
		"fmt":    {"3"},
	}

	urlString := constants.SMSCUrlPrefix + "send.php?" + urlValues.Encode()

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return "", errs.ServerNA
	}

	rep, err := smscHttpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err)
		return "", errs.ServerNA
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return "", errs.ServerNA
	}

	resultObj := &entities.SendRepSt{}

	c.lg.Infow("infow", "Body", rep.Body)

	if err = json.NewDecoder(rep.Body).Decode(resultObj); err != nil {
		c.lg.Errorw("Fail to parse http-body", err)
		return "", errs.ServerNA
	}

	if (resultObj.ErrorCode != 0) || (resultObj.Error != "") {
		c.lg.Infow("User phone", "phone", pars.Phones)
		if resultObj.ErrorCode != 8 && resultObj.ErrorCode != 7 && resultObj.ErrorCode != 6 {
			c.lg.Errorw("Bad response smsc.kz", nil, "error_code", resultObj.ErrorCode, "error", resultObj.Error)
		}
		return "", errs.ServerNA
	}

	return resultObj.CODE, nil
}

func (c *Core) Bcast(pars *entities.SendReqSt) error {
	var err error

	err = c.validateValues(pars, false)
	if err != nil {
		c.lg.Errorw("Not correct values", err)
		return err
	}

	urlValues := url.Values{
		"add":     {"1"},
		"login":   {c.smscUsername},
		"psw":     {c.smscPassword},
		"name":    {"bcast"},
		"phones":  {pars.Phones},
		"mes":     {pars.Message},
		"charset": {"utf-8"},
		"fmt":     {"3"},
	}

	if c.smscSender != "" {
		urlValues.Add("sender", c.smscSender)
	}

	urlString := constants.SMSCUrlPrefix + "send.php?" + urlValues.Encode()

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return errs.ServerNA
	}

	rep, err := smscHttpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err)
		return errs.ServerNA
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return errs.ServerNA
	}

	repObj := &entities.SendRepSt{}

	if err = json.NewDecoder(rep.Body).Decode(repObj); err != nil {
		c.lg.Errorw("Fail to parse http-body", err)
		return errs.ServerNA
	}

	if (repObj.ErrorCode != 0) || (repObj.Error != "") {
		c.lg.Infow("User phone", "phone", pars.Phones)
		if repObj.ErrorCode != 8 && repObj.ErrorCode != 7 && repObj.ErrorCode != 6 {
			c.lg.Errorw("Bad response smsc.kz", nil, "error_code", repObj.ErrorCode, "error", repObj.Error)
		}
		return errs.ServerNA
	}

	return nil
}

func (c *Core) GetBalance() float64 {
	var cacheObj entities.BalanceCacheSt

	cacheObjRaw, ok, err := c.cache.Get(constants.BalanceCacheKey)
	if err != nil {
		return 0
	}
	if !ok {
		return 0
	}

	err = json.Unmarshal(cacheObjRaw, &cacheObj)
	if err != nil {
		c.lg.Errorw("Fail to unmarshal json", err)
		return 0
	}

	return cacheObj.Balance
}

func (c *Core) setBalance(value float64) error {
	cacheObj := entities.BalanceCacheSt{
		Balance: value,
	}

	cacheObjRaw, err := json.Marshal(cacheObj)
	if err != nil {
		c.lg.Errorw("Fail to marshal json", err)
		return err
	}

	err = c.cache.Set(constants.BalanceCacheKey, cacheObjRaw, balanceCacheTimeout)
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) NotifyBalance(url string, balance float64) {
	reqBodyJson, err := json.Marshal(map[string]float64{"balance": balance})
	if err != nil {
		c.lg.Errorw("Fail to marshal json", err, "balance", balance)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return
	}

	rep, err := balanceNotifyHttpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err)
		return
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return
	}
}

func (c *Core) requestBalance() (float64, error) {
	urlValues := url.Values{
		"login": {c.smscUsername},
		"psw":   {c.smscPassword},
		"fmt":   {"3"},
	}.Encode()

	urlString := constants.SMSCUrlPrefix + "balance.php?" + urlValues

	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		c.lg.Errorw("Fail to create http-request", err)
		return 0, errs.ServerNA
	}

	rep, err := smscHttpClient.Do(req)
	if err != nil {
		c.lg.Errorw("Fail to send http-request", err)
		return 0, errs.ServerNA
	}
	defer rep.Body.Close()

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		c.lg.Errorw("Fail to send http-request, bad status code", nil, "status_code", rep.StatusCode)
		return 0, errs.ServerNA
	}

	repObj := &entities.GetBalanceRepSt{}

	if err = json.NewDecoder(rep.Body).Decode(repObj); err != nil {
		c.lg.Errorw("Fail to parse http-body", err)
		return 0, errs.ServerNA
	}

	if (repObj.ErrorCode != 0) || (repObj.Error != "") {
		c.lg.Errorw("Bad response smsc.kz", nil, "error_code", repObj.ErrorCode, "error", repObj.Error)
		return 0, errs.ServerNA
	}

	result, _ := strconv.ParseFloat(repObj.Balance, 64)

	return result, nil
}

func (c *Core) validateValues(pars *entities.SendReqSt, skipMsg bool) error {
	if len(pars.Phones) == 0 {
		c.lg.Warnw("Phones is empty", errs.PhonesRequired)
		return errs.PhonesRequired
	}

	if !skipMsg {
		if len(pars.Message) == 0 {
			c.lg.Warnw("Message is empty", errs.MessageRequired)
			return errs.MessageRequired
		}
	}

	return nil
}

func (c *Core) CheckBalance() {
	var err error

	currentBalance := c.GetBalance()

	newBalance, err := c.requestBalance()
	if err != nil {
		return
	}

	c.lg.Infow("Balance checked", "balance", newBalance)

	for vBalance, vUrl := range c.balanceNotifUrls {
		if newBalance < vBalance && currentBalance > vBalance {
			c.NotifyBalance(vUrl, newBalance)
		}
	}

	err = c.setBalance(newBalance)
	if err != nil {
		return
	}
}
