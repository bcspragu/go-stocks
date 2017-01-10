package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bcspragu/stocker/stock"
)

type TickerSymbol string

type Config struct {
	Holdings    map[TickerSymbol]float64
	TargetRatio map[TickerSymbol]int
}

type Holding struct {
	Ticker TickerSymbol
	Shares float64
}

type Fund struct {
	Ticker TickerSymbol
	Price  float64
}

var (
	holdings = flag.String("holdings", "", "file containing holding configuration")
	amount   = flag.Float64("amount", 0.0, "the amount to add to your funds")
)

func main() {
	flag.Parse()

	if *holdings == "" {
		log.Fatal("Need to specify --holdings filename")
	}

	if *amount == 0.0 {
		log.Fatal("Need to specify --amount to add")
	}

	config, err := loadConfig(*holdings)
	if err != nil {
		log.Fatal("Error loading holdings config: %v", err)
	}

	v := stock.NewVanguard()
	fChan := make(chan Fund)
	for ticker := range config.Holdings {
		go func(ticker TickerSymbol) {
			p, err := v.Price(string(ticker))
			if err != nil {
				log.Fatalf("Failed to retreive price for %s: %v", ticker, err)
			}
			fChan <- Fund{Ticker: ticker, Price: p}
		}(ticker)
	}

	funds := make(map[TickerSymbol]Fund)
	cashTotal := 0.0
	for i := 0; i < len(config.Holdings); i++ {
		fund := <-fChan
		funds[fund.Ticker] = fund
		cashTotal += config.Holdings[fund.Ticker] * fund.Price
	}
	close(fChan)

	ratioTotal := 0
	for _, ratioPart := range config.TargetRatio {
		ratioTotal += ratioPart
	}

	fmt.Println("You should buy:")
	for ticker, ratioPart := range config.TargetRatio {
		// Take the total amount we'll have and multiply it by the ratio we want to
		// get the amount we'd like to have. Subtrack how much we currently have to
		// get how much we should buy, then divide by the price of the stock
		fmt.Println(cashTotal, *amount, ratioPart, ratioTotal, float64(ratioPart)/float64(ratioTotal))
		shares := (((float64(cashTotal) + *amount) * float64(ratioPart) / float64(ratioTotal)) - config.Holdings[ticker]) / funds[ticker].Price / 100.0
		fmt.Printf("%s: %.3f shares\n", ticker, shares)
	}

}

func loadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
