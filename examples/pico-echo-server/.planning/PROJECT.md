# Pico Echo Server WebUI

## What This Is

一个用于测试 Pico Protocol 的 WebSocket 服务器，添加嵌入式 WebUI 用于可视化 Markdown 格式的消息收发。目标用户是开发者，用于快速测试和调试。

## Core Value

**轻量、零依赖、一键启动的 Markdown WebSocket 测试工具。**

如果用户能打开浏览器、连接 WebSocket、看到 Markdown 消息渲染，核心价值就达成了。

## Requirements

### Validated

- ✓ WebSocket server at `/ws` — existing (main.go)
- ✓ Bearer token authentication — existing (main.go)
- ✓ Ping/Pong protocol support — existing (main.go)
- ✓ message.send handling with stdout logging — existing (main.go)
- ✓ stdin broadcast as message.create — existing (main.go)
- ✓ Multiple client connection support — existing (main.go)

### Active

- [ ] Embedded WebUI served from `/` (Go embed)
- [ ] WebSocket connection from browser
- [ ] Message input form (text area + send button)
- [ ] Markdown rendering for received messages (marked.js)
- [ ] Auto-reconnect on connection loss

### Out of Scope

- Connection status display — minimal testing tool, not dashboard
- Message history persistence — ephemeral testing, no storage needed
- Session differentiation — single-view testing, multi-client not visualized
- Multiple server instances — single server scenario

## Context

- Part of PicoClaw project ecosystem
- Pico Protocol is a simple WebSocket messaging protocol
- Existing pico-echo-server is a minimal testing tool (160 lines Go)
- PicoClaw web/ directory has full launcher architecture (reference for patterns)
- Developers need visual feedback for Markdown content testing

## Constraints

- **Tech Stack**: Go + embed + pure HTML/JS — no build toolchain, minimal footprint
- **Size**: Stay lightweight (<500KB total) — testing tool should be fast
- **Browser**: Modern browser support only — no legacy IE compatibility

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Pure HTML + vanilla JS | No build step, embed directly, minimal complexity | — Pending |
| marked.js for Markdown | Lightweight, CDN available, simple integration | — Pending |
| Embed into binary | Single-file distribution, no external assets | — Pending |
| Replace stdin input | WebUI is the primary interface now | — Pending |

---
*Last updated: 2026-04-16 after initialization*

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state