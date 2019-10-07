package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Barrett17/market_prices/types"
)

var url = "http://127.0.0.1:8080/api/ticker/"

func TestMain(t *testing.T) {
	req, _ := http.NewRequest("GET", url+"USD", nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Req Error: " + err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Cannot GET data!")
	}

	response := types.HTTPResponse{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Errorf("Cannot decode data!")
	}

	if strings.Compare(response.Ticker, "USD") != 0 {
		t.Errorf("Response has wrong ticker!")
	}
}
