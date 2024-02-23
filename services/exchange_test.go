package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetExchangeRateService(t *testing.T) {
	os.Setenv("COINBASE_SERVER_URL", "http://invalid.com")
	service, err := GetExchangeRateService()
	require.Error(t, err)
	require.Equal(t, err, ErrGetExchangeRates)
	require.Nil(t, service)
}

func Test_GetExchangeRateServiceWithMockResponse(t *testing.T) {
	rateMap := make(map[string]interface{})
	rateMap["00"] = "14.7601476014760148"
	rateMap["1INCH"] = "2.2988505747126437"
	exchangeData := &ExchangeRateData{
		RateData: Data{Currency: "USD", Rates: rateMap},
	}
	responseBytes, err := json.Marshal(exchangeData)
	require.NoError(t, err)
	testTable := []struct {
		name             string
		server           *httptest.Server
		expectedResponse *ExchangeRateData
		expectedErr      error
	}{
		{
			name:             "happy-server-rxesponse",
			server:           GetMockServer(responseBytes, false),
			expectedResponse: exchangeData,
			expectedErr:      nil,
		},
		{
			name: "invalid-data-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data":{"currency":USD,"rates":{"00":"14.7601476014760148","1INCH":"2.2988505747126437"}}}`))
			})),
			expectedResponse: nil,
			expectedErr:      ErrUnmarshallData,
		},
		{
			name: "missing-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"data":{"dollars":"USD","rates":{"00":"14.7601476014760148","1INCH":"2.2988505747126437"}}}`))
			})),
			expectedResponse: nil,
			expectedErr:      ErrMissingExchangeRates,
		},
		{
			name: "unmarshal-error",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			expectedResponse: nil,
			expectedErr:      ErrUnmarshallData,
		},
		{
			name:             "failure-status",
			server:           GetMockServer(responseBytes, true),
			expectedResponse: nil,
			expectedErr:      ErrGetExchangeRates,
		},
		{
			name: "response body error",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1")
				w.WriteHeader(http.StatusOK)
			})),
			expectedResponse: nil,
			expectedErr:      ErrInvalidResponse,
		},
	}

	t.Run(testTable[0].name, func(t *testing.T) {
		defer testTable[0].server.Close()
		os.Setenv("COINBASE_SERVER_URL", testTable[0].server.URL)
		resp, err := GetExchangeRateService()
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, resp.ExchangeRateData.RateData.Currency, testTable[0].expectedResponse.RateData.Currency)
	})

	t.Run(testTable[1].name, func(t *testing.T) {
		defer testTable[1].server.Close()
		os.Setenv("COINBASE_SERVER_URL", testTable[1].server.URL)
		resp, err := GetExchangeRateService()
		require.Equal(t, err, testTable[1].expectedErr)
		require.Nil(t, resp)
	})
	t.Run(testTable[2].name, func(t *testing.T) {
		defer testTable[2].server.Close()
		os.Setenv("COINBASE_SERVER_URL", testTable[2].server.URL)
		resp, err := GetExchangeRateService()
		require.Equal(t, err, testTable[2].expectedErr)
		require.Nil(t, resp)
	})
	t.Run(testTable[3].name, func(t *testing.T) {
		defer testTable[3].server.Close()
		os.Setenv("COINBASE_SERVER_URL", testTable[3].server.URL)
		resp, err := GetExchangeRateService()
		require.Equal(t, err, testTable[3].expectedErr)
		require.Nil(t, resp)
	})
	t.Run(testTable[4].name, func(t *testing.T) {
		defer testTable[4].server.Close()
		os.Setenv("COINBASE_SERVER_URL", testTable[4].server.URL)
		resp, err := GetExchangeRateService()
		require.Equal(t, testTable[4].expectedErr, err)
		require.Nil(t, resp)
	})
	t.Run(testTable[5].name, func(t *testing.T) {
		defer testTable[5].server.Close()
		os.Setenv("COINBASE_SERVER_URL", testTable[5].server.URL)
		resp, err := GetExchangeRateService()
		require.Equal(t, testTable[5].expectedErr, err)
		require.Nil(t, resp)
	})
}

func Test_GetCryptoRate(t *testing.T) {
	validRates := make(map[string]interface{})
	validRates["BTC"] = "0.013"
	validRates["ETH"] = "0.003"
	t.Run("rate found ", func(t *testing.T) {
		validRates := make(map[string]interface{})
		validRates["BTC"] = "0.013"
		validRates["ETH"] = "0.003"
		cryptoEval := &ExchangeRateService{
			ExchangeRateData: ExchangeRateData{
				RateData: Data{
					Currency: "USD",
					Rates:    validRates,
				},
			},
		}
		rate, err := cryptoEval.GetCryptoRate("BTC")
		require.Nil(t, err)
		require.Equal(t, rate, 0.013)
	})
	t.Run("rate found ", func(t *testing.T) {

		cryptoEval := &ExchangeRateService{
			ExchangeRateData: ExchangeRateData{
				RateData: Data{
					Currency: "USD",
					Rates:    validRates,
				},
			},
		}
		rate, err := cryptoEval.GetCryptoRate("BTC")
		require.Nil(t, err)
		require.Equal(t, rate, 0.013)
	})
	t.Run("invalid configuration ", func(t *testing.T) {
		cryptoEval := &ExchangeRateService{
			ExchangeRateData: ExchangeRateData{
				RateData: Data{
					Currency: "USD",
					Rates:    nil,
				},
			},
		}
		_, err := cryptoEval.GetCryptoRate("BTC")
		require.Equal(t, err, ErrIncorrectConfiguration)
	})
	t.Run("rate not found ", func(t *testing.T) {
		cryptoEval := &ExchangeRateService{
			ExchangeRateData: ExchangeRateData{
				RateData: Data{
					Currency: "USD",
					Rates:    validRates,
				},
			},
		}
		_, err := cryptoEval.GetCryptoRate("GBP")
		require.Contains(t, err.Error(), RateNotFound)
	})
	t.Run("invalid rate found ", func(t *testing.T) {
		invalidRates := make(map[string]interface{})
		invalidRates["BTC"] = "0.013"
		invalidRates["ETH"] = "0.003ghk"
		cryptoEval := &ExchangeRateService{
			ExchangeRateData: ExchangeRateData{
				RateData: Data{
					Currency: "USD",
					Rates:    invalidRates,
				},
			},
		}
		_, err := cryptoEval.GetCryptoRate("ETH")
		require.Contains(t, err.Error(), IncorrectRate)
	})
}
