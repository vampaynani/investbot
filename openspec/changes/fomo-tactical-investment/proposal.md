# Proposal: Add an evidence-gated tactical assessment

## Why

Investbot's saved recommendation focuses on a 3 to 5 year investment thesis. It does not give a consistent, separately labeled view of short-term entry and exit conditions or FOMO risk. The tactical view must not change the long-term recommendation or imply current evidence when market data cannot be verified.

## What Changes

- Add an `investor-tactical` specialist for a default 1 to 3 month horizon.
- Require a qualitative FOMO category, tactical action, verified catalyst and source timestamp, horizon, entry and exit conditions, invalidation, evidence links, and confidence.
- Add a separate tactical snapshot and history section to each ticker record on its next evaluation. Do not mass-edit existing records.
- Require `/investor-hoy` to check every canonical ticker file on each invocation and show its tactical result separately from the long-term recommendation, valuation, and consensus.
- Update the report template and project documentation.

## Capabilities

### New Capabilities

- `investor-tactical-assessment`: evidence-gated FOMO risk and short-term entry and exit assessment, separate from the long-term recommendation.

### Modified Capabilities

None. The checkout has no existing project-local OpenSpec capabilities.

## Impact

Changes `.opencode/agent/investor-evaluator.md`, adds `.opencode/agent/investor-tactical.md`, updates `.opencode/commands/investor-hoy.md`, `docs/investor-report-template.md`, and `README.md`, and adds a separate tactical snapshot/history format for future evaluations. It does not change saved evaluations, assumption files, `.opencode/opencode.json`, backups, or perform live ticker research.
