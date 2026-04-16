# Phase 2: WebSocket Connection & Messaging - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-16
**Phase:** 02-websocket-connection-messaging
**Areas discussed:** WebSocket connection, Message framing, UI interaction

---

## WebSocket Connection

| Option | Description | Selected |
|--------|-------------|----------|
| Relative URL | `${location.host}/ws` — 自动匹配 origin | ✓ |
| Absolute URL | `ws://localhost:9090/ws` — 固定地址 | |

**User's choice:** Relative URL (Recommended from research)
**Notes:** Avoids Pitfall 2 (WebSocket origin mismatch on localhost variants)

---

## Message Framing

| Option | Description | Selected |
|--------|-------------|----------|
| Pico Protocol JSON | `{type: "message.send", ...}` | ✓ |
| Plain text | Raw string messages | |

**User's choice:** Pico Protocol JSON (from research)
**Notes:** Matches existing server format in main.go

---

## 发送方式 (UI Interaction)

| Option | Description | Selected |
|--------|-------------|----------|
| Button only | 点击按钮发送 | ✓ |
| Enter + Button | Enter 键快捷发送 + 按钮 | |
| Enter (Shift+Enter newline) | Enter 发送，换行快捷键 | |

**User's choice:** Button only (Recommended)
**Notes:** 简单交互，测试工具无需快捷键

---

## 消息显示样式 (UI Interaction)

| Option | Description | Selected |
|--------|-------------|----------|
| Side by side | 用户消息左对齐，服务器右对齐 | ✓ |
| Labels | [You] / [Server] 标签区分 | |
| Color distinction | 灰色/白色背景区分 | |

**User's choice:** Side by side (Recommended)
**Notes:** 类似聊天应用布局，视觉清晰

---

## the agent's Discretion

- Auto-scroll behavior
- Timestamp 显示格式
- 错误处理细节
- CSS 颜色/间距细节

---

## Deferred Ideas

- Connection status indicator — v2 requirement (PLSH-01)
- Enter 键快捷发送 — 测试工具简单即可
- Timestamp 显示 — Phase 3 可添加