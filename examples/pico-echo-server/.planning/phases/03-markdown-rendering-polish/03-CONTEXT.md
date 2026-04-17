# Phase 3: Markdown Rendering & Polish - Context

**Gathered:** 2026-04-17
**Status:** Ready for planning

<domain>
## Phase Boundary

添加 Markdown 渲染到 WebUI，服务器消息使用 marked.js 渲染为 HTML。用户消息保持纯文本显示。

**Scope:**
- Add marked.js CDN (v18.0.0) to web/index.html
- Modify appendMessage() to use marked.parse() for server messages
- Keep textContent for user messages (no XSS concern)
- Add CSS for rendered Markdown (code blocks, lists)

**Out of scope:**
- Auto-reconnect (deferred — testing tool doesn't need)
- DOMPurify (deferred — trusting server messages only)
- Connection status UI (deferred — v2 requirement)

</domain>

<decisions>
## Implementation Decisions

### Markdown Rendering
- **D-01:** Use marked.js CDN v18.0.0
  - URL: `https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.umd.min.js`
  - Rationale: Lightweight (~43KB), no build step, CDN loaded

### XSS Handling
- **D-02:** Trust server messages only
  - Server messages: `innerHTML = marked.parse(content)`
  - User messages: `textContent = content` (unchanged from Phase 2)
  - Rationale: Server content is trusted, user input is just echoed back

### Auto-Reconnect
- **D-03:** Skip auto-reconnect
  - Rationale: Testing tool, server restarts are infrequent
  - User can manually refresh page

### CSS for Markdown
- **D-04:** Add basic Markdown styling
  - Code blocks: `pre { background: #f5f5f5; padding: 8px; }`
  - Inline code: `code { background: #f5f5f5; }`
  - Lists: `ul, ol { margin-left: 20px; }`

### the agent's Discretion
- Markdown CSS details (colors, spacing)
- Additional Markdown elements styling

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Markdown Library
- `.planning/research/STACK.md` § CDN URLs for marked.js — pinned version 18.0.0
- `.planning/research/PITFALLS.md` § Pitfall 3: XSS via Markdown Rendering — sanitization options

### Project Context
- `.planning/research/SUMMARY.md` — Phase 3 implementation notes
- `.planning/REQUIREMENTS.md` — WS-03, MSG-03, MSG-04
- `web/index.html` — current HTML (Phase 2)

</canonical_refs>

<code_context>
## Existing Code Insights

### web/index.html (Phase 2)
- appendMessage() uses `textContent` for both user and server
- No Markdown rendering yet
- CSS: basic layout, `.user` and `.server` alignment

### Required Changes
1. Add `<script src="marked.js CDN">` in `<head>`
2. Modify appendMessage():
   ```javascript
   if (type === 'server') {
       div.innerHTML = marked.parse(content);
   } else {
       div.textContent = content;
   }
   ```
3. Add CSS for Markdown elements in `<style>`

</code_context>

<specifics>
## Specific Ideas

- marked.js CDN: `https://cdn.jsdelivr.net/npm/marked@18.0.0/lib/marked.umd.min.js`
- Code block styling: dark background, rounded corners
- Link styling: underline, hover effect

</specifics>

<deferred>
## Deferred Ideas

- DOMPurify sanitization — trusting server messages only
- Auto-reconnect with exponential backoff — testing tool doesn't need
- Connection status indicator — v2 requirement PLSH-01

</deferred>

---

*Phase: 03-markdown-rendering-polish*
*Context gathered: 2026-04-17*