package converter

import (
	"context"
	"net"
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
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
			}).DialContext,
			MaxIdleConns:          1,
			MaxConnsPerHost:       1,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
	return &Usecase{
		conversionClient: cmc.NewService(conversionHTTPClient, apiKey),
	}
}

func (u *Usecase) Conversion(ctx context.Context, from, to string, amount decimal.Decimal) (decimal.Decimal, error) {
	return u.conversionClient.Conversion(ctx, from, to, amount)
}
