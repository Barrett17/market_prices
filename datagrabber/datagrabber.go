package datagrabber

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"../types"
)

// Poll new data each Quantum seconds
const Quantum = 60
const RemoteUrl = "https://api.exchangeratesapi.io/"

var mutex = &sync.Mutex{}

var lastRates = types.ExchangeLatestResponse{}
var lastWeekRateUSD = 0.0
var lastWeekRateGBP = 0.0
var lastDay = time.Time{}

func init() {
	go func() {
		ticker := time.NewTicker(time.Second * Quantum)
		defer ticker.Stop()

		for {
			pollData();
			<-ticker.C
		}
	}()
}

func pollData() {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("Polling data");

	err := pollCurrentRates();
	if (err != nil) {
		fmt.Println(err)
	}

	// Check that 24 hours have passed
	now := time.Now()
	days := now.Sub(lastDay).Hours() / 24
	if (days > 1) {
		lastDay = now
		err = pollWeeklyRates();
		if (err != nil) {
			fmt.Println(err)
		}
	}
}

// Assumes the mutex is already locked when called
func pollCurrentRates() error {
	req, _ := http.NewRequest("GET", RemoteUrl+"latest?symbols=USD,GBP", nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// TODO define errors
		return errors.New("Wrong StatusCode from server")
	}

	response := types.ExchangeLatestResponse{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return errors.New("Cannot decode data")
	}

	lastRates = response

	return nil
}

// Assumes the mutex is already locked when called
func pollWeeklyRates() error {
	dt := time.Now()
	today := dt.Format("2006-01-02")

	// Sat and Sun the market is down, keep in mind this
	dt = dt.AddDate(0, 0, -7)
	oneWeekAgo := dt.Format("2006-01-02")

	url := RemoteUrl+"history?start_at=" + oneWeekAgo + "&end_at=" +
		today + "&symbols=USD,GBP"

	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("Wrong StatusCode from server")
	}

	response := types.ExchangeHistoryResponse{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return errors.New("Cannot decode data")
	}

	if (len(response.Rates) < 5) {
		return errors.New("Not enough data to process")
	}

	avgUSD := 0.0
	avgGBP := 0.0
	for _, value := range response.Rates {
		avgUSD += value.USD
		avgGBP += value.GBP
	}

	avgUSD = avgUSD/float64(len(response.Rates))
	lastWeekRateUSD = avgUSD
	avgGBP = avgGBP/float64(len(response.Rates))
	lastWeekRateGBP = avgGBP

	return nil
}

func GetData(ticker int) (types.HTTPResponse, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var ret types.HTTPResponse

	if (ticker == types.EUR_USD) {
		ret.Ticker = "USD";
		ret.Rate = lastRates.Rates.USD;
		ret.WeekRate = lastWeekRateUSD;
	} else if (ticker == types.EUR_GBP) {
		ret.Ticker = "GBP";
		ret.Rate = lastRates.Rates.GBP;
		ret.WeekRate = lastWeekRateGBP;
	} else {
		return types.HTTPResponse{}, errors.New("Wrong ticker")
	}

	if (ret.Rate >= ret.WeekRate) {
		ret.Prediction = false;
	} else {
		ret.Prediction = true;
	}

	return ret, nil
}
