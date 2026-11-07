package main

import (
	"math"
	"strings"
	"testing"
	"time"
)

var testNow = time.Date(2026, 11, 6, 21, 0, 0, 0, time.UTC)

func testOptions() Options {
	return Options{Now: testNow, MaxPriceAge: 24 * time.Hour}
}

func testInput() Input {
	return Input{
		Ticker: "TEST", Exchange: "NASDAQ", Currency: "USD",
		Price: PriceInput{Value: 200, AsOf: "2026-11-06T20:00:00Z", SourceURL: "https://example.com/quote"},
		FinancialBase: FinancialInput{
			SourceType: sourceAnalystNormalized, Basis: "normalized EPS", ValuePerShare: 10, Period: "FY2025",
			SourceURL: "https://example.com/filing-summary", Normalization: "test normalization",
		},
		Scenarios: []ScenarioInput{
			{Name: "conservative", AnnualGrowthPct: 0, ExitMultiple: 16},
			{Name: "base", AnnualGrowthPct: 8, ExitMultiple: 20},
			{Name: "optimistic", AnnualGrowthPct: 15, ExitMultiple: 24},
		},
		HorizonsYears: []int{3, 5},
		Sensitivity:   Sensitivity{AnnualGrowthPct: []float64{5, 10}, ExitMultiples: []float64{16, 20}, TargetCAGRPct: []float64{8}},
	}
}

func approx(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Fatalf("%s = %.6f, want %.6f (±%.6f)", name, got, want, tol)
	}
}

func TestBaseValuesAndCurrentMultiple(t *testing.T) {
	got, err := Calculate(testInput(), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	approx(t, "current multiple", *got.CurrentMultiple, 20, 1e-9)
	base5 := got.ScenarioResults[1].Horizons[1]
	// 10 * 1.08^5 = 14.69328 per share; x20 = 293.8656.
	approx(t, "base 5y ending price", base5.EndingPrice, 293.8656, 1e-4)
	approx(t, "base 5y price return", base5.CumulativePriceReturnPct, 46.9328, 1e-3)
	approx(t, "base 5y CAGR", base5.AnnualizedTotalReturnPct, 8.0, 1e-3)
	if base5.CumulativeDividendPct != 0 {
		t.Fatalf("dividends = %v, want 0 when no yield is given", base5.CumulativeDividendPct)
	}
}

func TestDividendsAddSimpleYieldTimesYears(t *testing.T) {
	in := testInput()
	in.Scenarios[1].AnnualDividendYieldPct = 2
	got, err := Calculate(in, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	h := got.ScenarioResults[1].Horizons[1] // 5 years
	approx(t, "dividend income", h.CumulativeDividendPct, 10, 1e-9)
	approx(t, "total return", h.CumulativeTotalReturnPct, h.CumulativePriceReturnPct+10, 1e-3)
	if h.AnnualizedTotalReturnPct <= h.AnnualizedPriceReturnPct {
		t.Fatal("annualized total return should exceed price CAGR when yield is positive")
	}
}

func TestAsymmetryUsesTotalReturn(t *testing.T) {
	got, err := Calculate(testInput(), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Asymmetry) != 2 {
		t.Fatalf("asymmetry rows = %d, want 2", len(got.Asymmetry))
	}
	// Conservative 5y: 160 / 200 - 1 = -20%. Base 5y total = +46.9328%.
	a5 := got.Asymmetry[1]
	if a5.ConservativeLossOverBaseGain == nil {
		t.Fatal("expected conservative loss / base gain ratio")
	}
	approx(t, "conservative loss over base gain", *a5.ConservativeLossOverBaseGain, 20/46.9328, 1e-3)
	if a5.OptimisticGainOverLoss == nil {
		t.Fatal("expected optimistic gain / loss ratio")
	}
}

func TestAsymmetryOmittedWhenConservativeNotLoss(t *testing.T) {
	in := testInput()
	in.Scenarios[0].ExitMultiple = 25 // conservative becomes a gain
	got, err := Calculate(in, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if got.Asymmetry[0].ConservativeLossOverBaseGain != nil || got.Asymmetry[0].OptimisticGainOverLoss != nil {
		t.Fatal("ratios must be N/D (omitted) when conservative case is not a loss")
	}
}

func TestRejectsStalePrice(t *testing.T) {
	in := testInput()
	in.Price.AsOf = "2026-11-04T20:00:00Z" // ~49h before testNow
	_, err := Calculate(in, testOptions())
	if err == nil || !strings.Contains(err.Error(), "older than") {
		t.Fatalf("error = %v, want stale-price rejection", err)
	}
}

func TestRejectsFuturePrice(t *testing.T) {
	in := testInput()
	in.Price.AsOf = "2026-11-07T20:00:00Z"
	_, err := Calculate(in, testOptions())
	if err == nil || !strings.Contains(err.Error(), "in the future") {
		t.Fatalf("error = %v, want future-price rejection", err)
	}
}

func TestMaxPriceAgeIsConfigurable(t *testing.T) {
	in := testInput()
	in.Price.AsOf = "2026-11-06T20:00:00Z" // 1h before testNow
	if _, err := Calculate(in, Options{Now: testNow, MaxPriceAge: 30 * time.Minute}); err == nil {
		t.Fatal("a 1h-old quote must be rejected with a 30m limit")
	}
	if _, err := Calculate(in, Options{Now: testNow, MaxPriceAge: 2 * time.Hour}); err != nil {
		t.Fatalf("a 1h-old quote must be accepted with a 2h limit: %v", err)
	}
}

func TestRejectsNonPositiveMaxAgeAndMissingClock(t *testing.T) {
	if _, err := Calculate(testInput(), Options{Now: testNow, MaxPriceAge: 0}); err == nil {
		t.Fatal("zero max age must be rejected")
	}
	if _, err := Calculate(testInput(), Options{MaxPriceAge: time.Hour}); err == nil {
		t.Fatal("missing Now must be rejected")
	}
}

func TestSECFilingRequiresSECHostAndFiledAt(t *testing.T) {
	in := testInput()
	in.FinancialBase.SourceType = sourceSECFiling
	in.FinancialBase.FiledAt = "2026-10-29"
	in.FinancialBase.SourceURL = "https://example.com/not-sec"
	if _, err := Calculate(in, testOptions()); err == nil || !strings.Contains(err.Error(), "sec.gov") {
		t.Fatalf("error = %v, want sec.gov host requirement", err)
	}

	in.FinancialBase.SourceURL = "https://www.sec.gov/Archives/edgar/data/1/filing"
	in.FinancialBase.FiledAt = ""
	if _, err := Calculate(in, testOptions()); err == nil || !strings.Contains(err.Error(), "filed_at") {
		t.Fatalf("error = %v, want filed_at requirement", err)
	}

	in.FinancialBase.FiledAt = "2026-10-29"
	if _, err := Calculate(in, testOptions()); err != nil {
		t.Fatalf("valid sec_filing rejected: %v", err)
	}
}

func TestAnalystNormalizedRequiresNormalizationNote(t *testing.T) {
	in := testInput()
	in.FinancialBase.Normalization = "   "
	if _, err := Calculate(in, testOptions()); err == nil || !strings.Contains(err.Error(), "normalization_note") {
		t.Fatalf("error = %v, want normalization_note requirement", err)
	}
}

func TestRejectsUnknownSourceTypeAndFutureFiling(t *testing.T) {
	in := testInput()
	in.FinancialBase.SourceType = "web_scrape"
	if _, err := Calculate(in, testOptions()); err == nil {
		t.Fatal("unknown source_type must be rejected")
	}
	in = testInput()
	in.FinancialBase.FiledAt = "2026-12-01"
	if _, err := Calculate(in, testOptions()); err == nil || !strings.Contains(err.Error(), "future") {
		t.Fatalf("error = %v, want future filed_at rejection", err)
	}
}

func TestCalculateWeightsOnlyFullySourcedProbabilities(t *testing.T) {
	in := testInput()
	probs := []float64{25, 50, 25}
	for i := range in.Scenarios {
		in.Scenarios[i].ProbabilityPct = &probs[i]
		in.Scenarios[i].ProbabilityBasis = "test calibration"
		in.Scenarios[i].ProbabilityURL = "https://example.com/probability-method"
	}
	got, err := Calculate(in, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ProbabilityWeighted) != 2 {
		t.Fatalf("weighted results = %d, want 2", len(got.ProbabilityWeighted))
	}
	if !strings.Contains(got.ProbabilitiesStatus, "not independently assessed") {
		t.Fatalf("unexpected probability status: %s", got.ProbabilitiesStatus)
	}
	w := got.ProbabilityWeighted[1]
	base := got.ScenarioResults[1].Horizons[1].CumulativeTotalReturnPct
	cons := got.ScenarioResults[0].Horizons[1].CumulativeTotalReturnPct
	opt := got.ScenarioResults[2].Horizons[1].CumulativeTotalReturnPct
	approx(t, "weighted cumulative", w.ExpectedCumulativeReturn, 0.25*cons+0.5*base+0.25*opt, 1e-3)
}

func TestRejectsUnsourcedOrIncompleteProbabilities(t *testing.T) {
	in := testInput()
	p := 100.0
	in.Scenarios[0].ProbabilityPct = &p
	in.Scenarios[0].ProbabilityBasis = "test basis"
	in.Scenarios[0].ProbabilityURL = "https://example.com/probability-method"
	_, err := Calculate(in, testOptions())
	if err == nil || !strings.Contains(err.Error(), "all three scenarios") {
		t.Fatalf("error = %v, want incomplete-probabilities validation", err)
	}
}

func TestRejectsProbabilitiesNotSummingToHundred(t *testing.T) {
	in := testInput()
	probs := []float64{30, 30, 30}
	for i := range in.Scenarios {
		in.Scenarios[i].ProbabilityPct = &probs[i]
		in.Scenarios[i].ProbabilityBasis = "b"
		in.Scenarios[i].ProbabilityURL = "https://example.com/p"
	}
	if _, err := Calculate(in, testOptions()); err == nil || !strings.Contains(err.Error(), "sum to") {
		t.Fatalf("error = %v, want sum-to-100 rejection", err)
	}
}

func TestReverseValuationRecoversKnownGrowth(t *testing.T) {
	in := testInput()
	in.Sensitivity.ExitMultiples = []float64{20}
	in.Sensitivity.TargetCAGRPct = []float64{8}
	in.HorizonsYears = []int{5}
	got, err := Calculate(in, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ReverseValuation) != 1 {
		t.Fatalf("reverse rows = %d, want 1", len(got.ReverseValuation))
	}
	// Price 200 grows 8%/yr for 5y, exit at 20x on EPS 10 -> EPS growth must be 8%.
	approx(t, "required growth", got.ReverseValuation[0].RequiredGrowthPct, 8.0, 1e-3)
}

func TestSensitivityGridSize(t *testing.T) {
	got, err := Calculate(testInput(), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Sensitivity) != 2*2*2 {
		t.Fatalf("sensitivity rows = %d, want 8", len(got.Sensitivity))
	}
}

func TestOutputIsRoundedToFourDecimals(t *testing.T) {
	got, err := Calculate(testInput(), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	for _, sr := range got.ScenarioResults {
		for _, h := range sr.Horizons {
			v := h.CumulativePriceReturnPct
			if math.Abs(v*1e4-math.Round(v*1e4)) > 1e-6 {
				t.Fatalf("value %.12f is not rounded to 4 decimals", v)
			}
		}
	}
}
