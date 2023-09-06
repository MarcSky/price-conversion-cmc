package converter

import (
	"context"

	"github.com/shopspring/decimal"
)

type ConversionClient interface {
	Conversion(ctx context.Context, from, to string, amount decimal.Decimal) (decimal.Decimal, error)
}
