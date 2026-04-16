# Phase 2: WebSocket Connection & Messaging - Context

**Gathered:** 2026-04-16
**Status:** Ready for planning

<domain>
## Phase Boundary

添加浏览器 WebSocket 客户端到 web/index.html，实现 bidirectional messaging。用户通过 WebUI 发送消息，服务器返回的消息在浏览器中显示。

**Scope:**
- WebSocket 客户端逻辑（内嵌 JS，无分离文件）
- message.send 发送格式（Pico Protocol JSON framing）
- message.create 接收和显示
- UI 交互（textarea 输入 + Send button）
- stdin goroutine 已在 Phase 1 删除（无需处理）

**Out of scope:**
- Markdown 渲染 (Phase 3)
- Auto-reconnect (Phase 3)
- Connection status UI (v2 requirement — explicitly deferred)
- Auth token handling (testing tool, open access)

</domain>

<decisions>
## Implementation Decisions

### WebSocket Connection
- **D-01:** WebSocket URL 使用 relative URL: `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws`
  - Rationale: 避免 Pitfall 2 (WebSocket origin mismatch)
  - Pattern: `new WebSocket(wsUrl)` in `<script>` block

### Message Framing
- **D-02:** 发送格式: JSON `{type: "message.send", timestamp: Date.now(), payload: {content: input}}`
- **D-03:** 接收格式: JSON `{type: "message.create", ...}` → 提取 `payload.content`
  - Pico Protocol framing (matching PicoClaw)

### UI Interaction
- **D-04:** 发送方式: 点击 Send button 发送（无 Enter 键快捷键）
  - Rationale: 简单交互，测试工具无需快捷键
  - Implementation: `sendBtn.onclick = () => { ws.send(...); input.value = ''; }`
  
- **D-05:** 消息显示: 用户消息左对齐，服务器消息右对齐
  - CSS: `.user { text-align: left; }`, `.server { text-align: right; }`
  - Implementation: appendMessage() 根据 type 设置 className

### the agent's Discretion
- Auto-scroll behavior (scrollTop = scrollHeight)
- Timestamp 显示格式
- 错误处理细节 (ws.onerror, ws.onclose)
- CSS 样式细节（颜色、间距）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### WebSocket Patterns
- `.planning/research/STACK.md` § Minimal HTML Template for WebSocket + Markdown — WebSocket URL construction pattern
- `.planning/research/PITFALLS.md` § Pitfall 2: WebSocket Origin Mismatch — 避免 localhost/127.0.0.1 冲突

### Project Context
- `.planning/research/SUMMARY.md` — Phase 2 implementation notes
- `.planning/REQUIREMENTS.md` — CORE-04, WS-01, WS-02, MSG-01, MSG-02
- `web/index.html` — 当前 HTML（Phase 1 创建）
- `main.go` — server broadcast() 函数（理解 message.create 格式）

</canonical_refs>

<code_context>
## Existing Code Insights

### web/index.html (Phase 1)
- 已有 placeholder elements: `#messages`, `#input`, `#send`
- CSS: max-width 800px, centered layout
- 无 JS 逻辑（Phase 1 只创建结构）

### main.go
- `handleWS()` — WebSocket handler (line 43-111)
- `broadcast()` — 发送 message.create (line 113-129)
- Message format: `picoMessage{Type, ID, SessionID, Timestamp, Payload}`
- `message.send` → server logs to stdout
- `message.create` → broadcast to all clients

### Required Changes to web/index.html
1. Add `<script>` block at end of `<body>`
2. WebSocket connection logic
3. `ws.onmessage` handler for message.create
4. `sendBtn.onclick` handler for message.send
5. `appendMessage()` helper function
6. CSS: `.user`, `.server` classes for message alignment

</code_context>

<specifics>
## Specific Ideas

- User message: 灰色背景，左对齐
- Server message: 白色背景，右对齐
- Auto-scroll to bottom on new message
- 简洁：无连接状态显示（deferred to v2）

</specifics>

<deferred>
## Deferred Ideas

- Connection status indicator ("Connected"/"Disconnected") — v2 requirement PLSH-01
- Enter 键快捷发送 — 简单测试工具不需要
- Timestamp 显示 — 可在 Phase 3 添加（与 Markdown 格式一起）

</deferred>

---

*Phase: 02-websocket-connection-messaging*
*Context gathered: 2026-04-16*