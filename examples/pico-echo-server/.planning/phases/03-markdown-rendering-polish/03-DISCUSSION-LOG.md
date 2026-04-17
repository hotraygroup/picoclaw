# Phase 3: Markdown Rendering & Polish - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-17
**Phase:** 03-markdown-rendering-polish
**Areas discussed:** XSS handling, auto-reconnect

---

## XSS Handling

| Option | Description | Selected |
|--------|-------------|----------|
| Trust server only | 服务器消息可信，innerHTML 渲染 | ✓ |
| DOMPurify sanitization | 严格白名单过滤，防止 XSS | |

**User's choice:** Trust server only (Recommended from research)
**Notes:** Server content is trusted. User input echoed back is displayed with textContent. DOMPurify adds complexity for testing tool.

---

## Auto-Reconnect

| Option | Description | Selected |
|--------|-------------|----------|
| Skip auto-reconnect | 测试工具重启频率低 | ✓ |
| Add auto-reconnect | 指数退避，最大 5 秒延迟 | |

**User's choice:** Skip auto-reconnect (Recommended)
**Notes:** Testing tool scenario — server restarts are infrequent, user can manually refresh.

---

## the agent's Discretion

- Markdown CSS details
- Code block styling
- Link styling

---

## Deferred Ideas

- DOMPurify sanitization — trusting server messages only
- Auto-reconnect — testing tool doesn't need
- Connection status indicator — v2 requirement PLSH-01