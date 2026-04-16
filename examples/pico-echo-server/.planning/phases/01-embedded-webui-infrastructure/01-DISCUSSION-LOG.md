# Phase 1: Embedded WebUI Infrastructure - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-16
**Phase:** 01-embedded-webui-infrastructure
**Areas discussed:** Embed pattern, stdin handling

---

## Embed Pattern

| Option | Description | Selected |
|--------|-------------|----------|
| 单文件 (string) | `//go:embed index.html` + `var indexHTML string` — 最简单，无路径问题 | |
| 目录 embed.FS | `//go:embed all:web` + `fs.Sub()` — 文件分离，便于扩展 | ✓ |

**User's choice:** 目录 embed.FS (Recommended)
**Notes:** Research recommended directory embed for Phase 2/3 extensibility. User accepted recommendation.

---

## stdin Handling

| Option | Description | Selected |
|--------|-------------|----------|
| 删除 stdin | WebUI 作为唯一输入界面，简化架构 | ✓ |
| 保留 stdin | 支持两种输入方式（stdin + WebUI） | |

**User's choice:** 删除 stdin (Recommended)
**Notes:** WebUI becomes primary interface. stdin goroutine (~10 lines) will be removed.

---

## the agent's Discretion

- HTML content details (Phase 2/3 will add WebSocket/Markdown logic)
- CSS styling specifics (minimal, may adjust in Phase 3)
- Code placement (embed variable location in main.go)

---

## Deferred Ideas

None — discussion stayed within phase scope.