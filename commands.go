package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type Response struct {
	Ticker     string `json:"ticker"`
	Rate       string `json:"rate"`
	WeekRate   string `json:"weekrate"`
	Prediction string `json:"prediction"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TickerDataGrabber\n")
}

func GetTicker(w http.ResponseWriter, r *http.Request) {
	
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
	w.Header().Set("Content-Type",
		"application/json")
	w.WriteHeader(http.StatusOK)
}
