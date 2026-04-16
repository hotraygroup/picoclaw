# Project Research Summary

**Project:** pico-echo-server WebUI
**Domain:** Embedded WebSocket WebUI for Go server
**Researched:** 2026-04-16
**Confidence:** HIGH

## Executive Summary

This is an embedded WebUI integration for a minimal Go WebSocket echo server. The project transforms a CLI-based stdin input system into a browser-based interface with Markdown rendering. Experts in this domain use Go's native `embed` package for static file distribution, vanilla browser WebSocket API for real-time communication, and CDN-hosted marked.js for lightweight Markdown rendering — achieving single-binary distribution with zero build steps.

The recommended approach is: embed a minimal `web/` directory (index.html, app.js, style.css) using `//go:embed all:web`, serve via `http.FileServer(http.FS(subFS))`, remove the stdin goroutine entirely, and use browser-native WebSocket with relative URL construction (`ws://${location.host}/ws`). The CDN approach for marked.js (43KB) is acceptable for a testing tool and avoids binary bloat.

Key risks are: XSS via unsanitized Markdown (DOMPurify required), WebSocket origin mismatch on localhost variants, and embed.FS path traps causing 404s. These are preventable with established patterns from PicoClaw's existing web/ implementation.

## Key Findings

### Recommended Stack

Use Go embed for static file distribution, marked.js via CDN for Markdown rendering, and vanilla browser WebSocket API — all proven patterns from PicoClaw's existing web/ implementation.

**Core technologies:**
- **Go embed (Go 1.16+, project uses 1.25)**: Embed static files into binary — single-file distribution, zero external dependencies, standard library support
- **marked.js v18.0.0 (CDN)**: Markdown to HTML rendering — lightweight (~43KB minified), zero dependencies, browser-native, actively maintained
- **Native WebSocket API**: Browser WebSocket client — no library needed, native browser support since IE10+, simple API
- **Vanilla HTML/CSS/JS**: Frontend structure — no build step, embed directly, minimal footprint

**CDN URL:** `https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.umd.min.js` (pinned for stability)

### Architecture Approach

The integration requires minimal changes to main.go: use `http.ServeMux` for explicit routing, add `//go:embed all:web` directive, register WebUI routes before WebSocket handler, and **remove the stdin goroutine entirely** (browser becomes primary interface).

**Major components:**
1. **web/index.html** — Main UI page with WebSocket client, textarea input, message display
2. **Embedded route (`/`)** — `http.FileServer(http.FS(subFS))` serving static files from `fs.Sub(webFS, "web")`
3. **WebSocket route (`/ws`)** — Existing handler unchanged, browser connects as new client type

**Route structure:**
- `/` → Embedded WebUI (http.FileServer)
- `/ws` → WebSocket handler (unchanged)

**File changes required:**
- Add imports: `embed`, `io/fs`
- Add directive: `//go:embed all:web`
- Add function: `registerWebUIRoutes(mux)`
- Change: Use `http.NewServeMux()` instead of global `http.HandleFunc`
- Delete: stdin goroutine (10 lines)

**Size estimate:** ~46KB embedded (HTML: 2KB, JS: 3KB, CSS: 1KB, marked: 40KB), final binary ~500KB

### Critical Pitfalls

**Top 5 pitfalls with prevention strategies:**

1. **embed.FS Path Traps** — Use `fs.Sub(webFS, "web")` to extract nested directory before serving. The `//go:embed all:web` directive embeds `web/` as a nested directory; serving requires stripping the prefix. Without this, requests for `index.html` fail because the embedded path is `web/index.html`.

2. **WebSocket Origin Mismatch** — Build WebSocket URL relative to page origin: `const wsUrl = `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws``. This handles localhost variants (`localhost`, `127.0.0.1`, `0.0.0.0`) automatically without cross-origin issues.

3. **XSS via Markdown Rendering** — marked.js outputs raw HTML by default. User content like `<script>alert('xss')</script>` executes. **DOMPurify is required** with strict whitelist: `ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'code', 'pre', 'a', 'ul', 'ol', 'li']`. Alternative: use marked.js CDN (trusting server-side messages only, not user input).

4. **WebSocket Auto-Reconnect Race Conditions** — Use generation tracking and exponential backoff with cap (5s max). Multiple reconnect timers without deduplication cause duplicate connections and message loss.

5. **SPA Fallback Hijacking Routes** — This is a single-page tool; simple `http.FileServer` is sufficient. No complex SPA fallback logic needed. If fallback is added, explicitly exclude `/api/*` routes before serving index.html.

## Implications for Roadmap

Based on research, suggested phase structure follows the natural integration order: static files first (foundation), WebSocket connection second (communication), Markdown third (rendering), polish fourth (reconnect).

### Phase 1: Embedded WebUI Setup
**Rationale:** Foundation for all subsequent work. Static file embedding must be correct before WebSocket or Markdown can function. Incorrect embed.FS setup causes 404s that block all testing.
**Delivers:** Single binary with embedded WebUI, serving index.html at `/`
**Addresses:** Basic UI structure, single-file distribution
**Avoids:** Pitfall 1 (embed.FS path traps), Pitfall 5 (SPA fallback hijacking)

**Implementation notes:**
- Create `web/index.html` with minimal structure
- Add `//go:embed all:web` and `fs.Sub()` pattern
- Register `/` route with `http.FileServer(http.FS(subFS))`
- Change from global `http.HandleFunc` to `http.ServeMux`

### Phase 2: WebSocket Integration
**Rationale:** Browser WebSocket client connects to existing `/ws` handler. Enables bidirectional communication for testing. Depends on Phase 1 for serving the HTML page.
**Delivers:** Browser WebSocket client, message send/receive, basic UI interaction
**Uses:** Vanilla browser WebSocket API, relative URL construction
**Implements:** Browser-side WebSocket client, message framing (Pico Protocol)
**Avoids:** Pitfall 2 (WebSocket origin mismatch)

**Implementation notes:**
- Add WebSocket client in `web/index.html` (or separate app.js)
- Use relative URL: `${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws`
- Handle `message.send` → server → `message.create` flow
- Remove stdin goroutine (browser becomes primary interface)

### Phase 3: Markdown Rendering
**Rationale:** User messages displayed as Markdown. Security critical — XSS must be prevented from day one. Depends on Phase 2 for receiving messages to render.
**Delivers:** Markdown-to-HTML rendering, styled message display
**Uses:** marked.js CDN (18.0.0), optional DOMPurify
**Implements:** `marked.parse()` for incoming messages, styled output container
**Avoids:** Pitfall 3 (XSS via Markdown)

**Implementation notes:**
- Include marked.js CDN: `<script src="https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.umd.min.js"></script>`
- Apply `marked.parse()` to `message.create` payload content
- Add CSS for rendered Markdown (code blocks, lists)
- **Security decision:** If rendering server messages only (trusted), CDN approach is safe. If user input rendered, DOMPurify required.

### Phase 4: Auto-Reconnect (Optional Polish)
**Rationale:** Improves UX for testing scenarios where server restarts. Not essential for core functionality but completes the testing tool experience.
**Delivers:** Automatic reconnect with exponential backoff, connection status indicator
**Avoids:** Pitfall 4 (WebSocket auto-reconnect race conditions)

**Implementation notes:**
- Add generation tracking to invalidate stale callbacks
- Exponential backoff: `Math.min(1000 * Math.pow(2, attempts), 5000)`
- Show connection status in UI ("Connected", "Disconnected", "Reconnecting...")
- Stop reconnect on explicit user disconnect

### Phase Ordering Rationale

1. **Dependency chain:** Embed first (static files) → WebSocket second (communication) → Markdown third (rendering received messages) → Polish fourth (UX improvements)
2. **Architecture grouping:** Phase 1 establishes route structure; Phase 2+ use that structure without modification
3. **Pitfall prevention:** Phase 1 fixes embed.FS (foundation); Phase 2 fixes origin mismatch (immediate failure); Phase 3 fixes XSS (security); Phase 4 fixes reconnect (polish)

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 3 (Markdown):** Security decision needed — CDN only (trusted server messages) vs embed marked.min.js + DOMPurify (user input rendered). Research the use case to determine.

Phases with standard patterns (skip research-phase):
- **Phase 1 (Embedded WebUI):** Well-documented in Go embed docs and PicoClaw `web/backend/embed.go` — pattern is established
- **Phase 2 (WebSocket Integration):** Vanilla browser WebSocket is standard; relative URL construction is documented in PITFALLS.md

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Verified from official Go docs, marked.js docs, CDN verification, existing PicoClaw patterns |
| Architecture | HIGH | Based on existing PicoClaw `web/backend/embed.go` implementation, verified patterns |
| Pitfalls | HIGH | Verified against PicoClaw web/ implementation, official docs for embed and marked.js |

**Overall confidence:** HIGH

All findings verified against:
- Go `embed` package documentation (official)
- marked.js official documentation (official)
- PicoClaw existing web/ implementation (project code)
- CDN jsdelivr file structure (verified)

### Gaps to Address

- **Security decision for Markdown:** Determine if user input is ever rendered (requires DOMPurify) or only server messages (CDN safe). This affects Phase 3 implementation complexity.
- **Auto-reconnect necessity:** Phase 4 is optional polish. Decide during planning if needed for testing tool use case.

## Sources

### Primary (HIGH confidence)
- `https://pkg.go.dev/embed` — Go embed directive syntax, `all:` prefix, `fs.Sub()` requirements
- `https://marked.js.org/` — Official usage docs, browser CDN example, XSS warning
- `https://cdn.jsdelivr.net/npm/marked@18.0.0/` — CDN file structure, version verification
- PicoClaw `web/backend/embed.go` — Existing embed.FS pattern with `fs.Sub()`
- PicoClaw `web/frontend/src/features/chat/websocket.ts` — WebSocket URL normalization

### Secondary (HIGH confidence, verified against project)
- PicoClaw `web/backend/api/pico.go` — WebSocket proxy pattern
- PicoClaw `web/frontend/src/features/chat/controller.ts` — Reconnect with generation tracking
- PicoClaw `web/frontend/src/components/chat/assistant-message.tsx` — `rehype-sanitize` usage

---
*Research completed: 2026-04-16*
*Ready for roadmap: yes*