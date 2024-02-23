package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
)

var (
	defaultServerUrl          = "https://api.coinbase.com/v2/exchange-rates?currency=USD"
	ErrGetExchangeRates       = fmt.Errorf("failed to get exchange rate data from coin base api")
	ErrUnmarshallData         = fmt.Errorf("failed to parse exchange rate data from coinbase api")
	ErrMissingExchangeRates   = fmt.Errorf("incorrect exchange rates from coinbase api")
	ErrIncorrectConfiguration = fmt.Errorf("invalid exchange rate service configuration")
	ErrInvalidResponse        = fmt.Errorf("invalid response from server")
	RateNotFound              = "exchange rate not found for crypto currency"
	IncorrectRate             = "invalid rate for crypto currency"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ExchangeRateService struct {
	ExchangeRateData ExchangeRateData
}

type Data struct {
	Currency string                 `json:"currency"`
	Rates    map[string]interface{} `json:"rates"`
}
type ExchangeRateData struct {
	RateData Data `json:"data"`
}

/*
*
Exchange service communicates with coinbase api to get data for USD
Env variable "COINBASE_SERVER_URL" can be set to different url incase we need exchange rates for different base currency .
Currently it fetches exchange rates for USD.
*
*/
func GetExchangeRateService() (*ExchangeRateService, error) {
	coinbaseService := ExchangeRateService{}
	serverUrl := os.Getenv("COINBASE_SERVER_URL")
	if serverUrl == "" {
		serverUrl = defaultServerUrl
	}
	resp, err := http.Get(serverUrl)
	if err != nil || resp.StatusCode != 200 {
		return nil, ErrGetExchangeRates
	}
	defer resp.Body.Close()
	respbytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrInvalidResponse
	}
	exchangeRates := ExchangeRateData{}
	if err = json.Unmarshal(respbytes, &exchangeRates); err != nil {
		return nil, ErrUnmarshallData
	}
	if exchangeRates.RateData.Currency == "" || exchangeRates.RateData.Rates == nil {
		return nil, ErrMissingExchangeRates
	}
	coinbaseService.ExchangeRateData = exchangeRates
	return &coinbaseService, nil
}

/*
*
Looksup currency to get the rate from exchange rate map.
*
*/
func (exchangeRateService *ExchangeRateService) GetCryptoRate(cryptoCurrency string) (float64, error) {
	if exchangeRateService.ExchangeRateData.RateData.Currency == "" || exchangeRateService.ExchangeRateData.RateData.Rates == nil {
		return 0, ErrIncorrectConfiguration
	}
	rateStr, ok := exchangeRateService.ExchangeRateData.RateData.Rates[cryptoCurrency]
	if !ok {
		return 0, fmt.Errorf("%s %s", RateNotFound, cryptoCurrency)
	}

	rate, err := strconv.ParseFloat(rateStr.(string), 64)
	if err != nil || rate == math.MaxFloat64 {
		return 0, fmt.Errorf("%s %s", IncorrectRate, cryptoCurrency)
	}
	return rate, nil
}

/*
*
Mock server for mocking different http reponse and statuses
*
*/
func GetMockServer(responseBytes []byte, withFailure bool) *httptest.Server {
	if !withFailure {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(responseBytes)
			if err != nil {
				panic(err)
			}
		}))
	} else {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
	}
}
