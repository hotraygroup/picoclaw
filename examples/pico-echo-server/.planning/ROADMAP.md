# Roadmap: Pico Echo Server WebUI

**Milestone:** v1 - Embedded WebUI for WebSocket Testing
**Granularity:** Coarse (3 phases)
**Created:** 2026-04-16

## Overview

Transforms pico-echo-server from stdin-based CLI to browser-based WebUI with Markdown rendering. Single binary distribution via Go embed.

## Phases

- [ ] **Phase 1: Embedded WebUI Infrastructure** - Single binary serves WebUI at root path
- [ ] **Phase 2: WebSocket Connection & Messaging** - Browser WebSocket client sends and receives messages
- [ ] **Phase 3: Markdown Rendering & Polish** - Messages render as formatted Markdown with auto-reconnect

---

## Phase Details

### Phase 1: Embedded WebUI Infrastructure

**Goal:** Single binary serves WebUI at root path, preserving WebSocket endpoint

**Depends on:** Nothing (first phase)

**Requirements:** CORE-01, CORE-02, CORE-03

**Success Criteria** (what must be TRUE):
  1. User can open browser to `http://localhost:8080/` and see the WebUI page
  2. Binary size stays under 500KB total (lightweight constraint satisfied)
  3. WebSocket endpoint `/ws` still accepts connections (backward compatibility preserved)

**Plans:** TBD

**UI hint:** yes

---

### Phase 2: WebSocket Connection & Messaging

**Goal:** Browser WebSocket client establishes bidirectional communication with server

**Depends on:** Phase 1

**Requirements:** CORE-04, WS-01, WS-02, MSG-01, MSG-02

**Success Criteria** (what must be TRUE):
  1. User can type message in textarea and click Send button to transmit via WebSocket
  2. Sent messages appear in the output area after server echoes them back
  3. stdin input is no longer required (WebUI is primary interface)

**Plans:** TBD

**UI hint:** yes

---

### Phase 3: Markdown Rendering & Polish

**Goal:** Messages display as formatted Markdown with resilient connection handling

**Depends on:** Phase 2

**Requirements:** WS-03, MSG-03, MSG-04

**Success Criteria** (what must be TRUE):
  1. Markdown content (code blocks, lists, bold/italic) displays correctly formatted
  2. Connection automatically reconnects after server restart with exponential backoff
  3. No XSS vulnerabilities in rendered content (security validated)

**Plans:** TBD

**UI hint:** yes

---

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Embedded WebUI Infrastructure | 0/0 | Not started | - |
| 2. WebSocket Connection & Messaging | 0/0 | Not started | - |
| 3. Markdown Rendering & Polish | 0/0 | Not started | - |

---

## Coverage Validation

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

**Coverage:** 11/11 requirements mapped ✓

---
*Roadmap created: 2026-04-16*
*Ready for planning: `/gsd-plan-phase 1`