package cmc_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"cryptoconverter/internal/service/cmc"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (s *mockClient) Do(req *http.Request) (*http.Response, error) {
	return s.do(req)
}

func TestService_ConversionSuccess(t *testing.T) {
	mClient := &mockClient{
		do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"data":[{"id":1,"symbol":"BTC","name":"Bitcoin","amount":1,"quote":{"USD":{"price":25652.814685744266}}}]}`)),
			}, nil
		},
	}

	svc := cmc.NewService(mClient, "")
	res, err := svc.Conversion(context.TODO(), "BTC", "USD", decimal.NewFromInt(1))
	assert.NoError(t, err)
	assert.Equal(t, "25652.814685744266", res.String())
}

func TestService_ConversionFailure(t *testing.T) {
	mClient := &mockClient{
		do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":{"error_code":1001,"error_message":"This API Key is invalid."}}`)),
			}, nil
		},
	}

	svc := cmc.NewService(mClient, "")
	res, err := svc.Conversion(context.TODO(), "BTC", "USD", decimal.NewFromInt(1))
	assert.Equal(t, decimal.Zero, res)
	assert.Error(t, err)
}
