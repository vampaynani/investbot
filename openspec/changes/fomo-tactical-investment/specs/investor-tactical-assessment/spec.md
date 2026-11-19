# Investor tactical assessment

## ADDED Requirements

### Requirement: Keep the tactical assessment separate from the long-term recommendation

The evaluator MUST request a tactical assessment from `investor-tactical` for each new ticker evaluation. The evaluator MUST keep the tactical action and FOMO risk separate from the profile-based 3 to 5 year recommendation, fundamental valuation, analyst consensus, and long-term thesis.

#### Scenario: New ticker evaluation

- **WHEN** the evaluator assesses a ticker for the first time or performs a full reevaluation
- **THEN** it requests both the existing long-term analysis and a tactical assessment
- **AND** it presents the two actions as separate conclusions

### Requirement: Gate tactical conclusions on current attributable evidence

The tactical assessment MUST report FOMO risk as `elevated`, `mixed`, `low`, or `insufficient evidence`, and action as `enter`, `wait`, `avoid`, or `N/D`. It MUST report a verified catalyst and source timestamp, horizon, entry condition, exit condition, invalidation, evidence links and source timestamps for quote, price-volume context, catalyst, and market context, and confidence. It MUST use `N/D` for fields that current attributable evidence cannot support. It MUST NOT invent a signal, source, price, catalyst, or numeric FOMO score.

#### Scenario: Required current evidence is missing

- **WHEN** quote, price-volume, catalyst, or market-context evidence is missing, stale, or unattributable
- **THEN** the specialist marks affected fields `N/D`
- **AND** marks FOMO risk `insufficient evidence` when required evidence is incomplete
- **AND** uses action `N/D` when the missing evidence prevents a supported action

### Requirement: Use a bounded tactical horizon

The tactical horizon MUST default to 1 to 3 months. The specialist MUST shorten it only for a verified event and MUST cite the event source and timestamp when it uses a shorter horizon. The tactical assessment MUST NOT include personalized position sizing.

#### Scenario: Verified near-term event

- **WHEN** a current, attributable source verifies an event within the default horizon
- **THEN** the assessment may set a shorter horizon tied to that event
- **AND** records the event date, source, and source timestamp

### Requirement: Persist tactical state separately and idempotently

The evaluator MUST store the tactical snapshot and tactical history in their own clearly titled sections of the canonical ticker record. It MUST preserve existing long-term snapshot and history content. A same-UTC-day rerun MUST update that day's tactical history entry rather than add a duplicate. Implementation MUST NOT mass-edit existing ticker records.

#### Scenario: Repeated daily tactical check

- **WHEN** a ticker is checked more than once on the same UTC date
- **THEN** the evaluator updates that date's tactical entry and current snapshot
- **AND** does not append a second tactical event for the same day

### Requirement: Review every canonical ticker in the daily command

`/investor-hoy` MUST check every canonical ticker file on each invocation and present a compact tactical view separately from the long-term recommendation. It MUST select only filenames matching `[A-Z0-9.-]+\.md` and MUST exclude supplemental valuation documents and noncanonical files.

#### Scenario: Daily review with canonical and supplemental files

- **WHEN** the tracker directory contains canonical ticker files and supplemental Markdown files
- **THEN** the command reviews every canonical ticker file
- **AND** does not treat supplemental files as assets
- **AND** shows FOMO risk, tactical action, horizon, verified catalyst and timestamp, entry condition, exit condition, invalidation, evidence links and timestamps for each evidence class, and confidence separately

### Requirement: Preserve existing investment lenses

The tactical change MUST preserve the existing 3 to 5 year recommendation, analyst-consensus rating and target, fundamental valuation, and same-day history idempotency rules. Tactical evidence MUST NOT mechanically change those conclusions.

#### Scenario: Tactical and long-term conclusions differ

- **WHEN** a tactical action differs from the long-term recommendation
- **THEN** the report labels each horizon and conclusion separately
- **AND** retains the long-term recommendation based on its existing rules
