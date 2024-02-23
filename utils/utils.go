package utils

import (
	"strconv"

	"fmt"

	"github.com/slipika/cryto-evaluator/services"
)

/*
Calculated the seventy thirty split for given two currencies

Future improvements/refactor could take an env variable to define the percentage split for currency1 and currency2
*/
func ProcessSeventyThirtySplit(inputAmount, seventyCurrency, thirtyCurrency string) ([]string, error) {
	amount, err := strconv.ParseFloat(inputAmount, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount %s", inputAmount)
	}
	service, err := services.GetExchangeRateService()
	if err != nil {
		return nil, err
	}
	seventyMultiplier, err := service.GetCryptoRate(seventyCurrency)
	if err != nil {
		return nil, err
	}

	thirtyMultipler, err := service.GetCryptoRate(thirtyCurrency)
	if err != nil {
		return nil, err
	}

	split1 := 0.7 * amount
	split2 := 0.3 * amount

	result := make([]string, 2)
	result[0] = fmt.Sprintf("%d=>%.4f", int(split1), seventyMultiplier*split1)
	result[1] = fmt.Sprintf("%d=>%.4f", int(split2), thirtyMultipler*split2)
	return result, nil
}
