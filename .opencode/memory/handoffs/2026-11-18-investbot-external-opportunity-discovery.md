---
project: investbot
task: external-opportunity-discovery
date: 2026-11-18
status: consumed
---

# External opportunity discovery handoff

## what changed

- Updated `.opencode/commands/investor-oportunidades.md` to discover U.S. and Mexican candidates from recent public lists and financial news, verify listings and evidence, compare available price performance with local benchmarks, and return a read-only prioritized list for full evaluation.
- Updated `README.md` to document the expanded discovery scope and `$ARGUMENTS` inputs.
- Created `.opencode/memory/handoffs/2026-11-18-investbot-external-opportunity-discovery.md`.
- `git diff --check` passed.
- `python3 -c 'from pathlib import Path; c=Path(".opencode/commands/investor-oportunidades.md").read_text(); r=Path("README.md").read_text(); required=["agent: investor-evaluator", "Estados Unidos y México", "listas públicas recientes y noticias financieras", "$ARGUMENTS", "moneda local", "hora UTC", "1, 3, 6 y 12 meses", "fuente y la fecha de ambos precios", "benchmark", "Continuidad de momentum", "Contrarian o reversión", "evidencia insuficiente", "N/D", "solo lectura", "No modifiques ni crees archivos", "No emitas una recomendación final `Buy`", "seguimiento en segundo plano"]; missing=[x for x in required if x not in c]; assert not missing, f"missing command requirements: {missing}"; assert "`$ARGUMENTS`" in r and "U.S. and Mexican markets" in r and "cites the sources, dates, and return basis" in r; assert all(line.rstrip()==line for text in (c,r) for line in text.splitlines()), "trailing whitespace"; print("PASS: command contract and README summary assertions; no trailing whitespace")'` passed.
- `git diff --quiet -- .opencode/commands/investor-hoy.md` passed.
- `opencode --help`, `opencode debug --help`, and `opencode debug config --help` showed no validation command. `opencode debug config` only displays resolved configuration, so it was not run.

## correction made 2026-11-19

- Replaced the final instruction that told `investor-evaluator` to delegate the full evaluation back to itself. It now tells the current evaluator to follow its existing full-evaluation procedure and delegate each selected candidate's analysis to the English-named specialists: `investor-business-quality`, `investor-valuation`, `investor-target-consensus`, `investor-macro-context`, and `investor-risk-killer`.
- Preserved the U.S.-and-Mexico discovery scope, read-only behavior, prioritized candidate shortlist, and no-final-`Buy` boundary.

Validation for this correction:

- `git diff --check`: PASS (exit 0).
- `python3 -c 'from pathlib import Path; c=Path(".opencode/commands/investor-oportunidades.md").read_text(); required=["agent: investor-evaluator", "Estados Unidos y México", "listas públicas recientes y noticias financieras", "$ARGUMENTS", "moneda local", "hora UTC", "1, 3, 6 y 12 meses", "fuente y la fecha de ambos precios", "benchmark", "Continuidad de momentum", "Contrarian o reversión", "evidencia insuficiente", "N/D", "solo lectura", "No modifiques ni crees archivos", "No emitas una recomendación final `Buy`", "seguimiento en segundo plano", "Como evaluador actual", "sigue tu procedimiento existente", "investor-business-quality", "investor-valuation", "investor-target-consensus", "investor-macro-context", "investor-risk-killer"]; missing=[x for x in required if x not in c]; assert not missing, f"missing command requirements: {missing}"; assert "lista priorizada para evaluación completa" in c; assert all(line.rstrip()==line for line in c.splitlines()), "trailing whitespace"; print("PASS: command routing, U.S./Mexico scope, discovery/evidence rules, read-only/no-auto-Buy boundary, shortlist, specialist delegation, and whitespace")'`: PASS.
- `git diff --quiet -- .opencode/commands/investor-hoy.md`: PASS (exit 0; unchanged).
- `/investor-oportunidades` was not invoked.

## correction made 2026-11-19 — shortlist handoff

- Clarified that `/investor-oportunidades` returns a prioritized shortlist and invites the user to choose candidates; the search itself does not run full evaluations. Only a later user request triggers the existing `investor-evaluator` procedure and delegation to `investor-business-quality`, `investor-valuation`, `investor-target-consensus`, `investor-macro-context`, and `investor-risk-killer`.
- Preserved the no-final-`Buy` rule and all English agent identifiers.

Validation for this correction:

- `git diff --check`: PASS (exit 0).
- `python3 -c 'from pathlib import Path; c=Path(".opencode/commands/investor-oportunidades.md").read_text(); required=["Esta búsqueda no ejecuta evaluaciones completas", "lista priorizada e invita al usuario a elegir", "Solo una solicitud posterior del usuario activa", "procedimiento existente de `investor-evaluator`", "investor-business-quality", "investor-valuation", "investor-target-consensus", "investor-macro-context", "investor-risk-killer", "No emitas una recomendación final `Buy`"]; missing=[x for x in required if x not in c]; assert not missing, f"missing requirements: {missing}"; assert all(line.rstrip()==line for line in c.splitlines()), "trailing whitespace"; print("PASS: shortlist-only boundary, later user-triggered evaluation routing, English agent identifiers, no-final-Buy rule, and whitespace")'`: PASS.
- `git diff --quiet -- .opencode/commands/investor-hoy.md`: PASS (exit 0; unchanged).
- `/investor-oportunidades` was not invoked.

## what another project needs to know

There is no cross-project contract.

## decisions or discoveries that might deserve durable memory

- External candidate discovery is limited to U.S. and Mexican markets, even though existing records cover other markets.
- Public lists and financial news provide leads only. A candidate needs a verified listing, quote, usable history, and current thesis evidence.
- `investor-evaluator` remains the only agent responsible for full evaluation and recommendation synthesis. The opportunities command remains read-only and does not claim background monitoring.

## resumed task: command names (2026-11-19)

- Renamed `/investor-hoy` to `/investor-today` in `.opencode/commands/investor-today.md`.
- Renamed `/investor-oportunidades` to `/investor-opportunities` in `.opencode/commands/investor-opportunities.md`; removed the old command files without adding Spanish aliases.
- Preserved the Spanish command instructions and English agent identifiers, including `investor-evaluator` and the specialist agents. Updated the root README command names and paths.
- `python3 -c '<rename, frontmatter, README, stale-name, contract, and whitespace assertions>'`: PASS. Both new command files exist, both route to `investor-evaluator`, README contains only the new command names, and the opportunity command retains U.S./Mexico scope, read-only shortlist behavior, and English specialist IDs.
- `git diff --check`: PASS (exit 0).
- `git status --short --untracked-files=all`: PASS. Only the approved paths changed; pre-existing unrelated files remain untouched.
- Neither command was invoked, and no live-market search was performed.
