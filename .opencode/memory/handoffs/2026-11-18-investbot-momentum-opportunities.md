---
project: investbot
task: momentum-opportunities
date: 2026-11-18
status: consumed
---

## what changed

- Added `.opencode/commands/investor-oportunidades.md`, routed to `investor-evaluator`. It scans only saved ticker records, separates momentum and contrarian or reversal candidates, requires current evidence, reports unsupported asymmetry as `N/D`, and remains read-only.
- Updated `README.md` to document the new command and distinguish its read-only behavior from full evaluations.
- Added this handoff at `.opencode/memory/handoffs/2026-11-18-investbot-momentum-opportunities.md`.

Validation commands and results:

- `opencode --help`: PASS. Confirmed the CLI command list.
- `opencode debug --help`: PASS. Confirmed available debug commands.
- `opencode debug config --help`: PASS. Confirmed the resolved-config diagnostic syntax. Did not run the diagnostic because it prints resolved configuration.
- `git diff --check -- README.md`: PASS.
- `python3 -c 'from pathlib import Path; s=Path(".opencode/commands/investor-oportunidades.md").read_text(); r=Path("README.md").read_text(); assert "\nagent: investor-evaluator\n---\n" in s; assert all(x in s for x in [".investor_evaluator/*.md", "Continuidad de momentum", "Contrarian o reversión", "nunca como señal de compra por sí solo", "Revalida los supuestos materiales", "indica `N/D`", "No deduzcas asimetría del precio objetivo de consenso ni del drawdown", "lista priorizada para seguimiento y evaluación completa", "No emitas una recomendación final `Buy`", "de solo lectura", "no afirmes que existe seguimiento en segundo plano"]); assert "commands/investor-oportunidades.md" in r and "/investor-oportunidades" in r; assert all(line == line.rstrip() for line in s.splitlines()); print("PASS: routing, scan boundaries, candidate groups, evidence rules, read-only limits, README entry, and whitespace")'`: PASS.
- `git diff --quiet -- .opencode/commands/investor-hoy.md`: PASS. The existing command is unchanged.
- `python3 -c 'from pathlib import Path; s=Path(".opencode/memory/handoffs/2026-11-18-investbot-momentum-opportunities.md").read_text(); assert all(x in s for x in ["project: investbot", "task: momentum-opportunities", "date: 2026-11-18", "status: open", "## what changed", "## what another project needs to know", "There is no cross-project contract", "## decisions or discoveries that might deserve durable memory"]); print("handoff frontmatter and required sections: PASS")'`: PASS.
- `python3 -c 'from pathlib import Path; paths=[Path(".opencode/commands/investor-oportunidades.md"),Path(".opencode/memory/handoffs/2026-11-18-investbot-momentum-opportunities.md")]; assert all(all(line==line.rstrip() for line in p.read_text().splitlines()) for p in paths); print("PASS: new Markdown whitespace")'`: PASS.
- `python3 -c 'import subprocess; st=subprocess.check_output(["git","status","--porcelain=v1","--untracked-files=all"],text=True).splitlines(); allowed={"README.md", ".opencode/commands/investor-oportunidades.md", ".opencode/memory/handoffs/2026-11-18-investbot-momentum-opportunities.md"}; own=[x for x in st if x[3:] in allowed]; outside=[x for x in st if x[3:] not in allowed]; print("authorized changed files: "+str(len(own))+"/3"); print("outside-scope changed files remain 12: "+("PASS" if len(outside)==12 else "FAIL ("+str(len(outside))+")")); raise SystemExit(len(own)!=3 or len(outside)!=12)'`: PASS. Exactly the three authorized files changed; the 12 pre-existing outside-scope changes remain.

The new market-search command was not invoked. The CLI help did not expose a command-file-specific validation command.

## what another project needs to know

There is no cross-project contract for this change.

## decisions or discoveries that might deserve durable memory

The command states its read-only boundary explicitly because `investor-evaluator` also handles full evaluations that update ticker records. No change to the agent definition or ticker data was needed.
