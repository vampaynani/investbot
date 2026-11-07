package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	inputPath := flag.String("input", "", "path to an assumptions JSON file")
	maxPriceAge := flag.Duration("max-price-age", 24*time.Hour, "reject a quote whose price.as_of is older than this")
	flag.Parse()
	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "usage: investor-metrics -input assumptions.json [-max-price-age 24h]")
		os.Exit(2)
	}

	file, err := os.Open(*inputPath)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	var in Input
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&in); err != nil {
		fatal(fmt.Errorf("decode input: %w", err))
	}
	if err := ensureEOF(decoder); err != nil {
		fatal(err)
	}

	result, err := Calculate(in, Options{Now: time.Now().UTC(), MaxPriceAge: *maxPriceAge})
	if err != nil {
		fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fatal(err)
	}
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("input must contain exactly one JSON value")
		}
		return fmt.Errorf("read trailing input: %w", err)
	}
	return nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "investor-metrics:", err)
	os.Exit(1)
}
