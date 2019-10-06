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

var lastRateUSD = 0.0
var lastRateGBP = 0.0
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

	rate, err := pollCurrentRates("USD");
	if (err != nil) {
		fmt.Println(err)
	}
	lastRateUSD = rate

	rate, err = pollCurrentRates("GBP");
	if (err != nil) {
		fmt.Println(err)
	}
	lastRateGBP = rate

	// Check that at least 24 hours have passed
	now := time.Now()
	days := now.Sub(lastDay).Hours() / 24
	if (days >= 1) {
		// Cool, it's another day, let's see what happened
		// last week.
		lastDay = now

		rate, err = pollWeeklyRates("USD");
		if (err != nil) {
			fmt.Println(err)
		}
		lastWeekRateUSD = rate

		rate, err = pollWeeklyRates("GBP");
		if (err != nil) {
			fmt.Println(err)
		}
		lastWeekRateGBP = rate
	}
}

// Assume the mutex is already locked when called
func pollCurrentRates(base string) (float64, error) {
	url := RemoteUrl + "latest?symbols=EUR&base=" + base
	response := types.ExchangeLatestResponse{}

	err := makeGetRequest(url, &response)
	if (err != nil) {
		return 0, err
	}

	return response.Rates["EUR"], nil
}

// Assume the mutex is already locked when called
func pollWeeklyRates(base string) (float64, error) {
	dt := time.Now()
	today := dt.Format("2006-01-02")

	// Sat and Sun the market is down, keep in mind this
	dt = dt.AddDate(0, 0, -7)
	oneWeekAgo := dt.Format("2006-01-02")

	url := RemoteUrl+"history?start_at=" + oneWeekAgo + "&end_at=" +
		today + "&symbols=EUR&base=" + base

	response := types.ExchangeHistoryResponse{}

	err := makeGetRequest(url, &response)
	if (err != nil) {
		return 0, err
	}

	if (len(response.Rates) < 5) {
		return 0, errors.New("Not enough data to process")
	}

	avg := 0.0
	for _, value := range response.Rates {
			avg += value["EUR"]
	}

	return avg/float64(len(response.Rates)), nil
}

func makeGetRequest(url string, data interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)

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

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(data); err != nil {
		return errors.New("Cannot decode data")
	}
	return nil
}

func GetData(ticker int) (types.HTTPResponse, error) {
	mutex.Lock()
	defer mutex.Unlock()

	ret := types.HTTPResponse{}

	if (ticker == types.EUR_USD) {
		ret.Ticker = "USD";
		// Round to float32 so that we round the number.
		// It is useful to do calculations using float64
		// so that we have always higher precision.
		ret.Rate = float32(lastRateUSD);
		ret.WeekRate = float32(lastWeekRateUSD);
	} else if (ticker == types.EUR_GBP) {
		ret.Ticker = "GBP";
		ret.Rate = float32(lastRateGBP);
		ret.WeekRate = float32(lastWeekRateGBP);
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
