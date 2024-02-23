package main

import (
	"fmt"
	"math"

	"strconv"

	"github.com/slipika/cryto-evaluator/utils"
	"github.com/spf13/cobra"
)

func main() {
	var cmd = &cobra.Command{
		Use:     "./crypto-evaluator usd-amount currency1 currency2",
		Example: "./crypto-evaluator 100 BTC ETH",
		Long:    "CLI that takes in a USD amount as holdings, and calculates the 70/30 split for 2 given crypto currencies",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.ExactArgs(3)(cmd, args); err != nil {
				return err
			}
			amount, err := strconv.Atoi(args[0])
			if err != nil || amount > math.MaxInt {
				return fmt.Errorf("invalid amount specified: %s", args[0])
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			result, err := utils.ProcessSeventyThirtySplit(args[0], args[1], args[2])
			if err != nil {
				fmt.Printf("error: %s", err.Error())
			} else {
				fmt.Printf("\n%s", result[0])
				fmt.Printf("\n%s", result[1])
				fmt.Print("\n")
			}

		},
	}
	err := cmd.Execute()
	if err != nil {
		fmt.Println("exit with error")
	}
}
