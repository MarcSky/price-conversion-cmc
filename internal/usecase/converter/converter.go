package converter

import (
	"context"

	"github.com/shopspring/decimal"
)

type Usecase struct {
	conversionClient ConversionClient
}

func New(conversionClient ConversionClient) *Usecase {
	return &Usecase{
		conversionClient: conversionClient,
	}
}

func (u *Usecase) Conversion(ctx context.Context, from, to string, amount decimal.Decimal) (decimal.Decimal, error) {
	return u.conversionClient.Conversion(ctx, from, to, amount)
}
