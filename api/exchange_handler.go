package api

import (
	"encoding/json"
	"fmt"
	"go-concurrency/types"
	"io"
	"net/http"
	"sync"
)

type ExchangeHandler struct {
	sync.Mutex
	Currencies map[string]types.Currency
}

func NewExchangeHandler(c map[string]types.Currency) *ExchangeHandler {
	return &ExchangeHandler{
		Currencies: c,
	}
}

func (ce *ExchangeHandler) HandleGetAllCurrencies() error {
	resp, err := http.Get("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies.json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	csMap := make(map[string]string)
	err = json.Unmarshal(cs, &csMap)
	if err != nil {
		return err
	}

	i := 0
	for code, name := range csMap {
		if i > 100 {
			break
		}
		c := types.Currency{
			Code:  code,
			Name:  name,
			Rates: make(map[string]float64),
		}
		ce.Currencies[code] = c
		i++
	}
	return nil
}

func (ce *ExchangeHandler) FetchCurrencyRates(currencyCode string) (map[string]float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@latest/v1/currencies/%s.json", currencyCode))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rates, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ratesStruct := make(map[string]interface{})
	err = json.Unmarshal(rates, &ratesStruct)
	if err != nil {
		return nil, err
	}

	//convert to map[string]float64
	ratesFloat := make(map[string]float64)
	for code, rate := range ratesStruct[currencyCode].(map[string]interface{}) {
		ratesFloat[code] = float64(rate.(float64))
	}
	return ratesFloat, nil
}
