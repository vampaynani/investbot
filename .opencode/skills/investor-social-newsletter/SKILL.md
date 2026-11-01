---
name: investor-social-newsletter
description: "Editorially adapt a saved `.investor_evaluator/<TICKER>.md` stock evaluation for a requested medium or audience (social post, newsletter, email, video script, etc.). Do not rerun the analysis or alter its investment view."
---

# Investor editorial adaptation

Turn a saved stock evaluation into editorial content for the medium(s) and audience the user requests. Read the saved record; do not repeat the five-specialist research process. The saved evaluation remains the source of truth; this skill changes presentation, not the investment judgment.

## Workflow

1. Resolve the exact ticker and listing from the request and the filename under `.investor_evaluator/`. Do not substitute another share class, ADR, or similarly named company.
2. Read the full current snapshot and the latest relevant history entries. Use the snapshot as the current conclusion; use history to understand material changes and avoid reviving an old thesis.
3. Check the analysis date, last verified price, last daily check, and next full review date. If the record is overdue, provisional, or missing current price/material evidence, say so in the editorial note and do not present stale data as current. Do not refresh data or run specialists unless the user also asks for a new evaluation.
4. Identify the requested medium, audience, language, length, and number of versions. Respect any specified editorial constraints. If the user names multiple media, create a clearly labeled version for each. If the medium is unspecified, default to one concise Spanish-language social post; ask a brief clarifying question only when the intended channel materially affects the deliverable.
5. Preserve the record’s conclusion exactly. Express its action in plain language without printing the scale label (`Strong Buy`, `Buy`, `Hold`, `Sell`, `Strong Sell`) or the 1–5 valuation score. For example, translate Hold into “la empresa sigue siendo interesante, pero al precio actual esperaría”; do not turn it into a Buy. If the record is Sell, communicate that the current risk/return does not justify investing or adding. If the record is N/D, say clearly that there is not enough current evidence for an investment view.
6. Explain valuation in ordinary language. If the record says the shares are expensive, explain **why at the current price**—for example, the earnings growth the price requires, a high normalized multiple, low cash-flow yield, or limited margin of safety. Do not mention the formal valuation score. If price or valuation is N/D, do not imply cheapness or expensiveness.
7. Preserve the core investment facts and caveats appropriate to each channel. For short public copy, normally include company/ticker, 3–5-year horizon, the same action in plain language, core thesis, concise valuation/return context, main risk, and a concrete condition that would change the view. Longer formats may organize these into a headline, summary, sections, or script; do not add unsupported details. Do not cite analyst targets as intrinsic value or use analyst labels to support the stance.
8. Keep sources traceable. Add a brief separate **“Respaldo editorial (no publicar)”** with the snapshot’s as-of date, internal recommendation/valuation level, and one or two source links already present in the record. Flag stale, provisional, overdue, or N/D information. Do not introduce new facts or claim fresh verification.

## Output

For each ticker, return one version per requested medium:

### [Medio / formato]

Use the structure, length, and editorial conventions appropriate to the channel. For social posts, default to one publishable Spanish paragraph of 60–100 words unless the user asks otherwise. Public-facing copy must not display the recommendation-scale label or the 1–5 valuation level; its plain-language stance must match them. A user may explicitly request a format that includes more analytical detail, but never misrepresent the saved view or expose an internal label as a public recommendation.

### Respaldo editorial (no publicar)

One short line with the data-as-of date, the underlying action and valuation level, plus source links from the saved record. Add a concise warning if the record is stale, provisional, overdue, or N/D.

## Guardrails

- Do not promise returns, imply certainty, hide a material risk, or upgrade/downgrade the saved view for a punchier post.
- Keep consensus, target price and fundamental value distinct. Normally omit consensus targets from the social copy; include one only if the user explicitly requests it and the record verifies its exact listing, currency, horizon and date.
- If the user asks only for an editorial adaptation, do not include the full analysis, re-evaluate the ticker, alter the saved file, or invent updated market data.
- If there is no matching saved record, ask for the ticker/listing or request a full evaluation first. Do not make a new stock recommendation from an empty record.
