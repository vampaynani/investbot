package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Portfolio mirrors .investor_portfolio/positions.json. Unknown fields are
// rejected on read so that a rewrite can never silently drop data. Prices,
// shares, and weights are decimal strings/numbers handled with exact
// rational arithmetic, never floats.
type Portfolio struct {
	SchemaVersion int        `json:"schema_version"`
	UpdatedAt     string     `json:"updated_at"`
	Coverage      string     `json:"coverage"`
	Positions     []Position `json:"positions"`
}

type Position struct {
	Ticker      string       `json:"ticker"`
	Exchange    string       `json:"exchange"`
	Currency    string       `json:"currency"`
	Shares      string       `json:"shares"`
	AverageCost *string      `json:"average_cost"`
	Price       *string      `json:"price"`
	PriceAsOf   *string      `json:"price_as_of"`
	WeightPct   *json.Number `json:"weight_pct"`
	AsOf        string       `json:"as_of"`
}

// WeightsOptions carries the clock and policy so results are reproducible.
type WeightsOptions struct {
	Now             time.Time
	MaxPriceAgeDays int
	BaseCurrency    string            // required only when currencies are mixed
	FX              map[string]string // currency -> decimal rate into BaseCurrency
}

// WeightsReport is the deterministic outcome of one computation.
type WeightsReport struct {
	Status       string            `json:"status"` // "ok" or "unavailable"
	Reasons      []string          `json:"reasons,omitempty"`
	Currency     string            `json:"currency,omitempty"`
	Coverage     string            `json:"coverage,omitempty"`
	FX           map[string]string `json:"fx,omitempty"`
	TotalValue   string            `json:"total_market_value,omitempty"`
	Weights      []WeightLine      `json:"weights,omitempty"`
	WeightsNote  string            `json:"note,omitempty"`
	SumWeightPct string            `json:"sum_weight_pct,omitempty"`
}

type WeightLine struct {
	Ticker      string `json:"ticker"`
	Exchange    string `json:"exchange"`
	MarketValue string `json:"market_value"`
	WeightPct   string `json:"weight_pct"`
	PriceAsOf   string `json:"price_as_of"`
}

const weightsBasis = "weight = shares × price ÷ sum of shares × price over the listed positions × 100, exact decimal arithmetic, rounded half up to 2 decimals; weights are over listed positions only and may sum to 100.00 ± rounding"

// ComputeWeights fills WeightPct on every position of p in place. If any
// input needed for a trustworthy denominator is missing, invalid for the
// policy, or stale, all weights are set to nil and the report says why. It
// returns an error only for a malformed portfolio.
func ComputeWeights(p *Portfolio, opts WeightsOptions) (WeightsReport, error) {
	report := WeightsReport{Coverage: p.Coverage, WeightsNote: weightsBasis}
	clearWeights := func(reasons ...string) WeightsReport {
		for i := range p.Positions {
			p.Positions[i].WeightPct = nil
		}
		report.Status = "unavailable"
		report.Reasons = reasons
		report.WeightsNote = ""
		return report
	}

	if p.SchemaVersion != 1 {
		return report, fmt.Errorf("unsupported schema_version %d", p.SchemaVersion)
	}
	if p.Coverage != "complete" && p.Coverage != "partial" {
		return report, fmt.Errorf("coverage must be complete or partial, got %q", p.Coverage)
	}

	type valued struct {
		value *big.Rat
	}
	values := make([]valued, len(p.Positions))
	total := new(big.Rat)
	var reasons []string
	seen := map[string]bool{}

	// Resolve the single calculation currency.
	currencies := map[string]bool{}
	for _, pos := range p.Positions {
		currencies[strings.ToUpper(strings.TrimSpace(pos.Currency))] = true
	}
	calcCurrency := ""
	switch {
	case len(currencies) == 1:
		for c := range currencies {
			calcCurrency = c
		}
	case len(currencies) > 1:
		calcCurrency = strings.ToUpper(strings.TrimSpace(opts.BaseCurrency))
		if calcCurrency == "" {
			reasons = append(reasons, "positions use mixed currencies and no base currency was given")
		}
	}

	today := opts.Now.UTC().Truncate(24 * time.Hour)
	fxUsed := map[string]string{}

	for i, pos := range p.Positions {
		label := pos.Ticker + "/" + pos.Exchange
		key := strings.ToUpper(pos.Ticker) + "/" + strings.ToUpper(pos.Exchange)
		if pos.Ticker == "" || pos.Exchange == "" || pos.Currency == "" {
			return report, fmt.Errorf("position %d: ticker, exchange, and currency are required", i)
		}
		if seen[key] {
			return report, fmt.Errorf("duplicate position %s", label)
		}
		seen[key] = true

		shares, ok := parseDecimal(pos.Shares)
		if !ok || shares.Sign() <= 0 {
			return report, fmt.Errorf("%s: shares must be a positive decimal, got %q", label, pos.Shares)
		}
		if pos.Price == nil || pos.PriceAsOf == nil {
			reasons = append(reasons, label+": no verified quote (price/price_as_of missing)")
			continue
		}
		price, ok := parseDecimal(*pos.Price)
		if !ok || price.Sign() <= 0 {
			return report, fmt.Errorf("%s: price must be a positive decimal, got %q", label, *pos.Price)
		}
		quoteDate, err := time.Parse("2006-01-02", *pos.PriceAsOf)
		if err != nil {
			return report, fmt.Errorf("%s: price_as_of must be YYYY-MM-DD, got %q", label, *pos.PriceAsOf)
		}
		if quoteDate.After(today) {
			reasons = append(reasons, fmt.Sprintf("%s: price_as_of %s is in the future", label, *pos.PriceAsOf))
			continue
		}
		if ageDays := int(today.Sub(quoteDate).Hours() / 24); ageDays > opts.MaxPriceAgeDays {
			reasons = append(reasons, fmt.Sprintf("%s: quote from %s is %d days old (max %d)", label, *pos.PriceAsOf, ageDays, opts.MaxPriceAgeDays))
			continue
		}

		value := new(big.Rat).Mul(shares, price)
		if cur := strings.ToUpper(strings.TrimSpace(pos.Currency)); calcCurrency != "" && cur != calcCurrency {
			rateStr, found := opts.FX[cur]
			rate, ok := parseDecimal(rateStr)
			if !found || !ok || rate.Sign() <= 0 {
				reasons = append(reasons, fmt.Sprintf("%s: no valid FX rate %s→%s", label, cur, calcCurrency))
				continue
			}
			value.Mul(value, rate)
			fxUsed[cur] = rateStr
		}
		values[i] = valued{value: value}
		total.Add(total, value)
	}

	if len(reasons) > 0 {
		return clearWeights(reasons...), nil
	}
	if len(p.Positions) == 0 {
		report.Status = "ok"
		report.Currency = calcCurrency
		return report, nil
	}
	if total.Sign() <= 0 {
		return clearWeights("total market value is not positive"), nil
	}

	report.Status = "ok"
	report.Currency = calcCurrency
	if len(fxUsed) > 0 {
		report.FX = fxUsed
	}
	report.TotalValue = total.FloatString(6)
	sum := new(big.Rat)
	for i := range p.Positions {
		pct := roundHalfUp2(new(big.Rat).Quo(new(big.Rat).Mul(values[i].value, big.NewRat(100, 1)), total))
		sum.Add(sum, pct)
		text := pct.FloatString(2)
		n := json.Number(text)
		p.Positions[i].WeightPct = &n
		report.Weights = append(report.Weights, WeightLine{
			Ticker:      p.Positions[i].Ticker,
			Exchange:    p.Positions[i].Exchange,
			MarketValue: values[i].value.FloatString(6),
			WeightPct:   text,
			PriceAsOf:   *p.Positions[i].PriceAsOf,
		})
	}
	report.SumWeightPct = sum.FloatString(2)
	return report, nil
}

// parseDecimal accepts plain decimals only ("12", "0.2163", "-1.5"); it
// rejects exponents, fractions, and surrounding junk that big.Rat would take.
func parseDecimal(s string) (*big.Rat, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, false
	}
	digits, dot := 0, 0
	for i, r := range s {
		switch {
		case r >= '0' && r <= '9':
			digits++
		case r == '.':
			dot++
		case r == '-' && i == 0:
		default:
			return nil, false
		}
	}
	if digits == 0 || dot > 1 {
		return nil, false
	}
	r, ok := new(big.Rat).SetString(s)
	return r, ok
}

// roundHalfUp2 rounds a non-negative rational to two decimals, half up.
func roundHalfUp2(r *big.Rat) *big.Rat {
	scaled := new(big.Rat).Mul(r, big.NewRat(100, 1))
	scaled.Add(scaled, big.NewRat(1, 2))
	floor := new(big.Int).Quo(scaled.Num(), scaled.Denom())
	return new(big.Rat).SetFrac(floor, big.NewInt(100))
}

type fxFlags map[string]string

func (f fxFlags) String() string { return "" }
func (f fxFlags) Set(v string) error {
	k, rate, ok := strings.Cut(v, "=")
	if !ok || k == "" || rate == "" {
		return fmt.Errorf("-fx wants CCY=rate, got %q", v)
	}
	f[strings.ToUpper(strings.TrimSpace(k))] = strings.TrimSpace(rate)
	return nil
}

// runWeights implements `investor-metrics weights`. It exits non-zero only
// when the file cannot be read or is malformed; "unavailable" weights are a
// valid, reported outcome and are written as nulls.
func runWeights(args []string, now time.Time, stdout io.Writer) error {
	fs := flag.NewFlagSet("weights", flag.ContinueOnError)
	path := fs.String("file", "", "path to positions.json")
	write := fs.Bool("write", false, "write weight_pct back to the file (default: report only)")
	maxAge := fs.Int("max-price-age-days", 3, "reject a quote whose price_as_of is older than this many days")
	base := fs.String("base", "", "base currency, required only when positions use mixed currencies")
	fx := fxFlags{}
	fs.Var(fx, "fx", "FX rate CCY=rate (1 unit of CCY in base currency); repeatable")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("usage: investor-metrics weights -file positions.json [-write] [-max-price-age-days 3] [-base USD -fx MXN=0.05]")
	}

	raw, err := os.ReadFile(*path)
	if err != nil {
		return err
	}
	var p Portfolio
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&p); err != nil {
		return fmt.Errorf("decode %s: %w", *path, err)
	}
	if err := ensureEOF(dec); err != nil {
		return err
	}

	report, err := ComputeWeights(&p, WeightsOptions{Now: now, MaxPriceAgeDays: *maxAge, BaseCurrency: *base, FX: fx})
	if err != nil {
		return err
	}
	if *write {
		if err := writePortfolio(*path, &p); err != nil {
			return err
		}
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func writePortfolio(path string, p *Portfolio) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(p); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".positions-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(buf.Bytes()); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), info.Mode().Perm()); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
