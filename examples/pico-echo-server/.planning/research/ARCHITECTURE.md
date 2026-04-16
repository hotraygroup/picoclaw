# Architecture Research

**Domain:** Embedded WebSocket WebUI
**Researched:** 2026-04-16
**Confidence:** HIGH (based on existing PicoClaw patterns)

## Current Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Server (main.go)                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   /ws ───────────────────────────────────── WebSocket       │
│        │                                     Handler         │
│        │    ┌─────────────────────────────┐                 │
│        └───→│  server.handleWS()           │                 │
│             │  - Auth check (Bearer token) │                 │
│             │  - Upgrade to WebSocket      │                 │
│             │  - Message handling loop     │                 │
│             └─────────────────────────────┘                 │
│                                                              │
│   stdin ────→ scanner ────→ broadcast() ────→ all clients   │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Current Components

| Component | Responsibility | Implementation |
|-----------|----------------|----------------|
| `server` struct | Connection management, auth | `sync.Mutex`, `map[*websocket.Conn]string` |
| `handleWS()` | WebSocket upgrade, message routing | gorilla/websocket upgrader |
| `broadcast()` | Send message.create to all clients | Mutex-protected iteration |
| stdin goroutine | Read lines, broadcast to clients | `bufio.Scanner` |

### Current Route Structure

```
/ws  → WebSocket handler (handleWS)
     → Auth: Bearer token (optional)
     → Upgrade: gorilla/websocket.Upgrader
     → Protocol: Pico Protocol (ping, message.send, typing.*)
```

## Proposed Architecture (WebUI Integration)

### System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Server (main.go)                     │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   / ──────────────────────────────────────── WebUI          │
│       │                                      (embedded)      │
│       │    ┌─────────────────────────────┐                  │
│       └───→│  http.FileServer(http.FS)   │                  │
│            │  - Serve index.html          │                  │
│            │  - Serve static assets       │                  │
│            └─────────────────────────────┘                  │
│                                                              │
│   /ws ───────────────────────────────────── WebSocket       │
│        │                                     Handler         │
│        │    ┌─────────────────────────────┐                 │
│        └───→│  server.handleWS()           │                 │
│             │  - Auth check (Bearer token) │                 │
│             │  - Upgrade to WebSocket      │                 │
│             │  - Message handling loop     │                 │
│             └─────────────────────────────┘                 │
│                                                              │
│   stdin ────→ (REMOVED)                                    │
│                                                              │
│   Browser ────→ WebSocket ────→ handleWS ────→ broadcast    │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Proposed Route Structure

| Route | Handler | Purpose |
|-------|---------|---------|
| `/` | `http.FileServer(http.FS(webFS))` | Serve embedded WebUI |
| `/ws` | `server.handleWS` | WebSocket connections (unchanged) |
| `/ws` (browser) | Same handler | Browser WebSocket client |

### WebUI Message Flow

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│    Browser UI    │────→│   WebSocket      │────→│   server.handleWS │
│                  │     │   Connection     │     │                  │
│  [textarea]      │     │                  │     │  - Validate      │
│  [send button]   │     │  message.send    │     │  - Log to stdout │
│                  │     │                  │     │                  │
│  [message list]  │←────│  message.create  │←────│  - Broadcast     │
│  (Markdown)      │     │                  │     │                  │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

## Recommended Project Structure

```
pico-echo-server/
├── main.go               # Core server (existing, modified)
├── web/                  # Embedded WebUI files
│   ├── index.html        # Main UI page
│   ├── app.js            # WebSocket client logic
│   └── marked.min.js     # Markdown renderer (CDN or local)
│   └── style.css         # Minimal styling
├── README.md             # Documentation
└── .planning/            # GSD planning directory
```

### Structure Rationale

- **web/ subdirectory:** Follows Go embed convention, cleanly separates static assets from server code
- **Single files (not nested):** Minimal complexity, no build step needed
- **marked.min.js local:** Avoids CDN dependency, keeps binary self-contained

## File Changes Required

### main.go Modifications

```go
// Add at package level
import (
    "embed"
    "io/fs"
    "net/http"
    // ... existing imports
)

//go:embed all:web
var webFS embed.FS

// Add before http.HandleFunc("/ws", ...)
func registerWebUIRoutes(mux *http.ServeMux) {
    subFS, _ := fs.Sub(webFS, "web")
    fileServer := http.FileServer(http.FS(subFS))
    
    mux.Handle("/", fileServer)
}

// In main():
// Change from http.HandleFunc to http.ServeMux for proper routing
mux := http.NewServeMux()
registerWebUIRoutes(mux)
mux.HandleFunc("/ws", s.handleWS)
log.Fatal(http.ListenAndServe(*addr, mux))
```

### stdin goroutine removal

The stdin goroutine (lines 147-157) should be **removed**:
- WebUI becomes the primary interface
- Browser sends messages via WebSocket (`message.send`)
- Server logs received messages to stdout (existing line 99)

## Architecture Patterns

### Pattern 1: Go embed for Static Files

**What:** Embed static files directly into binary using `//go:embed all:web`
**When:** Single-file distribution, no external web server needed
**Trade-offs:** 
- Pro: No external file dependencies, single binary deployment
- Pro: Works offline, no CDN dependencies
- Con: Binary size increase (typically ~50-100KB for minimal WebUI)

**Example:**
```go
//go:embed all:web
var webFS embed.FS

func registerWebUIRoutes(mux *http.ServeMux) {
    subFS, _ := fs.Sub(webFS, "web")
    mux.Handle("/", http.FileServer(http.FS(subFS)))
}
```

### Pattern 2: http.ServeMux Route Priority

**What:** Use `http.ServeMux` for explicit route registration order
**When:** Multiple routes on same server (static + WebSocket)
**Trade-offs:**
- Pro: Explicit route precedence (longest path match wins)
- Pro: Cleaner than `http.HandleFunc` global registration
- Con: Slightly more verbose

**Example:**
```go
mux := http.NewServeMux()
mux.Handle("/", webUIHandler)       // Catch-all for WebUI
mux.HandleFunc("/ws", wsHandler)    // Specific path for WebSocket
```

### Pattern 3: Browser WebSocket Client

**What:** Vanilla JS WebSocket client with auto-reconnect
**When:** Testing tool, no framework overhead needed
**Trade-offs:**
- Pro: Zero dependencies, minimal code
- Pro: Native browser API, fast
- Con: Manual reconnect logic

**Example:**
```javascript
let ws = new WebSocket(`ws://${location.host}/ws`);
ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    if (msg.type === 'message.create') {
        renderMarkdown(msg.payload.content);
    }
};
ws.onclose = () => { setTimeout(reconnect, 1000); };
```

## Integration with Existing Code

### Minimal Changes to server struct

**No changes needed to `server` struct** - it already handles:
- Connection management (`conns map`)
- Authentication (`token` field)
- Broadcasting (`broadcast()` method)

### Changes Required

| Component | Change | Impact |
|-----------|--------|--------|
| `main()` | Use `http.ServeMux` instead of global `http.HandleFunc` | 5 lines |
| Imports | Add `embed`, `io/fs` | 2 imports |
| Package | Add `//go:embed all:web` directive | 1 line |
| stdin goroutine | Remove entirely | Delete 10 lines |

### Total Integration Complexity

- **Additions:** ~15 lines (embed setup, route registration)
- **Deletions:** ~10 lines (stdin goroutine)
- **Net change:** +5 lines to existing 160-line file

## Data Flow

### Browser → Server → Browser

```
[Browser Input]
    ↓ (user types, clicks Send)
[app.js] sends message.send
    ↓ (WebSocket)
[handleWS] validates, logs to stdout
    ↓ (existing line 99)
[broadcast] sends message.create to ALL clients
    ↓ (WebSocket)
[app.js] receives message.create
    ↓
[marked.js] renders Markdown
    ↓
[DOM] displays in message list
```

### Protocol Flow

| Message Type | Direction | Purpose |
|--------------|-----------|---------|
| `message.send` | Browser → Server | User sends message |
| `message.create` | Server → Browser | Broadcast to all |
| `ping` | Browser → Server | Keepalive (optional) |
| `pong` | Server → Browser | Keepalive response |

## Anti-Patterns to Avoid

### Anti-Pattern 1: Using global http.HandleFunc

**What people do:** `http.HandleFunc("/", handler)` mixed with `http.HandleFunc("/ws", ...)`
**Why it's wrong:** Global registration order matters, unclear precedence
**Do this instead:** Use `http.ServeMux` for explicit route registration

### Anti-Pattern 2: CDN-only JavaScript

**What people do:** `<script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js">`
**Why it's wrong:** Binary not self-contained, requires internet at runtime
**Do this instead:** Embed marked.min.js in `web/` directory

### Anti-Pattern 3: SPA Fallback Logic (Overkill)

**What people do:** Complex index.html fallback for all unmatched routes
**Why it's wrong:** This is a single-page tool, no routing needed
**Do this instead:** Simple `http.FileServer` with single index.html

## Size Estimates

| Component | Size |
|-----------|------|
| index.html | ~2KB |
| app.js | ~3KB |
| style.css | ~1KB |
| marked.min.js | ~40KB |
| **Total embedded** | ~46KB |
| **Binary overhead** | ~50KB |
| **Final binary** | ~500KB (Go runtime + WebSocket + WebUI) |

**Result:** Well under the 500KB constraint.

## Sources

- PicoClaw `web/backend/embed.go` — Reference embed pattern (verified)
- PicoClaw `web/backend/api/pico.go` — WebSocket proxy pattern (verified)
- Go embed documentation — `//go:embed` directive (official)
- gorilla/websocket — Standard Go WebSocket library (verified)

---
*Architecture research for: pico-echo-server WebUI integration*
*Researched: 2026-04-16*