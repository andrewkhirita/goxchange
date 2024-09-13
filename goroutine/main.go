package main

import (
	"fmt"
	"go-concurrency/api"
	"go-concurrency/types"
	"sync"
	"time"
)

func main() {
	currencies := make(map[string]types.Currency)
	exchangeHandler := api.NewExchangeHandler(currencies)

	err := exchangeHandler.HandleGetAllCurrencies()
	if err != nil {
		panic(err)
	}

	startTime := time.Now()

	go func() {
		for {
			exchangeHandler.Lock()
			usd, ok := exchangeHandler.Currencies["usd"]
			exchangeHandler.Unlock()
			if ok {
				fmt.Println("USD:", usd.Rates)
			}
		}
	}()
	wg := sync.WaitGroup{}
	for code := range exchangeHandler.Currencies {
		wg.Add(1)
		go func(code string) {
			rates, err := exchangeHandler.FetchCurrencyRates(code)
			if err != nil {
				panic(err)
			}
			exchangeHandler.Lock()
			exchangeHandler.Currencies[code] = types.Currency{
				Code:  code,
				Name:  exchangeHandler.Currencies[code].Name,
				Rates: rates,
			}
			exchangeHandler.Unlock()
			wg.Done()
		}(code)
	}
	wg.Wait()

	endTime := time.Now()
	fmt.Println("======== RESULTS ========")
	for _, curr := range exchangeHandler.Currencies {
		fmt.Printf("%s (%s): %d rates\n", curr.Name, curr.Code, len(curr.Rates))
	}
	fmt.Println("=======================")
	fmt.Println("Time taken: ", endTime.Sub(startTime))
}
