package datagrabber

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Poll new data each 60 seconds
const Quantum = 60
const Timeout = Quantum
const RemoteUrl = "https://api.exchangeratesapi.io/"

var mutex = &sync.Mutex{}
var lastRates = ExchangeLatestResponse{}
var lastWeekRate float64

type Rates struct {
	USD float64 `json:"USD"`
	GBP float64 `json:"GBP"`
}

type ExchangeLatestResponse struct {
	Rates Rates `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}

type ExchangeHistoryResponse struct {
	Rates map[string]Rates `json:"rates"`
}

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

	response := ExchangeLatestResponse{}

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

	response := ExchangeHistoryResponse{}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		return errors.New("Cannot decode data")
	}

	//lastWeekRates = response

	//fmt.Println(response.Rates)

	return nil
}

func GetData() {
	mutex.Lock()
	defer mutex.Unlock()

	
}
