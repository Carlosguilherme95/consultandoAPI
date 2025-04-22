package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CurrencyInfo struct {
	Bid string `json:"bid"`
}
type CurrencyResponse struct {
	USDBRL CurrencyInfo `json:"USDBRL"`
}

func main() {
	//fazendo a chamada http
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	// decode do json
	var currency CurrencyResponse
	err = json.NewDecoder(req.Body).Decode(&currency)
	if err != nil {
		panic(err)
	}
	fmt.Println(currency.USDBRL.Bid)
}
