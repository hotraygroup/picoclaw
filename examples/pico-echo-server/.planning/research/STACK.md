# Stack Research

**Domain:** Embedded WebUI for Go WebSocket server
**Researched:** 2026-04-16
**Confidence:** HIGH (verified from official docs, existing project patterns, CDN verification)

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

**Recommendation:** Use pinned version `marked@18.0.0` for stability.

## Go Embed Patterns

### Pattern 1: Single File Embed (string)

```go
import _ "embed"

//go:embed index.html
var indexHTML string

func handleIndex(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(indexHTML))
}
```

**When:** Single HTML file, no external assets. Simplest approach.

### Pattern 2: Directory Embed (embed.FS)

```go
import "embed"
import "io/fs"
import "net/http"

//go:embed static/*
var staticFS embed.FS

func setupRoutes(mux *http.ServeMux) {
    // Strip "static/" prefix to serve files at root
    subFS, _ := fs.Sub(staticFS, "static")
    mux.Handle("/", http.FileServer(http.FS(subFS)))
}
```

**When:** Multiple files (HTML, CSS, JS). Full static directory structure.

### Pattern 3: Hybrid (FS with fallback)

```go
//go:embed all:static
var staticFS embed.FS

func handleSPA(w http.ResponseWriter, r *http.Request) {
    subFS, _ := fs.Sub(staticFS, "static")
    
    // Try exact file first
    path := r.URL.Path
    if _, err := fs.Stat(subFS, path); err == nil {
        http.FileServer(http.FS(subFS)).ServeHTTP(w, r)
        return
    }
    
    // Fallback to index.html for SPA routes
    r.URL.Path = "/index.html"
    http.FileServer(http.FS(subFS)).ServeHTTP(w, r)
}
```

**When:** SPA with client-side routing (not needed for this simple project).

## Minimal HTML Template for WebSocket + Markdown

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Pico Echo Server</title>
    <script src="https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.umd.min.js"></script>
    <style>
        body { font-family: system-ui, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        #messages { border: 1px solid #ccc; padding: 10px; height: 400px; overflow-y: auto; }
        #messages pre { background: #f5f5f5; padding: 8px; border-radius: 4px; }
        #input-area { margin-top: 10px; display: flex; gap: 10px; }
        textarea { flex: 1; resize: vertical; min-height: 60px; }
        button { padding: 8px 16px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>Pico Echo Server</h1>
    <div id="messages"></div>
    <div id="input-area">
        <textarea id="input" placeholder="Type Markdown message..."></textarea>
        <button id="send">Send</button>
    </div>

    <script>
        const ws = new WebSocket(`${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws`);
        const messages = document.getElementById('messages');
        const input = document.getElementById('input');
        const sendBtn = document.getElementById('send');

        ws.onopen = () => appendMessage('Connected', 'system');
        ws.onclose = () => appendMessage('Disconnected', 'system');
        ws.onerror = (e) => appendMessage('Error: ' + e.type, 'error');
        
        ws.onmessage = (e) => {
            const msg = JSON.parse(e.data);
            if (msg.type === 'message.create') {
                appendMessage(msg.payload.content, 'received');
            }
        };

        function appendMessage(content, type) {
            const div = document.createElement('div');
            div.className = type;
            if (type === 'received') {
                div.innerHTML = marked.parse(content);
            } else {
                div.textContent = content;
            }
            messages.appendChild(div);
            messages.scrollTop = messages.scrollHeight;
        }

        sendBtn.onclick = () => {
            const content = input.value.trim();
            if (!content) return;
            
            ws.send(JSON.stringify({
                type: 'message.send',
                timestamp: Date.now(),
                payload: { content }
            }));
            input.value = '';
        };
    </script>
</body>
</html>
```

**Key patterns:**
- WebSocket URL construction from `location.host` (works for both localhost and production)
- `marked.parse()` for Markdown rendering (UMD global)
- Auto-scroll to bottom (`scrollTop = scrollHeight`)
- JSON message framing (matches Pico Protocol)

## Auto-Reconnect Pattern

```javascript
let ws;
let reconnectDelay = 1000;

function connect() {
    ws = new WebSocket(`${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws`);
    
    ws.onclose = () => {
        appendMessage('Disconnected, reconnecting...', 'system');
        setTimeout(connect, reconnectDelay);
        reconnectDelay = Math.min(reconnectDelay * 2, 30000); // Cap at 30s
    };
    
    ws.onopen = () => {
        reconnectDelay = 1000; // Reset delay
        appendMessage('Connected', 'system');
    };
    
    // ... other handlers
}

connect();
```

## HTTP Server Setup for Embedded Content

```go
// Simple single-file approach
func setupServer() *http.ServeMux {
    mux := http.NewServeMux()
    
    // WebSocket handler (existing)
    mux.HandleFunc("/ws", s.handleWS)
    
    // Embedded index.html
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        w.Write([]byte(indexHTML))
    })
    
    return mux
}
```

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

**Recommendation:** Use CDN for marked.js. 43KB is acceptable for testing tool, and avoids binary bloat.

## Sources

- `https://marked.js.org/` — Official usage docs, browser CDN example (HIGH confidence)
- `https://cdn.jsdelivr.net/npm/marked@18.0.0/` — CDN file structure, version verification (HIGH confidence)
- `https://pkg.go.dev/embed` — Go embed directive syntax, FS usage (HIGH confidence)
- `/home/sgy/workspace/github/ihotray/picoclaw/web/backend/embed.go` — Existing project pattern reference (HIGH confidence)

---
*Stack research for: Embedded WebUI with WebSocket + Markdown*
*Researched: 2026-04-16*