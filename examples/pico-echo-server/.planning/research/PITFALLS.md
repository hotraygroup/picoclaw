# Pitfalls Research

**Domain:** Embedded WebSocket WebUI
**Project:** pico-echo-server WebUI
**Researched:** 2026-04-16
**Confidence:** HIGH (verified against PicoClaw web/ implementation patterns)

## Critical Pitfalls

### Pitfall 1: embed.FS Path Traps

**What goes wrong:**
Go's `embed.FS` has subtle path handling issues that cause 404s for static files. The most common failure is embedding files at wrong depth, then serving from incorrect paths.

**Why it happens:**
1. `//go:embed dist` embeds `dist/` as a nested directory — serving requires `fs.Sub()`
2. Hidden files (`.gitkeep`, `.htaccess`) aren't included without `all:` prefix
3. `http.FileServer(http.FS(embedFS))` doesn't strip the embedded prefix automatically

**How to avoid:**
```go
// Pattern 1: Use 'all:' to include hidden files
//go:embed all:dist
var frontendFS embed.FS

// Pattern 2: Extract subdirectory before serving
subFS, err := fs.Sub(frontendFS, "dist")  // "dist" is the embedded root name
if err != nil { log.Fatal(err) }

// Pattern 3: Use http.FS() to wrap for FileServer
fileServer := http.FileServer(http.FS(subFS))
```

**Warning signs:**
- 404s for files that exist in source tree
- `fs.Stat()` returning "file does not exist" for known paths
- Binary size smaller than expected (missing assets)

**Phase to address:** Phase 1 (Embedded WebUI Setup) — must be correct from the start, affects all subsequent development

---

### Pitfall 2: WebSocket Origin Mismatch

**What goes wrong:**
Browser WebSocket connections fail silently or reject when hostname mismatches between page origin and WebSocket URL. Particularly problematic for localhost variants (`localhost`, `127.0.0.1`, `0.0.0.0`).

**Why it happens:**
1. `ws://localhost:9090/ws` from page at `http://127.0.0.1:9090` triggers cross-origin check
2. Cookies are hostname-specific — session cookies won't transfer between localhost variants
3. Development servers often bind to `0.0.0.0` but browser uses `localhost`

**How to avoid:**
```javascript
// Normalize WebSocket URL to match page origin
function normalizeWsUrl(wsUrl) {
  const parsed = new URL(wsUrl)
  const isLocal = ['localhost', '127.0.0.1', '0.0.0.0'].includes(parsed.hostname)
  const browserIsLocal = ['localhost', '127.0.0.1'].includes(window.location.hostname)
  
  if (isLocal && browserIsLocal && parsed.hostname !== window.location.hostname) {
    parsed.hostname = window.location.hostname  // Match the page
  }
  return parsed.toString()
}

// Use relative WebSocket URL (same origin)
const wsScheme = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
const wsUrl = `${wsScheme}//${window.location.host}/ws`
```

**Warning signs:**
- WebSocket connection "opening" but never receiving messages
- CheckOrigin rejects valid localhost connections
- Auth cookies missing on WebSocket upgrade request

**Phase to address:** Phase 2 (WebSocket Integration) — first connection attempt reveals this immediately

---

### Pitfall 3: XSS via Markdown Rendering

**What goes wrong:**
Rendering user-provided Markdown without sanitization enables XSS attacks. `marked.js` and similar libraries output raw HTML — malicious content like `<script>alert('xss')</script>` executes.

**Why it happens:**
Markdown parsers are designed for flexibility, not security. They allow raw HTML by default for features like `<details>` blocks. This is intentional but dangerous for user-generated content.

**How to avoid:**
```javascript
// Option A: Use DOMPurify (recommended for vanilla JS)
import DOMPurify from 'dompurify'
const html = marked.parse(content)
const safeHtml = DOMPurify.sanitize(html, {
  ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'code', 'pre', 'a', 'ul', 'ol', 'li'],
  ALLOWED_ATTR: ['href', 'title']
})
element.innerHTML = safeHtml

// Option B: Disable HTML in marked.js config
marked.setOptions({
  mangle: false,
  headerIds: false,
})
// Still needs sanitization for edge cases

// Option C: Use safe-by-default renderer (react-markdown + rehype-sanitize)
<ReactMarkdown rehypePlugins={[rehypeSanitize]}>{content}</ReactMarkdown>
```

**Warning signs:**
- Raw HTML appearing in rendered output without filtering
- `<script>` or `<iframe>` tags in message content
- Missing sanitization library in dependencies

**Phase to address:** Phase 3 (Markdown Rendering) — security must be baked in from day one, not retrofitted

---

### Pitfall 4: WebSocket Auto-Reconnect Race Conditions

**What goes wrong:**
Naive auto-reconnect creates multiple connections, loses messages, or reconnects when user explicitly disconnected. Stale sockets send messages to nowhere.

**Why it happens:**
1. Multiple `setTimeout` timers fire without deduplication
2. Reconnecting on explicit user disconnect (toggle off)
3. Not tracking connection "generation" to invalidate old callbacks
4. Infinite reconnect loop when server is permanently down

**How to avoid:**
```javascript
let socket = null
let reconnectTimer = null
let reconnectAttempts = 0
let shouldMaintainConnection = false
let generation = 0  // Invalidates stale callbacks

function scheduleReconnect() {
  if (!shouldMaintainConnection || reconnectTimer) return
  
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 5000)  // Exponential backoff, max 5s
  reconnectAttempts++
  
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    if (shouldMaintainConnection) connect()
  }, delay)
}

function onclose() {
  socket = null
  updateState('disconnected')
  scheduleReconnect()
}

function disconnect(explicit) {
  generation++  // Invalidate all callbacks from old generation
  clearTimeout(reconnectTimer)
  reconnectTimer = null
  
  if (explicit) shouldMaintainConnection = false
  
  const oldSocket = socket
  socket = null
  oldSocket?.close()
}
```

**Warning signs:**
- Multiple WebSocket objects in browser DevTools
- Messages appearing multiple times
- "Connecting..." state after explicit disconnect
- Network tab shows rapid retry loop

**Phase to address:** Phase 4 (Auto-Reconnect) — UI polish phase, but foundation must be laid in Phase 2

---

### Pitfall 5: SPA Fallback Hijacking API Routes

**What goes wrong:**
SPA "fallback to index.html" pattern catches `/api/*` routes, returning HTML instead of proper 404 or API response. Debuggers show bizarre errors like "SyntaxError in HTML at line 1".

**Why it happens:**
The catch-all fallback handler (`/` → index.html) is too greedy. It catches everything that isn't a static file, including API routes that should 404 or be handled separately.

**How to avoid:**
```go
mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  // Explicitly exclude API routes from SPA fallback
  if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api" {
    http.NotFound(w, r)  // Let API handler deal with it, or 404 if no handler
    return
  }
  
  // Check for static file
  if _, err := fs.Stat(subFS, cleanPath); err == nil {
    fileServer.ServeHTTP(w, r)
    return
  }
  
  // Fallback to index.html for SPA routes
  indexReq := r.Clone(r.Context())
  indexReq.URL.Path = "/"
  fileServer.ServeHTTP(w, indexReq)
}))
```

**Warning signs:**
- API requests returning HTML content type
- JavaScript parsing errors on JSON API responses
- DevTools showing `<html>` in response to `/api/endpoint`

**Phase to address:** Phase 1 (Embedded WebUI Setup) — route structure must be correct from the start

---

## Technical Debt Patterns

Shortcuts that seem reasonable but create long-term problems.

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Skip MIME type registration | Simpler code | SVG files served as `image/svg` (wrong) instead of `image/svg+xml` | Never — browsers may reject |
| `CheckOrigin: return true` | No CORS headaches | Security vulnerability, any origin can connect | Local dev only, never production |
| Raw `innerHTML` for Markdown | Simple integration | XSS vulnerability, security audit failure | Never — always sanitize |
| Fixed reconnect delay (no backoff) | Predictable behavior | Thundering herd on server restart, DDoS-like retry storm | Never — exponential backoff is mandatory |

## Integration Gotchas

Common mistakes when connecting to external services.

| Integration | Common Mistake | Correct Approach |
|-------------|----------------|------------------|
| `marked.js` CDN | Assume it's XSS-safe by default | DOMPurify must sanitize output |
| WebSocket subprotocol | Not using auth in subprotocol header | Use `new WebSocket(url, ['token.' + authToken])` for auth |
| Browser WebSocket URL | Use absolute URL from config | Build relative URL from `window.location` for same-origin |
| Go embed | Embed root directory directly | Use `fs.Sub()` to extract nested directory |

## Performance Traps

Patterns that work at small scale but fail as usage grows.

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| No reconnect backoff cap | 100+ clients reconnecting simultaneously after restart | Cap at 5s max delay | 10+ concurrent users |
| Stale socket accumulation | Memory growth, ghost connections | Track generation, invalidate old callbacks | Long-running sessions |
| Sync WebSocket send in UI thread | Input lag during slow networks | Use async send with queue | High-latency networks |

## Security Mistakes

Domain-specific security issues beyond general web security.

| Mistake | Risk | Prevention |
|---------|------|------------|
| No Markdown sanitization | XSS, credential theft, session hijacking | DOMPurify with strict whitelist |
| `CheckOrigin: return true` | Cross-origin WebSocket hijacking | Check `r.Header.Get("Origin")` against whitelist |
| Auth token in URL query | Token leaked in logs, browser history | Use WebSocket subprotocol header `['token.' + token]` |
| Missing HTTPS upgrade | Token exposed on network | Enforce `wss://` in production |

## UX Pitfalls

Common user experience mistakes in this domain.

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Silent connection failure | User thinks app is broken | Show "Disconnected" status, attempt reconnect |
| Instant reconnect notification | Spammy UI, flickering | Debounce status changes, show after 2+ failures |
| No send confirmation | Message appears to vanish | Optimistic UI update + send confirmation |
| No typing indicator | Users don't know if message sent | Show "Sending..." briefly before confirmed |

## "Looks Done But Isn't" Checklist

Things that appear complete but are missing critical pieces.

- [ ] **WebSocket connects:** Often missing auth token in subprotocol header — verify auth works on reconnect
- [ ] **Markdown renders:** Often missing XSS sanitization — verify `<script>` tags are stripped
- [ ] **Auto-reconnect works:** Often missing generation tracking — verify explicit disconnect stops reconnect
- [ ] **Static files serve:** Often missing MIME types — verify `.svg` serves as `image/svg+xml`
- [ ] **SPA fallback:** Often missing API exclusion — verify `/api/*` returns proper 404, not HTML

## Recovery Strategies

When pitfalls occur despite prevention, how to recover.

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| embed.FS path wrong | LOW | Adjust `fs.Sub()` path, rebuild binary |
| XSS in Markdown | HIGH | Add DOMPurify immediately, audit for compromised content |
| Race in reconnect | MEDIUM | Add generation tracking, clear timers on disconnect |
| API fallback hijack | LOW | Add prefix check before SPA fallback |

## Pitfall-to-Phase Mapping

How roadmap phases should address these pitfalls.

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| embed.FS Path Traps | Phase 1 (Embedded WebUI) | Test `fs.Stat()` for all static files, verify MIME types |
| WebSocket Origin Mismatch | Phase 2 (WebSocket Integration) | Test from both `localhost` and `127.0.0.1` URLs |
| XSS via Markdown | Phase 3 (Markdown Rendering) | Inject `<script>` in test message, verify it's stripped |
| Auto-Reconnect Race | Phase 4 (Auto-Reconnect) | Toggle connect/disconnect rapidly, verify no duplicate sockets |
| SPA Fallback Hijack | Phase 1 (Embedded WebUI) | Request `/api/nonexistent`, verify 404 not HTML |

## Sources

- PicoClaw `web/backend/embed.go` — Verified embed.FS patterns with `fs.Sub()`
- PicoClaw `web/frontend/src/features/chat/controller.ts` — Verified reconnect with generation tracking
- PicoClaw `web/frontend/src/components/chat/assistant-message.tsx` — Verified `rehype-sanitize` usage
- PicoClaw `web/frontend/src/features/chat/websocket.ts` — Verified URL normalization for localhost
- Go `embed` package documentation — `all:` prefix and `fs.Sub()` requirements
- marked.js documentation — XSS warning, sanitization required

---
*Pitfalls research for: embedded WebSocket WebUI*
*Researched: 2026-04-16*