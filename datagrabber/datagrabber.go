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

var mutex = &sync.RWMutex{}

// Poll new data each Quantum seconds
const Quantum = 60
const RemoteUrl = "https://api.exchangeratesapi.io/"

var lastRateUSD = 0.0
var lastRateGBP = 0.0
var lastWeekRateUSD = 0.0
var lastWeekRateGBP = 0.0
var lastPredictionGBP = false
var lastPredictionUSD = false

func init() {
	go func() {
		ticker := time.NewTicker(time.Second * Quantum)
		defer ticker.Stop()

		for {
			pollData();
			<-ticker.C
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Hour * 24)
		defer ticker.Stop()

		for {
			pollWeeklyData();
			<-ticker.C
		}
	}()
}

func pollWeeklyData() {
	mutex.Lock()
	defer mutex.Unlock()

	// Cool, it's another day, let's see what happened
	// last week.

	c1 := make(chan float64, 1)
	c2 := make(chan float64, 1)

	go func() {
		err := pollWeeklyRates("USD", c1);
		if (err != nil) {
			fmt.Println(err)
		}
	}()

	go func() {
		err := pollWeeklyRates("GBP", c2);
		if (err != nil) {
			fmt.Println(err)
		}
	}()

	lastWeekRateUSD = <-c1
	lastWeekRateGBP = <-c2

	go func() {
		close(c1)
		close(c2)
	}()
}

func pollData() {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("Polling data");

	c1 := make(chan float64, 1)
	c2 := make(chan float64, 1)

	go func() {
		err := pollCurrentRates("USD", c1);
		if (err != nil) {
			fmt.Println(err)
		}
	}()

	go func() {
		err := pollCurrentRates("GBP", c2);
		if (err != nil) {
			fmt.Println(err)
		}
	}()

	lastRateUSD = <-c1
	lastRateGBP = <-c2

	go func() {
		close(c1)
		close(c2)
	}()

	if (lastRateUSD >= lastWeekRateUSD) {
		lastPredictionUSD = false;
	} else {
		lastPredictionUSD = true;
	}

	if (lastRateGBP >= lastWeekRateGBP) {
		lastPredictionGBP = false;
	} else {
		lastPredictionGBP = true;
	}
}

// Assume the mutex is already locked when called
func pollCurrentRates(base string, out chan<- float64) error {
	url := RemoteUrl + "latest?symbols=EUR&base=" + base
	response := types.ExchangeLatestResponse{}

	err := makeGetRequest(url, &response)
	if (err != nil) {
		return err
	}

	out <- response.Rates["EUR"]
	return nil
}

// Assume the mutex is already locked when called
func pollWeeklyRates(base string, out chan<- float64) error {
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
		return err
	}

	if (len(response.Rates) < 5) {
		return errors.New("Not enough data to process")
	}

	avg := 0.0
	for _, value := range response.Rates {
		avg += value["EUR"]
	}

	out <- avg/float64(len(response.Rates))

	return nil
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
		return errors.New("Wrong StatusCode from server")
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(data); err != nil {
		return errors.New("Cannot decode data")
	}
	return nil
}

func GetData(ticker int) (*types.HTTPResponse, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	ret := types.HTTPResponse{}

	if (ticker == types.EUR_USD) {
		ret.Ticker = "USD"
		// Round to float32 so that we round the number.
		// It is useful to do calculations using float64
		// so that we have always higher precision.
		ret.Rate = float32(lastRateUSD)
		ret.WeekRate = float32(lastWeekRateUSD)
		ret.Prediction = lastPredictionUSD
	} else if (ticker == types.EUR_GBP) {
		ret.Ticker = "GBP"
		ret.Rate = float32(lastRateGBP)
		ret.WeekRate = float32(lastWeekRateGBP)
		ret.Prediction = lastPredictionGBP
	} else {
		return nil, errors.New("Wrong ticker")
	}

	return &ret, nil
}
