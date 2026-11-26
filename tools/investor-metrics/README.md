# Investor Metrics (Go)

Deterministic scenario arithmetic for the evaluator. From supplied inputs it computes the current per-share multiple, three- and five-year scenario returns (price and total), annualized figures, scenario asymmetry, optional probability-weighted returns, a growth/multiple sensitivity grid, and reverse-valuation growth requirements.

It does **not** make a `Buy`/`Hold`/`Sell` recommendation, does not assign a valuation level, does not verify sources, and does not claim its assumptions are forecasts.

When capex conversion is material to the thesis, one input may also include exactly two optional `synthetic_fcf_bases` (e.g. `partial_conversion` and `strong_conversion`). Each entry contains a name, an explicitly normalized/synthetic `financial_base`, and its own conservative/base/optimistic scenarios. The same run returns `synthetic_fcf_results` with current implied P/FCF, 3-/5-year scenarios, asymmetry, sensitivity, and reverse-valuation results for each alternative, alongside the primary scenario output. These are supplemental hypotheses and must never be confused with reported FCF or replace the primary evidence-based valuation.

## Portfolio weights

`investor-metrics weights` computes `weight_pct` for `.investor_portfolio/positions.json`, used by `/investor-portfolio` on every update.

```sh
go run . weights -file ../../.investor_portfolio/positions.json            # report only
go run . weights -file ../../.investor_portfolio/positions.json -write     # also write weight_pct
go run . weights -file ... -write -base USD -fx MXN=0.054                  # mixed currencies
```

- Weight = `shares × price` ÷ sum of `shares × price` over the listed positions × 100, in exact decimal arithmetic (no floats), rounded half up to 2 decimals. Same file and date give byte-identical output; a second `-write` is a no-op. Rounded weights may sum to 100.00 ± 0.01.
- Quotes (`price`, `price_as_of`) must already be in the file, verified by the caller. The program is offline. It rejects a quote older than `-max-price-age-days` (default 3, which covers a weekend) or dated in the future.
- If any quote is missing, stale, or future-dated, or currencies are mixed without `-base` and an `-fx` rate for each other currency, no weight can be trusted: every `weight_pct` is written as `null`, the output is `status: "unavailable"` with the reasons, and the exit code is 0. A malformed file (unknown field, bad decimal, duplicate listing, unsupported schema) exits 1 and writes nothing.
- Weights are over the listed positions only; with `coverage: partial` they are not weights of the whole portfolio.

## Folder convention

```
tools/investor-metrics/
├── fixtures/          # verified input assumptions (tracked in git)
│   └── <ticker>_assumptions.json
├── evals/             # dated evaluation outputs (ignored, regenerable)
│   └── <ticker>_eval_YYYYMMDD.json
├── examples/          # synthetic placeholders, not investment data
│   └── sample.json
├── main.go, metrics.go, weights.go, *_test.go
└── go.mod
```

- **`fixtures/`**: inputs de supuestos verificados para una empresa. Son parte del proyecto y se rastrean en git. Un agente o skill que necesite un archivo de entrada para una empresa existente lo lee desde aquí; uno nuevo se crea aquí con nombre `<ticker>_assumptions.json`.
- **`evals/`**: salidas con fecha de corridas de evaluación (escenarios, FCF sintéticos, bases sintéticas). Son resultado regenerable, no se rastrean en git (`tools/investor-metrics/evals/` está en `.gitignore`). Los agentes escriben sus resultados aquí con nombre `<ticker>_eval_YYYYMMDD.json` (o similar con fecha).
- **`examples/`**: placeholders sintéticos, nunca datos de inversión reales.

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
- Optional `synthetic_fcf_bases`: if present, it must contain exactly two uniquely named alternatives. Each requires a positive per-share financial base with a full normalization note and exactly three scenarios in the required order. Label hypothetical values as synthetic in both `basis` and `normalization_note`.
- Growth is annual **per-share metric** growth. Reflect dilution or buybacks in this assumption.
- The exit multiple applies to the chosen per-share metric (EPS with P/E, FCF/share with P/FCF).
- The financial base must be positive. If earnings or FCF are negative or distorted, choose a positive normalized basis and say so, or do not use this simple multiple model.
- Optional probabilities must be given for all three scenarios, sum to 100%, and include a `probability_basis` and a `probability_source_url`. The program checks completeness and arithmetic only; it **cannot validate that probabilities are evidence-based or calibrated**. Omit them when no defensible basis exists.
- Output values are rounded to four decimals.

The tool returns calculations, not a valuation level or action. The analyst still judges business quality, permanent-loss risk, plausibility of the assumptions, and whether the return compensates for the downside.

## Known limits

- The probability-weighted annualized figure is the weighted average of scenario CAGRs, not the CAGR of the weighted cumulative return. The output labels this.
- Earnings-based scenarios assume the chosen metric is a usable proxy. Accounting distortions (investment gains, one-offs, stock-based compensation) must be handled in the normalization note, not by the program.
