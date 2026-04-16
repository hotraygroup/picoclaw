# Requirements: Pico Echo Server WebUI

**Defined:** 2026-04-16
**Core Value:** 轻量、零依赖、一键启动的 Markdown WebSocket 测试工具

## v1 Requirements

Requirements for WebUI addition. Each maps to roadmap phases.

### Core

- [ ] **CORE-01**: WebUI embedded into binary via Go embed
- [ ] **CORE-02**: HTTP route `/` serves embedded HTML
- [ ] **CORE-03**: WebSocket endpoint `/ws` unchanged
- [ ] **CORE-04**: stdin input removed (WebUI is primary interface)

### WebSocket Client

- [ ] **WS-01**: Browser WebSocket connects to `/ws` with relative URL
- [ ] **WS-02**: Connection uses same-origin (window.location.hostname)
- [ ] **WS-03**: Auto-reconnect on connection loss (exponential backoff)

### Message Handling

- [ ] **MSG-01**: Input form sends message to WebSocket
- [ ] **MSG-02**: Received messages displayed in output area
- [ ] **MSG-03**: Markdown rendered using marked.js (CDN)
- [ ] **MSG-04**: XSS prevented (DOMPurify or trusted content)

## v2 Requirements

Deferred to future. Not in current roadmap.

### Polish

- **PLSH-01**: Connection status indicator (connected/disconnected)
- **PLSH-02**: Message history scroll with timestamps
- **PLSH-03**: Session ID display for multi-client testing

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Connection status UI | Minimal testing tool, not dashboard |
| Message persistence | Ephemeral testing, no storage |
| Session differentiation | Single-view, multi-client not visualized |
| Build toolchain | Stay lightweight, no npm/vite |
| Authentication UI | Token passed via URL param, no login form |

## Traceability

Which phases cover which requirements. Updated during roadmap creation.

| Requirement | Phase | Status |
|-------------|-------|--------|
| CORE-01 | Phase 1 | Pending |
| CORE-02 | Phase 1 | Pending |
| CORE-03 | Phase 1 | Pending |
| CORE-04 | Phase 2 | Pending |
| WS-01 | Phase 2 | Pending |
| WS-02 | Phase 2 | Pending |
| WS-03 | Phase 3 | Pending |
| MSG-01 | Phase 2 | Pending |
| MSG-02 | Phase 2 | Pending |
| MSG-03 | Phase 3 | Pending |
| MSG-04 | Phase 3 | Pending |

**Coverage:**
- v1 requirements: 11 total
- Mapped to phases: 11
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-16*
*Last updated: 2026-04-16 after initial definition*