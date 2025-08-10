# GitHub Copilot Agent Mode – Repo Instructions (umatare5/controld-exporter)

### Scope & Metadata

- **Last Updated**: 2025-08-11

- **Precedence**: Highest in this repository (see §2)

- **Goals**:

  - **Primary Goal**: Build and maintain a **Prometheus exporter** for **Control D** (CLI binary + container image).
  - Keep public surface (CLI flags, exported metrics, HTTP endpoints) **clear, small, and stable**; maximize **readability/maintainability**.
  - Prefer **minimal diffs** (avoid unnecessary churn).
  - Defaults are secure (no secret logging; TLS and credentials handled safely).

- **Non‑Goals**:

  - Creating unrelated apps, SDKs, GUIs, or dashboards in this repo.
  - Introducing editor‑external, ad‑hoc lint rules or style changes.
  - Emitting or persisting secrets/test credentials.
  - Adding metrics that deviate from Prometheus exporter best practices (naming/types/labels) without explicit approval.

---

## 0. Normative Keywords

- **NK-001 (MUST)** Interpret **MUST / MUST NOT / SHOULD / SHOULD NOT / MAY** per RFC 2119/8174.

## 1. Repository Purpose & Scope

- **GP-001 (MUST)** Treat this repository as a **third‑party Prometheus Exporter for Control D**. Outputs: a CLI binary (`controld-exporter`) and a container image.
- **GP-002 (MUST)** Scope the code around: configuration/flags, HTTP server for `/metrics`, collectors, Control D API client, logging, and packaging (Docker/Goreleaser).
- **GP-003 (MUST NOT)** Assume library/SDK usage by other Go programs; the primary artifact is the exporter binary/image.
- **GP-004 (SHOULD)** Keep **metric names, help strings, types, and label sets** stable. Breaking changes must follow SemVer (see §13) and be called out in release notes.

## 2. Precedence & Applicability

- **GA-001 (MUST)** When editing/generating code, Copilot **must follow** this document.
- **GA-002 (MUST)** In this repository, **this file (`copilot-instructions.md`) has the highest precedence** over any other instruction set. **On conflict, prioritize this file**.
- **GA-003 (MUST)** Lint/format rules follow **repository settings only** (see §5).
- **GA-004 (MUST)** There is **no separate review instruction**. Review behavior is defined by this file as well.

## 3. Expert Personas (for AI edits/reviews)

- **EP-001 (MUST)** Act as a **Go 1.24 expert**.
- **EP-002 (MUST)** Act as an **exporter author** versed in **Prometheus best practices** (metric naming, types, label cardinality, help, `/metrics` HTTP handling, graceful shutdown).
- **EP-003 (SHOULD)** Be familiar with **Control D** domain concepts (personal vs business mode, endpoints/profiles/stats) and practical API client patterns.

## 4. Security & Privacy

- **SP-001 (MUST NOT)** Log tokens or credentials. **MUST** mask secrets (e.g., `${TOKEN:0:6}…`).
- **SP-002 (MUST)** Keep defaults secure (bind address/port, log level). Any insecure toggles are **opt‑in** and clearly labeled for dev/testing.
- **SP-003 (MUST)** Never write credentials to disk or VCS. Use env vars and process memory only.

## 5. Editor‑Driven Tooling (single source of truth)

- **ED-001 (MUST)** Lint/format/type checks follow repository settings (e.g., `.golangci.yml`, `.editorconfig`, `.markdownlint.json`, `.air.toml`, `.goreleaser.yml`).
- **ED-002 (MUST NOT)** Add flags/rules or inline disables that are not configured.
- **ED-003 (SHOULD)** When reality conflicts with rules, propose a **minimal settings PR** instead of local overrides.

## 6. Coding Principles (Basics)

- **GC-001 (MUST)** Apply **KISS/DRY** and keep code quality high.
- **GC-002 (MUST)** Avoid magic numbers; **use named constants** proactively (e.g., default port/path).
- **GC-003 (MUST)** Use **predicate helpers** (`is*`, `has*`) to keep call‑sites expressive.
- **GC-004 (SHOULD)** Isolate **Control D API** interactions behind a small client interface; keep collectors thin and testable.
- **GC-005 (SHOULD)** Prefer context‑aware functions and timeouts for outbound API calls.

## 7. Coding Principles (Conditionals)

- **CF-001 (MUST)** Prefer predicate helpers in conditions.
- **CF-002 (MUST)** Prefer **early returns** to reduce nesting and improve clarity.

## 8. Coding Principles (Loops)

- **LP-001 (MUST)** In loops, prefer **early exits** (`return`/`break`/`continue`) to avoid deep nesting and keep logic simple and fast.

## 9. Working Directory / Temp Files

- **WD-001 (MUST)** Place all temporary artifacts (work files, coverage, test binaries, etc.) **under `./tmp`**.
- **WD-002 (MUST)** Before completion, delete **zero‑byte files** (**exception**: keep `.keep`).

## 10. Model‑Aware Execution Workflow (when shell execution is available)

- **WF-001 (MUST)** Use `bash` explicitly (no shell auto‑detection).
- **WF-002 (MUST)** After editing Go code: run `go build ./...` (or `make build` if provided) and fix until it passes.
- **WF-003 (MUST)** After editing Docker‑related files: run `make image` (or `docker build`) and fix until success.
- **WF-004 (SHOULD)** After editing Prometheus examples/rules: run `promtool check config examples/prometheus.yml` and `promtool check rules examples/prometheus.alert_rules.yml`.
- **WF-005 (SHOULD)** After editing **shell scripts** under `./scripts/`: execute with documented options; ensure any related `make` targets succeed.
- **WF-006 (MUST)** On completion: summarize actions/results into `./.copilot_reports/<prompt_title>_<YYYY-MM-DD_HH-mm-ss>.md`.

> **Note**: This repo is **under construction**. Some targets/tests may not exist yet; **skip safely** but leave notes in the report for follow‑up.

## 11. Tests / Quality Gate (for AI reviewers)

- **QG-001 (MUST)** Keep CI green. Do not merge code that violates configured lint/format rules.
- **QG-002 (SHOULD)** When adding collectors or API calls, add **unit tests** around parsing/mapping, and **integration tests** guarded by env vars (e.g., `CTRLD_API_KEY`).
- **QG-003 (SHOULD)** Include **fixture‑based tests** for metric exposure (golden output or `testutil.CollectAndCompare`).

## 12. Change Scope & Tone (for AI reviewers)

- **CS-001 (MUST)** Focus on the **diff**; propose wide refactors only with explicit request/label (e.g., `allow-wide`).
- **CS-002 (SHOULD)** Tag comments with **\[BLOCKER] / \[MAJOR] / \[MINOR (Nit)] / \[QUESTION] / \[PRAISE]**.
- **CS-003 (SHOULD)** Structure comments as “**TL;DR → Evidence (rule/proof) → Minimal‑diff proposal**”.

## 13. Quick Checklist (before completion)

- **QC-001 (MUST, v1.0.0+)** **SemVer**: Metric **name/type/label** changes are **breaking** → require a major bump and release notes.
- **QC-002 (MUST)** Follow **`README.md`** for baseline usage (flags, ports, endpoints) and examples.
- **QC-003 (MUST)** Follow **packaging/release** rules (VERSION file, tag workflow, Goreleaser/GitHub Actions as configured).
- **QC-004 (MUST)** Lint/format are clean per editor settings (no ad‑hoc flags/inline disables).
- **QC-005 (MUST)** Temp artifacts under `./tmp`, zero‑byte files removed, and report written to `./.copilot_reports/`.
- **QC-006 (SHOULD)** Validate example Prometheus config and alert rules with `promtool`.
- **QC-007 (SHOULD)** Ensure container image listens on the documented port and exposes the configured metrics path.

---

### Repository‑Specific Notes (FYI for agents)

- **CLI surface**: `--web.listen-address`, `--web.listen-port` (default `10034`), `--web.telemetry-path` (default `/metrics`), `--controld.api-key` (`$CTRLD_API_KEY`), `--controld.business-mode`, `--log.level`, `--help`, `--version`.
- **Modes**: Default **personal** mode; **business** mode is opt‑in and may add org‑scoped labels/metrics.
- **Examples**: Prometheus scrape config & alert rules live under `examples/`; a Grafana dashboard JSON is provided.
- **Layout (indicative)**: `cmd/`, `internal/collector`, `internal/server`, `internal/config`, `internal/controld`, `internal/log`.
- **Packaging**: `Dockerfile`, `VERSION`, `.goreleaser.yml`/workflows. Use `make image` where available.

> Keep edits **conservative** and **operator‑friendly**: predictable flags, stable metrics, low cardinality labels, and clear release notes.
