package cmc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
)

const (
	baseURL = "https://pro-api.coinmarketcap.com"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Service struct {
	apiKey string
	client httpClient
}

type Result struct {
	Status struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message,omitempty"`
	}
	Data []struct {
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

func NewService(client httpClient, apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
		client: client,
	}
}

func (c *Service) Run() {

}

func (c *Service) Conversion(ctx context.Context, from, to string, amount decimal.Decimal) (decimal.Decimal, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/v2/tools/price-conversion", http.NoBody)
	if err != nil {
		return decimal.Zero, err
	}

	query := request.URL.Query()
	query.Add("amount", amount.String())
	query.Add("symbol", from)
	query.Add("convert", to)
	request.URL.RawQuery = query.Encode()
	request.Header.Add("X-CMC_PRO_API_KEY", c.apiKey)

	response, err := c.client.Do(request)
	if err != nil {
		return decimal.Zero, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return decimal.Zero, errors.New("invalid response status code from cmc api")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return decimal.Zero, err
	}

	var result Result
	if err = json.Unmarshal(body, &result); err != nil {
		return decimal.Zero, err
	}

	if result.Status.ErrorCode > 0 && result.Status.ErrorMessage != "" {
		return decimal.Zero, fmt.Errorf("error from cmc api %s", result.Status.ErrorMessage)
	}

	if len(result.Data) == 0 {
		return decimal.Zero, errors.New("empty response data from cmc api")
	}

	quote, ok := result.Data[0].Quote[to]
	if !ok {
		return decimal.Zero, errors.New("invalid quote symbol")
	}

	return decimal.NewFromFloat(quote.Price), nil
}
