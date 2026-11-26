package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var weightsNow = time.Date(2026, 11, 26, 15, 0, 0, 0, time.UTC)

func sp(s string) *string { return &s }

func pos(ticker, ccy, shares, price, date string) Position {
	p := Position{Ticker: ticker, Exchange: "NASDAQ", Currency: ccy, Shares: shares, AsOf: "2026-11-25"}
	if price != "" {
		p.Price, p.PriceAsOf = sp(price), sp(date)
	}
	return p
}

func opts() WeightsOptions { return WeightsOptions{Now: weightsNow, MaxPriceAgeDays: 3} }

func weightsOf(p *Portfolio) []string {
	var out []string
	for _, x := range p.Positions {
		if x.WeightPct == nil {
			out = append(out, "null")
		} else {
			out = append(out, x.WeightPct.String())
		}
	}
	return out
}

func TestWeightsExactAndDeterministic(t *testing.T) {
	mk := func() *Portfolio {
		return &Portfolio{SchemaVersion: 1, Coverage: "complete", Positions: []Position{
			pos("A", "USD", "1", "1", "2026-11-26"),
			pos("B", "USD", "1", "1", "2026-11-26"),
			pos("C", "USD", "1", "1", "2026-11-26"),
		}}
	}
	p1, p2 := mk(), mk()
	r1, err := ComputeWeights(p1, opts())
	if err != nil || r1.Status != "ok" {
		t.Fatalf("status=%s err=%v", r1.Status, err)
	}
	r2, _ := ComputeWeights(p2, opts())
	if !reflect.DeepEqual(r1, r2) || !reflect.DeepEqual(weightsOf(p1), weightsOf(p2)) {
		t.Fatal("same inputs produced different outputs")
	}
	if got := weightsOf(p1); !reflect.DeepEqual(got, []string{"33.33", "33.33", "33.33"}) {
		t.Fatalf("weights = %v", got)
	}
	if r1.SumWeightPct != "99.99" {
		t.Fatalf("sum = %s", r1.SumWeightPct)
	}
}

func TestWeightsRoundHalfUp(t *testing.T) {
	// 1/8 of 100 = 12.5 exactly; 0.125 of 100 -> use 1:7 split for 12.5/87.5.
	p := &Portfolio{SchemaVersion: 1, Coverage: "complete", Positions: []Position{
		pos("A", "USD", "1", "1", "2026-11-26"),
		pos("B", "USD", "7", "1", "2026-11-26"),
	}}
	ComputeWeights(p, opts())
	if got := weightsOf(p); !reflect.DeepEqual(got, []string{"12.50", "87.50"}) {
		t.Fatalf("weights = %v", got)
	}
	// 0.5% boundary: 1 of 200 -> 0.5 -> "0.50"; 1 of 800 -> 0.125 -> "0.13" (half up).
	q := &Portfolio{SchemaVersion: 1, Coverage: "complete", Positions: []Position{
		pos("A", "USD", "1", "1", "2026-11-26"),
		pos("B", "USD", "799", "1", "2026-11-26"),
	}}
	ComputeWeights(q, opts())
	if got := weightsOf(q); got[0] != "0.13" {
		t.Fatalf("0.125 should round half up to 0.13, got %v", got)
	}
}

func TestWeightsFractionalSharesAndPrices(t *testing.T) {
	p := &Portfolio{SchemaVersion: 1, Coverage: "complete", Positions: []Position{
		pos("NVDA", "USD", "0.21633", "180.25", "2026-11-26"),
		pos("META", "USD", "0.03190", "600", "2026-11-25"),
		pos("MSFT", "USD", "0.03902", "500.5", "2026-11-24"),
	}}
	r, _ := ComputeWeights(p, opts())
	// Expected values computed independently with Python fractions.
	if r.Status != "ok" || r.TotalValue != "77.662993" {
		t.Fatalf("unexpected report: %+v", r)
	}
	if got := weightsOf(p); !reflect.DeepEqual(got, []string{"50.21", "24.64", "25.15"}) {
		t.Fatalf("weights = %v", got)
	}
}

func TestWeightsUnavailableNullsEverything(t *testing.T) {
	stale := "2026-11-01"
	n := json.Number("50.00")
	cases := map[string]*Portfolio{
		"missing quote": {SchemaVersion: 1, Coverage: "complete", Positions: []Position{
			pos("A", "USD", "1", "1", "2026-11-26"), pos("B", "USD", "1", "", "")}},
		"stale quote": {SchemaVersion: 1, Coverage: "complete", Positions: []Position{
			pos("A", "USD", "1", "1", "2026-11-26"), pos("B", "USD", "1", "1", stale)}},
		"future quote": {SchemaVersion: 1, Coverage: "complete", Positions: []Position{
			pos("A", "USD", "1", "1", "2026-11-27")}},
		"mixed currency, no base": {SchemaVersion: 1, Coverage: "complete", Positions: []Position{
			pos("A", "USD", "1", "1", "2026-11-26"), pos("B", "MXN", "1", "1", "2026-11-26")}},
	}
	for name, p := range cases {
		p.Positions[0].WeightPct = &n // stale value must be wiped
		r, err := ComputeWeights(p, opts())
		if err != nil || r.Status != "unavailable" || len(r.Reasons) == 0 {
			t.Fatalf("%s: report=%+v err=%v", name, r, err)
		}
		for _, w := range weightsOf(p) {
			if w != "null" {
				t.Fatalf("%s: weights not nulled: %v", name, weightsOf(p))
			}
		}
	}
}

func TestWeightsFX(t *testing.T) {
	p := &Portfolio{SchemaVersion: 1, Coverage: "complete", Positions: []Position{
		pos("A", "USD", "1", "100", "2026-11-26"),
		pos("B", "MXN", "10", "100", "2026-11-26"),
	}}
	o := opts()
	o.BaseCurrency, o.FX = "USD", map[string]string{"MXN": "0.05"}
	r, _ := ComputeWeights(p, o)
	if r.Status != "ok" || r.Currency != "USD" {
		t.Fatalf("%+v", r)
	}
	if got := weightsOf(p); !reflect.DeepEqual(got, []string{"66.67", "33.33"}) {
		t.Fatalf("weights = %v", got)
	}
}

func TestWeightsMalformed(t *testing.T) {
	bad := []*Portfolio{
		{SchemaVersion: 2, Coverage: "complete"},
		{SchemaVersion: 1, Coverage: "unknown"},
		{SchemaVersion: 1, Coverage: "complete", Positions: []Position{pos("A", "USD", "abc", "1", "2026-11-26")}},
		{SchemaVersion: 1, Coverage: "complete", Positions: []Position{pos("A", "USD", "1e3", "1", "2026-11-26")}},
		{SchemaVersion: 1, Coverage: "complete", Positions: []Position{pos("A", "USD", "0", "1", "2026-11-26")}},
		{SchemaVersion: 1, Coverage: "complete", Positions: []Position{pos("A", "USD", "1", "-5", "2026-11-26")}},
		{SchemaVersion: 1, Coverage: "complete", Positions: []Position{pos("A", "USD", "1", "1", "26-11-2026")}},
		{SchemaVersion: 1, Coverage: "complete", Positions: []Position{pos("A", "USD", "1", "1", "2026-11-26"), pos("a", "USD", "1", "1", "2026-11-26")}},
	}
	for i, p := range bad {
		if _, err := ComputeWeights(p, opts()); err == nil {
			t.Fatalf("case %d: expected error", i)
		}
	}
}

func TestWeightsEmptyComplete(t *testing.T) {
	r, err := ComputeWeights(&Portfolio{SchemaVersion: 1, Coverage: "complete"}, opts())
	if err != nil || r.Status != "ok" {
		t.Fatalf("%+v %v", r, err)
	}
}

func TestRunWeightsWriteIsIdempotentAndPreservesFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "positions.json")
	src := `{
  "schema_version": 1,
  "updated_at": "2026-11-26",
  "coverage": "partial",
  "positions": [
    {
      "ticker": "NVDA",
      "exchange": "NASDAQ",
      "currency": "USD",
      "shares": "2",
      "average_cost": "184.8",
      "price": "100",
      "price_as_of": "2026-11-26",
      "weight_pct": null,
      "as_of": "2026-11-25"
    },
    {
      "ticker": "META",
      "exchange": "NASDAQ",
      "currency": "USD",
      "shares": "1",
      "average_cost": null,
      "price": "100",
      "price_as_of": "2026-11-26",
      "weight_pct": null,
      "as_of": "2026-11-25"
    }
  ]
}
`
	os.WriteFile(path, []byte(src), 0o600)

	var out bytes.Buffer
	if err := runWeights([]string{"-file", path}, weightsNow, &out); err != nil {
		t.Fatal(err)
	}
	if after, _ := os.ReadFile(path); string(after) != src {
		t.Fatal("report-only run modified the file")
	}

	if err := runWeights([]string{"-file", path, "-write"}, weightsNow, &out); err != nil {
		t.Fatal(err)
	}
	first, _ := os.ReadFile(path)
	want := strings.Replace(strings.Replace(src, `"weight_pct": null`, `"weight_pct": 66.67`, 1), `"weight_pct": null`, `"weight_pct": 33.33`, 1)
	if string(first) != want {
		t.Fatalf("unexpected file:\n%s", first)
	}
	if err := runWeights([]string{"-file", path, "-write"}, weightsNow, &out); err != nil {
		t.Fatal(err)
	}
	if second, _ := os.ReadFile(path); !bytes.Equal(first, second) {
		t.Fatal("second run changed the file")
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %v", info.Mode().Perm())
	}
}

func TestRunWeightsRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "p.json")
	os.WriteFile(path, []byte(`{"schema_version":1,"updated_at":"x","coverage":"complete","positions":[],"broker":"x"}`), 0o600)
	if err := runWeights([]string{"-file", path, "-write"}, weightsNow, &bytes.Buffer{}); err == nil {
		t.Fatal("expected unknown field error")
	}
}
