# Tasks: Add an evidence-gated tactical assessment

## 1. Define and wire the tactical specialist

- [x] 1.1 Add `investor-tactical` with evidence freshness, attribution, output fields, qualitative FOMO rules, action definitions, horizon, and no-personalized-sizing safeguards.
- [x] 1.2 Update `investor-evaluator` to invoke the tactical specialist for new evaluations and every ticker in daily review.
- [x] 1.3 Define separate tactical snapshot and history sections, including same-UTC-day idempotency and preservation of existing records.

## 2. Update user-facing workflows

- [x] 2.1 Update `/investor-hoy` to scan every canonical filename and present the tactical lens separately.
- [x] 2.2 Update the report template and README to document the new specialist and tactical output.

## 3. Verify the change

- [x] 3.1 Review the changes for scope, field consistency, and preservation of existing long-term, consensus, valuation, and narrative rules.
- [x] 3.2 Run applicable local checks and report unavailable checks.
