package main

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

const (
	sourceSECFiling         = "sec_filing"
	sourceAnalystNormalized = "analyst_normalized"

	returnBasis = "total return = price return + simple cumulative dividend yield on the initial price (dividends not reinvested); taxes, fees, and FX excluded"
)

// Options carries the evaluation clock and freshness policy. Passing Now
// explicitly keeps staleness checks deterministic and testable.
type Options struct {
	Now         time.Time
	MaxPriceAge time.Duration
}

type Input struct {
	Ticker        string          `json:"ticker"`
	Exchange      string          `json:"exchange"`
	Currency      string          `json:"currency"`
	Price         PriceInput      `json:"price"`
	FinancialBase FinancialInput  `json:"financial_base"`
	Scenarios     []ScenarioInput `json:"scenarios"`
	HorizonsYears []int           `json:"horizons_years"`
	Sensitivity   Sensitivity     `json:"sensitivity"`
}

type PriceInput struct {
	Value     float64 `json:"value"`
	AsOf      string  `json:"as_of"`
	SourceURL string  `json:"source_url"`
}

// FinancialInput is either an SEC filing figure (source_type "sec_filing") or
// an analyst-normalized figure derived from filings (source_type
// "analyst_normalized"), which must explain every adjustment.
type FinancialInput struct {
	SourceType    string  `json:"source_type"`
	Basis         string  `json:"basis"`
	ValuePerShare float64 `json:"value_per_share"`
	Period        string  `json:"period"`
	FiledAt       string  `json:"filed_at,omitempty"`
	SourceURL     string  `json:"source_url"`
	Normalization string  `json:"normalization_note,omitempty"`
}

type ScenarioInput struct {
	Name                   string   `json:"name"`
	AnnualGrowthPct        float64  `json:"annual_growth_pct"`
	ExitMultiple           float64  `json:"exit_multiple"`
	AnnualDividendYieldPct float64  `json:"annual_dividend_yield_pct,omitempty"`
	ProbabilityPct         *float64 `json:"probability_pct,omitempty"`
	ProbabilityBasis       string   `json:"probability_basis,omitempty"`
	ProbabilityURL         string   `json:"probability_source_url,omitempty"`
}

type Sensitivity struct {
	AnnualGrowthPct        []float64 `json:"annual_growth_pct"`
	ExitMultiples          []float64 `json:"exit_multiples"`
	TargetCAGRPct          []float64 `json:"target_cagr_pct"`
	AnnualDividendYieldPct float64   `json:"annual_dividend_yield_pct,omitempty"`
}

type Result struct {
	Ticker              string              `json:"ticker"`
	Exchange            string              `json:"exchange"`
	Currency            string              `json:"currency"`
	Price               PriceInput          `json:"price"`
	FinancialBase       FinancialInput      `json:"financial_base"`
	CurrentMultiple     *float64            `json:"current_multiple,omitempty"`
	ReturnBasis         string              `json:"return_basis"`
	ScenarioResults     []ScenarioResult    `json:"scenario_results"`
	Asymmetry           []AsymmetryResult   `json:"scenario_asymmetry"`
	ProbabilityWeighted []WeightedResult    `json:"probability_weighted_return,omitempty"`
	ProbabilitiesStatus string              `json:"probabilities_status"`
	Sensitivity         []SensitivityResult `json:"sensitivity"`
	ReverseValuation    []ReverseResult     `json:"reverse_valuation"`
	Limitations         []string            `json:"limitations"`
}

type ScenarioResult struct {
	Name           string        `json:"name"`
	ProbabilityPct *float64      `json:"probability_pct,omitempty"`
	Horizons       []HorizonCase `json:"horizons"`
}

type HorizonCase struct {
	Years                    int     `json:"years"`
	EndingMetricPerShare     float64 `json:"ending_metric_per_share"`
	EndingPrice              float64 `json:"ending_price"`
	CumulativePriceReturnPct float64 `json:"cumulative_price_return_pct"`
	AnnualizedPriceReturnPct float64 `json:"annualized_price_return_pct"`
	CumulativeDividendPct    float64 `json:"cumulative_dividend_income_pct"`
	CumulativeTotalReturnPct float64 `json:"cumulative_total_return_pct"`
	AnnualizedTotalReturnPct float64 `json:"annualized_total_return_pct"`
}

type AsymmetryResult struct {
	Years                        int      `json:"years"`
	ConservativeLossOverBaseGain *float64 `json:"conservative_loss_over_base_gain,omitempty"`
	OptimisticGainOverLoss       *float64 `json:"optimistic_gain_over_conservative_loss,omitempty"`
	Notes                        string   `json:"notes"`
}

type WeightedResult struct {
	Years                    int     `json:"years"`
	ExpectedCumulativeReturn float64 `json:"expected_cumulative_total_return_pct"`
	ExpectedAnnualizedReturn float64 `json:"probability_weighted_annualized_total_return_pct"`
	Caveat                   string  `json:"caveat"`
}

type SensitivityResult struct {
	GrowthPct    float64 `json:"annual_growth_pct"`
	ExitMultiple float64 `json:"exit_multiple"`
	Years        int     `json:"years"`
	EndingPrice  float64 `json:"ending_price"`
	ReturnPct    float64 `json:"cumulative_total_return_pct"`
	CAGRPct      float64 `json:"annualized_total_return_pct"`
}

// ReverseResult is price-only: it solves for the per-share growth needed to
// reach a target annualized price return, ignoring dividends.
type ReverseResult struct {
	ExitMultiple      float64 `json:"assumed_exit_multiple"`
	TargetCAGRPct     float64 `json:"target_price_cagr_pct"`
	Years             int     `json:"years"`
	RequiredGrowthPct float64 `json:"required_annual_metric_growth_pct"`
}

func Calculate(in Input, opts Options) (Result, error) {
	if err := validate(in, opts); err != nil {
		return Result{}, err
	}
	horizons := in.HorizonsYears
	if len(horizons) == 0 {
		horizons = []int{3, 5}
	}
	probabilityStatus, weightedOK, err := validateProbabilities(in.Scenarios)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		Ticker: strings.ToUpper(strings.TrimSpace(in.Ticker)), Exchange: in.Exchange,
		Currency: strings.ToUpper(in.Currency), Price: in.Price, FinancialBase: in.FinancialBase,
		ReturnBasis:         returnBasis,
		ProbabilitiesStatus: probabilityStatus,
		Limitations: []string{
			"The program calculates from supplied inputs; it does not download, verify, or refresh filings, prices, or assumptions.",
			"The financial base and the quote must refer to the same issuer, share class, currency, and split-adjusted basis.",
			"Scenario growth is per-share growth; dilution and buybacks must be reflected in that assumption.",
			"Dividends are a simple sum of the stated yield on the initial price, not reinvested; taxes, fees, and FX are excluded.",
			"An exit multiple is an assumption, not a forecast; sensitivity ranges are not probabilities.",
		},
	}

	if in.FinancialBase.ValuePerShare > 0 {
		m := round4(in.Price.Value / in.FinancialBase.ValuePerShare)
		result.CurrentMultiple = &m
	}

	for _, scenario := range in.Scenarios {
		sr := ScenarioResult{Name: scenario.Name, ProbabilityPct: scenario.ProbabilityPct}
		for _, years := range horizons {
			sr.Horizons = append(sr.Horizons, makeHorizon(in.Price.Value, in.FinancialBase.ValuePerShare,
				scenario.AnnualGrowthPct, scenario.ExitMultiple, scenario.AnnualDividendYieldPct, years))
		}
		result.ScenarioResults = append(result.ScenarioResults, sr)
	}

	result.Asymmetry = makeAsymmetry(result.ScenarioResults, horizons)
	if weightedOK {
		result.ProbabilityWeighted = makeWeighted(result.ScenarioResults, in.Scenarios, horizons)
	}
	result.Sensitivity = makeSensitivity(in, horizons)
	result.ReverseValuation = makeReverse(in, horizons)
	return result, nil
}

func validate(in Input, opts Options) error {
	if opts.Now.IsZero() {
		return errors.New("options.Now is required")
	}
	if opts.MaxPriceAge <= 0 {
		return errors.New("max price age must be greater than zero")
	}
	if blank(in.Ticker) || blank(in.Exchange) {
		return errors.New("ticker and exchange are required to disambiguate the listing")
	}
	if blank(in.Currency) {
		return errors.New("currency is required")
	}
	if err := validatePrice(in.Price, opts); err != nil {
		return err
	}
	if err := validateFinancialBase(in.FinancialBase, opts.Now); err != nil {
		return err
	}
	if len(in.Scenarios) != 3 {
		return errors.New("exactly three scenarios are required: conservative, base, optimistic")
	}
	want := []string{"conservative", "base", "optimistic"}
	for i, s := range in.Scenarios {
		if strings.ToLower(strings.TrimSpace(s.Name)) != want[i] {
			return fmt.Errorf("scenarios[%d].name must be %q", i, want[i])
		}
		if !finite(s.AnnualGrowthPct) || s.AnnualGrowthPct <= -100 {
			return fmt.Errorf("scenarios[%d].annual_growth_pct must be finite and greater than -100", i)
		}
		if !finitePositive(s.ExitMultiple) {
			return fmt.Errorf("scenarios[%d].exit_multiple must be finite and positive", i)
		}
		if !validYield(s.AnnualDividendYieldPct) {
			return fmt.Errorf("scenarios[%d].annual_dividend_yield_pct must be finite, >= 0 and < 100", i)
		}
		if s.ProbabilityPct != nil && (!finite(*s.ProbabilityPct) || *s.ProbabilityPct < 0 || *s.ProbabilityPct > 100) {
			return fmt.Errorf("scenarios[%d].probability_pct must be between 0 and 100", i)
		}
	}
	for _, years := range in.HorizonsYears {
		if years < 1 || years > 50 {
			return errors.New("horizons_years values must be between 1 and 50")
		}
	}
	for _, growth := range in.Sensitivity.AnnualGrowthPct {
		if !finite(growth) || growth <= -100 {
			return errors.New("sensitivity annual growth rates must be finite and greater than -100")
		}
	}
	for _, multiple := range in.Sensitivity.ExitMultiples {
		if !finitePositive(multiple) {
			return errors.New("sensitivity exit multiples must be finite and positive")
		}
	}
	for _, target := range in.Sensitivity.TargetCAGRPct {
		if !finite(target) || target <= -100 {
			return errors.New("sensitivity target CAGR values must be finite and greater than -100")
		}
	}
	if !validYield(in.Sensitivity.AnnualDividendYieldPct) {
		return errors.New("sensitivity annual_dividend_yield_pct must be finite, >= 0 and < 100")
	}
	return nil
}

func validatePrice(p PriceInput, opts Options) error {
	if !finitePositive(p.Value) {
		return errors.New("price.value must be finite and greater than zero")
	}
	asOf, err := time.Parse(time.RFC3339, p.AsOf)
	if err != nil {
		return fmt.Errorf("price.as_of must be RFC3339: %w", err)
	}
	if asOf.After(opts.Now.Add(time.Minute)) {
		return fmt.Errorf("price.as_of %s is in the future", p.AsOf)
	}
	if age := opts.Now.Sub(asOf); age > opts.MaxPriceAge {
		return fmt.Errorf("price.as_of %s is older than the allowed %s; obtain a current verified quote or mark the valuation N/D", p.AsOf, opts.MaxPriceAge)
	}
	return validateURL("price.source_url", p.SourceURL)
}

func validateFinancialBase(fb FinancialInput, now time.Time) error {
	if !finitePositive(fb.ValuePerShare) {
		return errors.New("financial_base.value_per_share must be finite and positive; use a positive, normalized per-share basis")
	}
	if blank(fb.Basis) || blank(fb.Period) {
		return errors.New("financial_base.basis and financial_base.period are required")
	}
	if err := validateURL("financial_base.source_url", fb.SourceURL); err != nil {
		return err
	}
	hasFiled := !blank(fb.FiledAt)
	if hasFiled {
		filed, err := time.Parse("2006-01-02", fb.FiledAt)
		if err != nil {
			return fmt.Errorf("financial_base.filed_at must be YYYY-MM-DD: %w", err)
		}
		if filed.After(now) {
			return fmt.Errorf("financial_base.filed_at %s is in the future", fb.FiledAt)
		}
	}
	switch fb.SourceType {
	case sourceSECFiling:
		if !hasFiled {
			return errors.New("sec_filing requires financial_base.filed_at")
		}
		if !isSECHost(fb.SourceURL) {
			return errors.New("sec_filing requires financial_base.source_url on sec.gov")
		}
	case sourceAnalystNormalized:
		if blank(fb.Normalization) {
			return errors.New("analyst_normalized requires normalization_note explaining every adjustment from the filing")
		}
	default:
		return fmt.Errorf("financial_base.source_type must be %q or %q", sourceSECFiling, sourceAnalystNormalized)
	}
	return nil
}

func validateProbabilities(scenarios []ScenarioInput) (string, bool, error) {
	count := 0
	total := 0.0
	for i, s := range scenarios {
		if s.ProbabilityPct == nil {
			continue
		}
		count++
		total += *s.ProbabilityPct
		if blank(s.ProbabilityBasis) {
			return "not calculated: probability basis missing", false, fmt.Errorf("scenarios[%d] has a probability but no probability_basis", i)
		}
		if err := validateURL(fmt.Sprintf("scenarios[%d].probability_source_url", i), s.ProbabilityURL); err != nil {
			return "not calculated: probability source missing", false, err
		}
	}
	if count == 0 {
		return "not calculated: no scenario probabilities supplied", false, nil
	}
	if count != len(scenarios) {
		return "not calculated: probabilities incomplete", false, errors.New("provide probabilities for all three scenarios or for none")
	}
	if math.Abs(total-100) > 0.01 {
		return "not calculated: probabilities do not sum to 100%", false, fmt.Errorf("scenario probabilities sum to %.4f%%, not 100%%", total)
	}
	return "calculated from supplied probabilities; evidence/source not independently assessed", true, nil
}

func makeHorizon(price, base, growthPct, multiple, yieldPct float64, years int) HorizonCase {
	n := float64(years)
	endingMetric := base * math.Pow(1+growthPct/100, n)
	endingPrice := endingMetric * multiple
	priceReturn := (endingPrice/price - 1) * 100
	priceCAGR := (math.Pow(endingPrice/price, 1/n) - 1) * 100
	dividends := yieldPct * n
	totalReturn := priceReturn + dividends
	totalCAGR := (math.Pow(1+totalReturn/100, 1/n) - 1) * 100
	return HorizonCase{
		Years:                    years,
		EndingMetricPerShare:     round4(endingMetric),
		EndingPrice:              round4(endingPrice),
		CumulativePriceReturnPct: round4(priceReturn),
		AnnualizedPriceReturnPct: round4(priceCAGR),
		CumulativeDividendPct:    round4(dividends),
		CumulativeTotalReturnPct: round4(totalReturn),
		AnnualizedTotalReturnPct: round4(totalCAGR),
	}
}

func makeAsymmetry(results []ScenarioResult, horizons []int) []AsymmetryResult {
	var out []AsymmetryResult
	for hi, years := range horizons {
		conservative := results[0].Horizons[hi].CumulativeTotalReturnPct
		base := results[1].Horizons[hi].CumulativeTotalReturnPct
		optimistic := results[2].Horizons[hi].CumulativeTotalReturnPct
		item := AsymmetryResult{Years: years, Notes: "total-return magnitudes only; not probabilities or maximum-loss estimates"}
		if conservative < 0 && base > 0 {
			ratio := round4(math.Abs(conservative) / base)
			item.ConservativeLossOverBaseGain = &ratio
		}
		if conservative < 0 && optimistic > 0 {
			ratio := round4(optimistic / math.Abs(conservative))
			item.OptimisticGainOverLoss = &ratio
		}
		out = append(out, item)
	}
	return out
}

func makeWeighted(results []ScenarioResult, scenarios []ScenarioInput, horizons []int) []WeightedResult {
	var out []WeightedResult
	for hi, years := range horizons {
		weightedCumulative := 0.0
		weightedAnnualized := 0.0
		for si, scenario := range scenarios {
			weight := *scenario.ProbabilityPct / 100
			weightedCumulative += weight * results[si].Horizons[hi].CumulativeTotalReturnPct
			weightedAnnualized += weight * results[si].Horizons[hi].AnnualizedTotalReturnPct
		}
		out = append(out, WeightedResult{
			Years:                    years,
			ExpectedCumulativeReturn: round4(weightedCumulative),
			ExpectedAnnualizedReturn: round4(weightedAnnualized),
			Caveat:                   "probability-weighted average of scenario outcomes; the annualized figure is the weighted average of scenario CAGRs, not the CAGR of the weighted cumulative return. The program does not validate calibration.",
		})
	}
	return out
}

func makeSensitivity(in Input, horizons []int) []SensitivityResult {
	var out []SensitivityResult
	for _, growth := range in.Sensitivity.AnnualGrowthPct {
		for _, multiple := range in.Sensitivity.ExitMultiples {
			for _, years := range horizons {
				c := makeHorizon(in.Price.Value, in.FinancialBase.ValuePerShare, growth, multiple, in.Sensitivity.AnnualDividendYieldPct, years)
				out = append(out, SensitivityResult{GrowthPct: growth, ExitMultiple: multiple, Years: years,
					EndingPrice: c.EndingPrice, ReturnPct: c.CumulativeTotalReturnPct, CAGRPct: c.AnnualizedTotalReturnPct})
			}
		}
	}
	return out
}

func makeReverse(in Input, horizons []int) []ReverseResult {
	var out []ReverseResult
	for _, multiple := range in.Sensitivity.ExitMultiples {
		for _, target := range in.Sensitivity.TargetCAGRPct {
			for _, years := range horizons {
				n := float64(years)
				targetEndingPrice := in.Price.Value * math.Pow(1+target/100, n)
				requiredGrowth := (math.Pow(targetEndingPrice/(in.FinancialBase.ValuePerShare*multiple), 1/n) - 1) * 100
				out = append(out, ReverseResult{ExitMultiple: multiple, TargetCAGRPct: target, Years: years, RequiredGrowthPct: round4(requiredGrowth)})
			}
		}
	}
	return out
}

func validateURL(field, raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return fmt.Errorf("%s must be an absolute http(s) URL", field)
	}
	return nil
}

func isSECHost(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	h := strings.ToLower(u.Hostname())
	return h == "sec.gov" || strings.HasSuffix(h, ".sec.gov")
}

func blank(s string) bool           { return strings.TrimSpace(s) == "" }
func round4(x float64) float64      { return math.Round(x*1e4) / 1e4 }
func validYield(x float64) bool     { return finite(x) && x >= 0 && x < 100 }
func finite(x float64) bool         { return !math.IsNaN(x) && !math.IsInf(x, 0) }
func finitePositive(x float64) bool { return finite(x) && x > 0 }
