package main

import (
	"context"
	"cryptoconverter/internal/service/cmc"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cryptoconverter/internal/usecase/converter"
	"github.com/shopspring/decimal"
)

var help = "usage: ./converter <amount> <from> <to>"

// 123.45 USD BTC
func main() {
	errExit := func(err error) {
		fmt.Println(err)
		fmt.Println(help)
		os.Exit(1)
	}

	args, err := parseArgs(os.Args)
	if err != nil {
		errExit(err)
	}

	conversionHTTPClient := &http.Client{
		Timeout:   15 * time.Second,
		Transport: http.DefaultTransport,
	}

	cmcSvc := cmc.NewService(conversionHTTPClient, os.Getenv("CMC_API_KEY"))
	converterUs := converter.New(cmcSvc)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := converterUs.Conversion(ctx, args.From, args.To, args.Amount)
	if err != nil {
		errExit(err)
	}
	log.Println("result", result.String())

	os.Exit(0)
}

type Args struct {
	From   string
	To     string
	Amount decimal.Decimal
}

func parseArgs(args []string) (Args, error) {
	if len(args) != 4 {
		return Args{}, fmt.Errorf("invalid count arguments")
	}

	amount, err := decimal.NewFromString(args[1])
	if err != nil {
		return Args{}, err
	}

	return Args{
		Amount: amount,
		From:   args[2],
		To:     args[3],
	}, nil
}
