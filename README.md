# OpenCode configuration

The agent, command, and skill definitions in this directory configure stock-analysis workflows and provide reusable guidance for development tasks in OpenCode.

## Agents

Files in `agent/` define agents. Their frontmatter sets the agent description and whether the agent is a primary agent or a subagent.

| Agent | Role |
| --- | --- |
| `investor-evaluator` | Primary stock evaluator. Delegates work to the five specialist agents, combines their findings into a profile-based recommendation, and keeps one follow-up record per ticker in `.investor_evaluator/` at the project root. |
| `investor-business-quality` | Assesses the business model, competitive advantage, management, capital allocation, and durability. It does not make valuation conclusions. |
| `investor-valuation` | Compares price with earnings, cash flow, and growth over a 3 to 5 year holding period. It provides conservative, base, and optimistic scenarios. |
| `investor-target-consensus` | Verifies published analyst ratings and price targets as market-expectation data, separate from fundamental valuation and the evaluator's recommendation. |
| `investor-macro-context` | Assesses industry structure, demand, regulation, geopolitics, and supply-chain exposure. |
| `investor-risk-killer` | Stress-tests the thesis, identifies the condition that would invalidate it, and assesses downside and concentration risk. |

The evaluator requires specialists to separate facts from interpretation, cite sources, and avoid made-up figures. Its 3-to-5-year recommendation is based on the investor profile and verified evidence. Before finalizing, it reconciles the action with base-case return, downside, valuation, and unresolved thesis tests. Analyst ratings and 12-month consensus targets remain separate market-expectations lenses: they neither set fundamental valuation nor mechanically determine the action, but any apparent conflict with the 3-to-5-year recommendation must be explained by horizon and evidence. Full evaluations end with the bounded conclusion, review date, and sources; editorial adaptations are handled separately by the `investor-social-newsletter` skill.

For the standard user-facing structure, see [`docs/investor-report-template.md`](docs/investor-report-template.md). It provides a reusable report outline and a final consistency checklist for both people and LLMs; the evaluator agent instructions remain authoritative for research, recommendations, and follow-up records.

Ticker records live outside this directory, in `.investor_evaluator/` at the project root. The evaluator reads and updates those records when invoked. It does not run background checks between sessions.

## Command

`commands/investor-hoy.md` defines the `/investor-hoy` command and routes it to `investor-evaluator`. The command asks the evaluator to review every tracked ticker, check whether data is current, report changes, and keep its recommendation distinct from analyst consensus and fundamental valuation. If no ticker records exist, it asks the user to evaluate a ticker first. The command does not imply that checks run between invocations.

## Skills

Each directory in `skills/` contains a `SKILL.md` with the skill's purpose, selection guidance, and workflow. Some skills also include reference prompts, playbooks, scripts, or tests. OpenCode can apply a skill when its description matches the task; several also name explicit user phrases or slash commands.

### Task workflows

| Skill | What it does |
| --- | --- |
| `architect` | Sketches types, signatures, and module structure before a non-trivial implementation. Reference files provide design red flags, a rationale template, and a runner prompt. |
| `arena` | Compares parallel candidate solutions and combines the strongest parts. |
| `figure-it-out` | Plans and tracks a large or multi-phase task when no narrower workflow fits. It uses hypothesis checks and can keep a decision log with `show-me-your-work`. |
| `how` | Explains how code works, how responsibilities are divided, and where changes belong. Its reference prompts guide exploration and explanation. |
| `investor-social-newsletter` | Adapts an existing `.investor_evaluator/` evaluation for requested media (such as a social post, newsletter, email, or video script) without rerunning the analysis; preserves the saved investment view and source traceability. |
| `interrogate` | Runs an adversarial review from multiple angles. References define reviewer prompts, review criteria, and lead judgment. |
| `no-comments` | Reviews comments, fixes accepted findings, and suggests code structures for constraints that comments claim. |
| `poteto-mode` | Sets a concise, detailed work style with playbooks for bugs, features, investigations, refactors, shipping, worktrees, and other workflows. |
| `reflect` | Reviews the active session through tooling, judgment, and divergent-review lenses, then turns useful findings into edits to existing skills. |
| `setup-pstack` | Selects models and reasoning budgets for pstack roles, then writes a rule that overrides the skill defaults. |
| `show-me-your-work` | Keeps a TSV decision log for long or unattended work. It includes a log script and a decision-log template. |
| `swarm` | Runs parallel workers for broad coverage, races, or exploration and combines their reports. |
| `tdd` | Uses a failing test first when the user requests TDD or when a bug has a clear, inexpensive local test target. |
| `technical-writing` | Applies a layered style guide for documentation and other technical writing. |
| `unslop` | Removes common AI writing patterns. This skill applies to writing work throughout the repository. |
| `why` | Investigates design rationale and past decisions. References provide evidence guidance and source-specific instructions for code history, Linear, Sentry, Databricks, Notion, Datadog, Slack, and incident reports. |

### Engineering principles

These skills give focused guidance for recurring design and implementation decisions:

| Skill | Focus |
| --- | --- |
| `principle-attack-the-premise` | Re-examine the underlying premise after repeated fixes fail the same check. |
| `principle-boundary-discipline` | Validate at system boundaries and keep internal business logic direct. |
| `principle-build-the-lever` | Build a reusable tool that performs or verifies non-trivial work. |
| `principle-encode-lessons-in-structure` | Turn repeated instructions into checks, metadata, or code. |
| `principle-exhaust-the-design-space` | Compare multiple prototypes for a novel design problem. |
| `principle-experience-first` | Favor a polished user experience when product scope requires a trade-off. |
| `principle-fix-root-causes` | Reproduce bugs, trace symptoms to their cause, and fix the cause. |
| `principle-foundational-thinking` | Choose sound types and data structures before writing logic. |
| `principle-guard-the-context-window` | Keep large payloads and repetitive detail out of the main working context. |
| `principle-laziness-protocol` | Prefer the smallest change and remove needless abstraction. |
| `principle-make-operations-idempotent` | Design commands and processing loops to converge after retries or restarts. |
| `principle-migrate-callers-then-delete-legacy-apis` | Move callers to a replacement API and remove the old API in the same change. |
| `principle-minimize-reader-load` | Reduce indirection, hidden state, and the effort needed to understand code. |
| `principle-model-the-domain` | Represent domain rules in data structures instead of repeated conditionals. |
| `principle-never-block-on-the-human` | Proceed with reversible work and reserve questions for decisions that need approval. |
| `principle-outcome-oriented-execution` | Keep multi-phase rewrites focused on the final design, not temporary compatibility. |
| `principle-prove-it-works` | Verify the real behavior and inspect the resulting artifact before reporting completion. |
| `principle-redesign-from-first-principles` | Rework an existing design around a new foundational requirement. |
| `principle-separate-before-serializing-shared-state` | Remove shared writes where possible before adding serialization. |
| `principle-sequence-verifiable-units` | Split multi-step work into ordered changes that can each be checked. |
| `principle-subtract-before-you-add` | Remove dead or redundant pieces before adding a replacement. |
| `principle-test-behavior-not-implementation` | Assert observable results through the same interface users call. |
| `principle-type-system-discipline` | Use types to represent valid states, parse external data at boundaries, and cover variants. |

### Skill support files

- `skills/architect/references/` contains a design checklist, rationale template, and runner prompt.
- `skills/how/references/` contains prompts for exploring code and writing an explanation.
- `skills/interrogate/references/` contains review prompts, a rubric, and lead-judgment guidance.
- `skills/poteto-mode/playbooks/` contains workflow guides. Its `scripts/` directory contains plan checking, bootstrap and worktree tools, orchestration code, and the `watch-pr` command with tests and supporting modules. `scripts/package.json` and `scripts/bun.lock` describe that script package.
- `skills/poteto-mode/references/` contains BugBot triage guidance.
- `skills/reflect/references/` contains three reviewer prompts and a synthesis prompt.
- `skills/show-me-your-work/` contains the decision-log template and the script that appends log entries.
- `skills/why/references/` contains the investigation and synthesis prompts, evidence guidance, and source-specific playbooks.

## Local package files

`package.json` declares `@opencode-ai/plugin` at version `1.18.35`. `package-lock.json` and `node_modules/` are local npm install artifacts. The `.gitignore` in this directory excludes `node_modules`, package and lock files, `bun.lock`, and `.gitignore` itself.

The agent, command, and skill files are the configuration and instructions. The package dependency does not define a plugin entry point in this directory.

The `.DS_Store` files are Finder metadata and do not configure OpenCode.
