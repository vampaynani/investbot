# Design: Add an evidence-gated tactical assessment

## Context

The evaluator owns one Markdown follow-up record per ticker. Its current snapshot and history cover the 3 to 5 year thesis, valuation, analyst consensus, and narratives. The `/investor-hoy` command asks the evaluator to review all tracked tickers, but it has no distinct short-term assessment.

The approved scope adds an evidence-gated tactical view without changing those long-term responsibilities. Existing ticker records remain untouched until their next evaluation.

## Goals / Non-Goals

**Goals:**

- Add a specialist that returns the same qualitative tactical fields for every ticker.
- Keep tactical analysis distinct from the long-term recommendation and history.
- Make repeated daily checks idempotent by UTC date.
- Check every canonical ticker record in `/investor-hoy` and exclude supplemental files.

**Non-Goals:**

- Re-evaluate or rewrite saved ticker records during implementation.
- Change the fundamental recommendation, consensus, valuation, or narrative rules.
- Perform live market research, calculate a numeric FOMO score, or recommend personalized position sizing.
- Add a programmatic market-data integration or a new external dependency.

## Decisions

### Keep the tactical result separate

`investor-tactical` owns the short-term evidence review. `investor-evaluator` remains responsible for the 3 to 5 year recommendation, consensus and valuation synthesis, ticker-record ownership, and overall history. The evaluator calls the tactical specialist for every new evaluation and for every ticker in a daily review.

Store the tactical result in separate snapshot and history sections in the existing canonical ticker file. A separate file would split one ticker's follow-up record and add another discovery path. Adding tactical fields to the long-term snapshot would mix distinct horizons and increase risk when updating existing conclusions.

### Use one qualitative result shape

Each result includes:

- `fomo_risk`: `elevated`, `mixed`, `low`, or `insufficient evidence`.
- `action`: `enter`, `wait`, `avoid`, or `N/D`.
- `catalyst`: verified event or `N/D`, with source and source timestamp.
- `horizon`: normally 1 to 3 months. Use a shorter horizon only for a verified event, and state its date and source.
- `entry_condition`, `exit_condition`, and `invalidation`: observable conditions or `N/D`.
- `evidence`: attributable links and timestamps for quote, price and volume, catalyst, and market context.
- `confidence`: qualitative confidence with a brief evidence-based reason.

Do not produce a numeric FOMO score. `elevated`, `mixed`, and `low` are qualitative judgments grounded in cited evidence, not values mapped from a hidden numeric threshold. Use `insufficient evidence` when an essential evidence class is missing, stale, or unattributable.

### Gate actions on evidence

Use a quote and price-volume evidence that meet the source's stated freshness. Attribute market context and any catalyst to a source with its timestamp. If freshness or attribution cannot be established, mark affected fields `N/D`. Do not repeat historical prices, events, or conditions as current.

- Use `enter` only when current evidence supports the setup and the entry condition is met. State an observable exit condition and invalidation.
- Use `wait` when the setup remains possible but a stated condition or verified catalyst has not occurred.
- Use `avoid` when current evidence shows that the setup is invalidated or does not support the entry.
- Use `N/D` when the evidence does not support an action.

FOMO risk and tactical action remain separate judgments. Elevated FOMO does not by itself require `avoid`, and low FOMO does not by itself support `enter`. Do not provide personalized position sizing.

### Select only canonical ticker files

Scan `.investor_evaluator/` at invocation time. Process filenames matching `[A-Z0-9.-]+\.md`. Exclude supplemental valuation documents and any other file that does not match. Run the tactical review for every canonical ticker, even when its long-term review is not due. If no canonical ticker files exist, say so and do not claim a ticker was checked.

### Keep same-day history updates idempotent

Record the initial tactical result and material tactical changes in tactical history. A daily rerun on the same UTC date updates that date's entry instead of appending a duplicate. A check with no material change updates the tactical snapshot and daily-check timestamp without creating a new history event. Keep history chronological.

## Risks / Trade-offs

- **Stale or unattributed data can imply a false signal.** Mark affected fields `N/D` and do not use unsupported data to choose an action.
- **The short-term action can be mistaken for the long-term recommendation.** Use separate sections and keep the evaluator's fundamental synthesis unchanged.
- **Repeated checks can duplicate history.** Key tactical history updates by UTC date and update the day's entry on reruns.
- **Supplemental Markdown files can look like ticker records.** Apply the canonical filename rule and explicitly exclude valuation supplements.

## Migration Plan

Add the tactical snapshot and history sections when each ticker is next evaluated. Do not migrate saved ticker records in bulk. Restart OpenCode after adding or changing the specialist so it loads the updated agent definition.

## Open Questions

None for the approved scope. Source freshness remains governed by the source's stated freshness and attribution, not a new arbitrary age threshold.
