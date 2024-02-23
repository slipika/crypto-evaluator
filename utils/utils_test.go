package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/slipika/cryto-evaluator/services"
	"github.com/stretchr/testify/require"
)

func Test_ProcessSeventyThirtySplit(t *testing.T) {

	rateMap := make(map[string]interface{})
	max := math.MaxFloat64
	rateMap["00"] = "14.7601476014760148"
	rateMap["1INCH"] = "2.2988505747126437"
	rateMap["MAX"] = fmt.Sprintf("%f", max)

	exchangeData := &services.ExchangeRateData{
		RateData: services.Data{Currency: "USD", Rates: rateMap},
	}
	responseBytes, err := json.Marshal(exchangeData)
	require.NoError(t, err)
	server := services.GetMockServer(responseBytes, false)
	defer server.Close()
	os.Setenv("COINBASE_SERVER_URL", server.URL)
	t.Run("success", func(t *testing.T) {
		result, err := ProcessSeventyThirtySplit("100", "00", "1INCH")
		require.Equal(t, result[0], "70=>1033.2103")
		require.Equal(t, result[1], "30=>68.9655")
		require.NoError(t, err)
	})
	t.Run("invalid amount", func(t *testing.T) {
		result, err := ProcessSeventyThirtySplit("100f", "00", "1INCH")
		require.Nil(t, result)
		require.Contains(t, err.Error(), "invalid amount")
	})
	t.Run("currency not found fro seventy", func(t *testing.T) {
		result, err := ProcessSeventyThirtySplit("100", "8900", "1INCH")
		require.Nil(t, result)
		require.Contains(t, err.Error(), services.RateNotFound)
	})
	t.Run("currency not found for thirty", func(t *testing.T) {
		result, err := ProcessSeventyThirtySplit("100", "00", "2345")
		require.Nil(t, result)
		require.Contains(t, err.Error(), services.RateNotFound)
	})
	t.Run("currency not found for thirty", func(t *testing.T) {
		result, err := ProcessSeventyThirtySplit("100", "00", "MAX")
		require.Nil(t, result)
		require.Contains(t, err.Error(), services.IncorrectRate)
	})
	t.Run("service error", func(t *testing.T) {
		failedExchangeServer := services.GetMockServer(responseBytes, true)
		defer failedExchangeServer.Close()
		os.Setenv("COINBASE_SERVER_URL", failedExchangeServer.URL)
		result, err := ProcessSeventyThirtySplit("100", "00", "MAX")
		require.Nil(t, result)
		require.Equal(t, err, services.ErrGetExchangeRates)
	})
}
