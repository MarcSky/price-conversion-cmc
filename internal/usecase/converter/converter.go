package converter

import (
	"context"
	"net/http"
	"time"

	"cryptoconverter/internal/service/cmc"

	"github.com/shopspring/decimal"
)

type Usecase struct {
	conversionClient ConversionClient
}

func New(apiKey string) *Usecase {
	conversionHTTPClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: http.DefaultTransport,
	}
	return &Usecase{
		conversionClient: cmc.NewService(conversionHTTPClient, apiKey),
	}
}

func (u *Usecase) Conversion(ctx context.Context, from, to string, amount decimal.Decimal) (decimal.Decimal, error) {
	return u.conversionClient.Conversion(ctx, from, to, amount)
}
