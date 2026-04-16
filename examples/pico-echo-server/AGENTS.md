<!-- GSD:project-start source:PROJECT.md -->
## Project

**Pico Echo Server WebUI**

一个用于测试 Pico Protocol 的 WebSocket 服务器，添加嵌入式 WebUI 用于可视化 Markdown 格式的消息收发。目标用户是开发者，用于快速测试和调试。

**Core Value:** **轻量、零依赖、一键启动的 Markdown WebSocket 测试工具。**

如果用户能打开浏览器、连接 WebSocket、看到 Markdown 消息渲染，核心价值就达成了。

### Constraints

- **Tech Stack**: Go + embed + pure HTML/JS — no build toolchain, minimal footprint
- **Size**: Stay lightweight (<500KB total) — testing tool should be fast
- **Browser**: Modern browser support only — no legacy IE compatibility
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Recommended Stack
### Core Technologies
| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go embed | Go 1.16+ (project uses 1.25) | Embed static files into binary | Single-file distribution, zero external dependencies, standard library support |
| marked.js | 18.0.0 | Markdown to HTML rendering | Lightweight (~43KB minified), zero dependencies, browser-native, actively maintained |
| Native WebSocket API | Built-in | Browser WebSocket client | No library needed, native browser support since IE10+, simple API |
| Vanilla HTML/CSS/JS | — | Frontend structure | No build step, embed directly, minimal footprint |
### Supporting Libraries
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| gorilla/websocket | 1.5.x | Go WebSocket server (already used) | Project dependency, well-maintained, production-ready |
### CDN URLs for marked.js
| URL | Size | Use Case |
|-----|------|----------|
| `https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.umd.min.js` | ~43KB | Production (Terser minified) |
| `https://cdn.jsdelivr.net/npm/marked/lib/marked.umd.min.js` | ~43KB | Always latest (may break) |
| `https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.esm.js` | — | ES modules (if using `<script type="module">`) |
## Go Embed Patterns
### Pattern 1: Single File Embed (string)
### Pattern 2: Directory Embed (embed.FS)
### Pattern 3: Hybrid (FS with fallback)
## Minimal HTML Template for WebSocket + Markdown
- WebSocket URL construction from `location.host` (works for both localhost and production)
- `marked.parse()` for Markdown rendering (UMD global)
- Auto-scroll to bottom (`scrollTop = scrollHeight`)
- JSON message framing (matches Pico Protocol)
## Auto-Reconnect Pattern
## HTTP Server Setup for Embedded Content
## Alternatives Considered
| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| CDN marked.js | Embed marked.min.js | Offline operation, security policy blocks CDN |
| Vanilla WebSocket | Socket.IO | Need fallback transports, broadcast features, IE support |
| Single HTML file | React/Vue build | Complex UI, component architecture, future extensibility |
| Go embed | External static files | Dynamic content updates, separate deployment |
## What NOT to Use
| Avoid | Why | Use Instead |
|-------|-----|-------------|
| `marked.parse()` without sanitization | XSS vulnerability (marked docs warn) | DOMPurify for user input, or trust server messages only |
| `//go:embed .` | Embeds everything, includes hidden files | Use `//go:embed index.html` or `//go:embed static/*` |
| `http.Handle("/", http.FileServer(...))` without strip | Breaks path resolution | Use `http.StripPrefix` or handle `/` explicitly |
| WebSocket without `CheckOrigin` in dev | CORS rejection | Set `CheckOrigin: func(*http.Request) bool { return true }` for testing |
## Size Considerations
| Component | Size | Notes |
|-----------|------|-------|
| marked.min.js (CDN) | ~43KB | Downloaded by browser, not embedded |
| Single HTML embed | ~3KB | Embedded into binary |
| Total embedded | ~3KB | Well under 500KB target |
| Total with CDN download | ~46KB | Browser fetches marked.js |
## Sources
- `https://marked.js.org/` — Official usage docs, browser CDN example (HIGH confidence)
- `https://cdn.jsdelivr.net/npm/marked@18.0.0/` — CDN file structure, version verification (HIGH confidence)
- `https://pkg.go.dev/embed` — Go embed directive syntax, FS usage (HIGH confidence)
- `/home/sgy/workspace/github/ihotray/picoclaw/web/backend/embed.go` — Existing project pattern reference (HIGH confidence)
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, or `.github/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
