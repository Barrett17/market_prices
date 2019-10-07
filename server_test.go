package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
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
		t.Errorf("Cannot GET data")
	}

	response := types.HTTPResponse{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		t.Errorf("Cannot decode data")
	}

	if strings.Compare(response.Ticker, "USD") != 0 {
		t.Errorf("Response has wrong ticker")
	}
}

func TestMultipleRequests(t *testing.T) {
	var wg sync.WaitGroup
	count := 1000
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(){
			req, _ := http.NewRequest("GET", url+"GBP", nil)

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("Req Error: " + err.Error())
			}

			defer res.Body.Close()
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestResponseTime(t *testing.T) {
	req, _ := http.NewRequest("GET", url+"USD", nil)

	client := &http.Client{}
	start := time.Now()
	res, err := client.Do(req)
	end := time.Now()
	if err != nil {
		t.Errorf("Req Error: " + err.Error())
	}

	defer res.Body.Close()
	elapsed := end.Sub(start)
	if (time.Duration(elapsed)/time.Millisecond > 1) {
		t.Errorf("Request over 1ms")
	}
}
