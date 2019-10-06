package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"./datagrabber"
	"./types"
)

type responseErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TickerDataGrabber\n")
}

func GetTicker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	tickerStr := mux.Vars(r)["ticker"]

	ticker := types.INVALID_TICKER
	if strings.Compare(tickerStr, "USD") == 0 {
		ticker = types.EUR_USD
	} else if strings.Compare(tickerStr, "GBP") == 0 {
		ticker = types.EUR_GBP
	} else {
		replyError(http.StatusBadRequest, w, r, "Wrong ticker format")
		return
	}

	data, err := datagrabber.GetData(ticker)
	if err != nil {
		replyError(http.StatusInternalServerError, w, r, "Server in trouble")
		return
	}

	replyJSON(w, http.StatusOK, *data)
}

func replyJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func replyError(err int, w http.ResponseWriter, r *http.Request, text string) {
	replyJSON(w, err, map[string]string{"error": text})
}

func replyOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
