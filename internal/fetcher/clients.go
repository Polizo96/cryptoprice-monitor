package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BinanceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type CoinbaseResponse struct {
	Data struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

type KrakenResponse struct {
	Result map[string]struct {
		C []string `json:"c"`
	} `json:"result"`
}

func FetchBinancePrice(apiUrl string, coin string) (float64, error) {
	url := fmt.Sprintf(apiUrl, strings.ToUpper(coin))

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var br BinanceResponse
	if err := json.Unmarshal(body, &br); err != nil {
		return 0, err
	}

	var price float64
	_, err = fmt.Scanf(br.Price, "%f", &price)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func FetchCoinbasePrice(apiUrl string, coin string) (float64, error) {
	url := fmt.Sprintf(apiUrl, strings.ToUpper(coin))

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("coinbase API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var cr CoinbaseResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return 0, err
	}

	var price float64
	_, err = fmt.Sscanf(cr.Data.Amount, "%f", &price)
	if err != nil {
		return 0, err
	}

	return price, nil
}

func FetchKrakenPrice(apiUrl string, coin string) (float64, error) {
	pair := strings.ToUpper(coin) + "USD"
	url := fmt.Sprintf(apiUrl, pair)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("kraken API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var kr KrakenResponse
	if err := json.Unmarshal(body, &kr); err != nil {
		return 0, err
	}

	pairData, ok := kr.Result[pair]
	if !ok {
		return 0, fmt.Errorf("kraken pair %s not found", pair)
	}

	if len(pairData.C) == 0 {
		return 0, fmt.Errorf("kraken pair %s closing price missing", pair)
	}

	var price float64
	_, err = fmt.Sscanf(pairData.C[0], "%f", &price)
	if err != nil {
		return 0, err
	}

	return price, nil
}
