# Phase 1: Embedded WebUI Infrastructure - Context

**Gathered:** 2026-04-16
**Status:** Ready for planning

<domain>
## Phase Boundary

添加 Go embed 基础设施，使单文件二进制能够通过 HTTP 根路径 (`/`) 服务 WebUI。保持现有 WebSocket endpoint (`/ws`) 正常工作，删除 stdin goroutine（WebUI 成为主界面）。

**Scope:**
- 创建 `web/` 目录结构
- 添加 `//go:embed all:web` 指令
- 设置 `http.ServeMux` 路由
- 保持 `/ws` endpoint 不变
- 删除 stdin 输入代码

**Out of scope:**
- WebSocket 客户端逻辑 (Phase 2)
- Markdown 渲染 (Phase 3)
- UI 样式细节 (后续 phases)

</domain>

<decisions>
## Implementation Decisions

### Embed Pattern
- **D-01:** 使用目录 embed.FS (`//go:embed all:web` + `fs.Sub()`)
  - Rationale: 文件分离便于后续扩展，Phase 2/3 需要添加 JS/CSS
  - Pattern: `subFS, _ := fs.Sub(webFS, "web")` + `http.FileServer(http.FS(subFS))`

### stdin Handling
- **D-02:** 删除 stdin goroutine（约10行代码）
  - Rationale: WebUI 作为唯一输入界面，简化架构
  - stdin 输入在 Phase 1 后不再可用

### Route Structure
- **D-03:** 使用 `http.ServeMux` 替代全局 `http.HandleFunc`
  - Rationale: 明确路由映射，避免 route 冲突
  - Routes: `/` → WebUI, `/ws` → WebSocket handler

### the agent's Discretion
- HTML 文件内容（基础结构，Phase 2/3 会添加内容）
- CSS 样式（最小样式，Phase 3 可能调整）
- 具体代码位置（embed 变量放在 main.go 或分离文件）

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Embed Pattern
- `.planning/research/STACK.md` § Go Embed Patterns — Pattern 2: Directory Embed
- `.planning/research/PITFALLS.md` § Pitfall 1: embed.FS Path Traps — 防止 404 错误
- `https://pkg.go.dev/embed` — Go embed 官方文档

### Project Context
- `.planning/research/SUMMARY.md` — Phase 1 实现说明
- `.planning/PROJECT.md` — Core Value: 轻量、零依赖
- `.planning/REQUIREMENTS.md` — CORE-01, CORE-02, CORE-03

### Existing Code Reference
- `main.go` — 当前实现，需要修改

</canonical_refs>

<code_context>
## Existing Code Insights

### Current Implementation (main.go)
- `http.HandleFunc("/ws", s.handleWS)` — 全局路由，需要改为 ServeMux
- stdin goroutine (line 147-157) — 需要删除
- `upgrader = websocket.Upgrader{CheckOrigin: ...}` — 保持不变
- `server` struct with `token`, `mu`, `conns` — 保持不变

### Required Changes
1. Add imports: `embed`, `io/fs`
2. Add directive: `//go:embed all:web`
3. Add embed variable: `var webFS embed.FS`
4. Change: `http.HandleFunc` → `http.NewServeMux()` + `mux.Handle`
5. Delete: stdin goroutine (10 lines)
6. Add: WebUI route registration

### File to Create
- `web/index.html` — 基础 HTML 结构（Phase 2 会添加 WebSocket 逻辑）

</code_context>

<specifics>
## Specific Ideas

- 使用 PicoClaw `web/backend/embed.go` 作为 embed pattern 参考
- HTML 基础结构：`<title>Pico Echo Server</title>`, 简洁布局

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 01-embedded-webui-infrastructure*
*Context gathered: 2026-04-16*