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

	// TODO check if a day has passed
	err = pollWeeklyRates();
	if (err != nil) {
		fmt.Println(err)
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

	s := fmt.Sprintf("%f %f", response.Rates.GBP, response.Rates.USD);
	fmt.Println(s)

	return nil
}

// Assumes the mutex is already locked when called
func pollWeeklyRates() error {
	req, _ := http.NewRequest("GET", RemoteUrl+"history?start_at=2018-01-01&end_at=2018-09-01&symbols=USD,GBP", nil)

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

	avgUSD := 0.0
	avgGBP := 0.0
	for _, value := range response.Rates {
		avgUSD += value.USD;
		avgGBP += value.GBP;
	}

	avgUSD = avgUSD/float64(len(response.Rates))
	lastWeekRateUSD = avgUSD
	avgGBP = avgGBP/float64(len(response.Rates))
	lastWeekRateGBP = avgGBP

	s := fmt.Sprintf("%f", lastWeekRateUSD);
	fmt.Println(s)

	return nil
}

func GetData() {
	mutex.Lock()
	defer mutex.Unlock()
	
}
