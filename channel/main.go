package main

import (
	"fmt"
	"go-concurrency/api"
	"go-concurrency/types"
	"time"
)

func runCurrencyWorker(
	workerId int,
	currencyChan <-chan types.Currency,
	resultChan chan<- types.Currency) {
	currencies := make(map[string]types.Currency)
	exchangeHandler := api.NewExchangeHandler(currencies)
	fmt.Printf("Worker %d started\n", workerId)
	for c := range currencyChan {
		rates, err := exchangeHandler.FetchCurrencyRates(c.Code)
		if err != nil {
			panic(err)
		}
		c.Rates = rates
		resultChan <- c
	}
	fmt.Printf("Worker %d stopped", workerId)
}

func main() {
	currencies := make(map[string]types.Currency)
	exchangeHandler := api.NewExchangeHandler(currencies)

	err := exchangeHandler.HandleGetAllCurrencies()
	if err != nil {
		panic(err)
	}

	currencyChan := make(chan types.Currency, len(exchangeHandler.Currencies))
	resultChan := make(chan types.Currency, len(exchangeHandler.Currencies))
	for i := 0; i < 5; i++ {
		go runCurrencyWorker(i, currencyChan, resultChan)
	}
	startTime := time.Now()
	resultCount := 0

	for _, curr := range exchangeHandler.Currencies {
		currencyChan <- curr
	}

	for {
		if resultCount == len(exchangeHandler.Currencies) {
			fmt.Println("Closing resultChain")
			close(currencyChan)
			break
		}
		select {
		case currency := <-resultChan:
			exchangeHandler.Currencies[currency.Code] = currency
			resultCount++
		case <-time.After(3 * time.Second):
			fmt.Println("Timeout")
			break
		}
	}
	endTime := time.Now()
	fmt.Println("======== RESULTS ========")
	for _, curr := range exchangeHandler.Currencies {
		fmt.Printf("%s (%s): %d rates\n", curr.Name, curr.Code, len(curr.Rates))
	}
	fmt.Println("=======================")
	fmt.Println("Time taken: ", endTime.Sub(startTime))
}
