# Investor Metrics (Go)

Deterministic scenario arithmetic for the evaluator. From supplied inputs it computes the current per-share multiple, three- and five-year scenario returns (price and total), annualized figures, scenario asymmetry, optional probability-weighted returns, a growth/multiple sensitivity grid, and reverse-valuation growth requirements.

It does **not** make a `Buy`/`Hold`/`Sell` recommendation, does not assign a valuation level, does not verify sources, and does not claim its assumptions are forecasts.

## Data sources

The program is offline: it reads a JSON file and makes no network requests. You fill in each value from a source you have checked.

- **Fundamentals, U.S. issuers only.** SEC EDGAR filings are the regulatory source for U.S. domestic filers. Use `"source_type": "sec_filing"`, which requires a `filed_at` date and a `source_url` on `sec.gov`. EDGAR does not cover most non-U.S. listings (e.g. `005930.KS`, `WALMEX`, `AMXB`, `CEMEXCPO`, `BIMBOA`, `GCARSOA1`, `KOFUBL`, `TLEVISACPO`, `VOLAR`, `ASURB`). For those, and for any analyst-normalized figure, use `"source_type": "analyst_normalized"`. This requires a `normalization_note` listing every adjustment from the reported figures and a `source_url` for the underlying documents. The filing date is optional.
- **Share price.** Provide a separately verified quote with `as_of` (RFC 3339), `source_url`, and currency. EDGAR is not a quote source. The program rejects a quote older than `-max-price-age` (default 24h) and any quote dated in the future, so a stale price cannot silently pass as current. A rejected quote means the valuation must be marked `N/D` until a current quote is obtained.
- Price and financial base must refer to the same issuer, share class, currency, and split-adjusted basis.

SEC API reference: <https://www.sec.gov/edgar/sec-api-documentation>.

## Run

```sh
cd tools/investor-metrics
go test ./...
go run . -input examples/sample.json
go run . -input examples/sample.json -max-price-age 2h
```

The sample numbers and quote are synthetic placeholders, **not investment data**. Its `as_of` will eventually be older than 24 hours, and the program will then refuse it. That is intended.

## Return definitions

- **Price return:** ending price / current price − 1, cumulative and annualized.
- **Dividend income:** `annual_dividend_yield_pct × years`, a simple sum on the initial price. Dividends are **not reinvested**, so this understates compounding when the yield is material. Taxes, fees, and FX are excluded.
- **Total return:** price return + dividend income. Asymmetry ratios and probability-weighted figures use total return.
- Sensitivity grid rows report total return using `sensitivity.annual_dividend_yield_pct`. Reverse valuation is **price-only**: it gives the per-share growth required for a target annualized price return at each assumed exit multiple.

## Input rules

- Exactly three scenarios, in order: `conservative`, `base`, `optimistic`.
- Growth is annual **per-share metric** growth. Reflect dilution or buybacks in this assumption.
- The exit multiple applies to the chosen per-share metric (EPS with P/E, FCF/share with P/FCF).
- The financial base must be positive. If earnings or FCF are negative or distorted, choose a positive normalized basis and say so, or do not use this simple multiple model.
- Optional probabilities must be given for all three scenarios, sum to 100%, and include a `probability_basis` and a `probability_source_url`. The program checks completeness and arithmetic only; it **cannot validate that probabilities are evidence-based or calibrated**. Omit them when no defensible basis exists.
- Output values are rounded to four decimals.

The tool returns calculations, not a valuation level or action. The analyst still judges business quality, permanent-loss risk, plausibility of the assumptions, and whether the return compensates for the downside.

## Known limits

- The probability-weighted annualized figure is the weighted average of scenario CAGRs, not the CAGR of the weighted cumulative return. The output labels this.
- Earnings-based scenarios assume the chosen metric is a usable proxy. Accounting distortions (investment gains, one-offs, stock-based compensation) must be handled in the normalization note, not by the program.
