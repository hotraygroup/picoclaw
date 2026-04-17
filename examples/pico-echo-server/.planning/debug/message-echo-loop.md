---
status: resolved
trigger: 消息死循环了，自己发送的，又发回给自己了
created: 2026-04-17
updated: 2026-04-17
---

# Debug Session: message-echo-loop

## Symptoms

**Expected:** Messages received from client A should NOT be sent back to A, only broadcast to OTHER clients.

**Actual:** Client A sends message → server broadcasts to ALL clients → client A receives its own message back.

**Repro:** Always. Every message.send triggers echo to sender.

**Timeline:** Started immediately after Phase 2 completion (s.broadcast added).

**Error:** None. Functional bug, not crash.

## Current Focus

**Hypothesis:** `main.go:102 s.broadcast(content)` sends message.create to ALL connected clients without excluding the sender.

**Test:** Check broadcast() function - does it iterate over ALL conns?

**Expecting:** broadcast() iterates over s.conns map without filtering.

**Next action:** Confirm broadcast() iterates over all conns, then modify to exclude sender.

**Reasoning checkpoint:** 
- Phase 2 added s.broadcast(content) to echo messages
- broadcast() sends to ALL s.conns (lines 125-130)
- No exclusion logic for the sender conn
- Browser receives its own message via ws.onmessage

## Evidence

- timestamp: 2026-04-17T02:00:00Z
  observation: main.go:102 calls s.broadcast(content) on message.send
  source: code review

- timestamp: 2026-04-17T02:00:30Z
  observation: broadcast() iterates over ALL s.conns without exclusion
  source: main.go:125-130

## Eliminated

(none yet)

## Resolution

**Root cause:** `broadcast()` function iterated over ALL `s.conns` without excluding the sender connection, causing every client (including the message sender) to receive the echoed message.

**Fix:** 
1. Renamed `broadcast()` to `broadcastExcept(content, senderConn)` with sender exclusion parameter
2. Added `if conn == senderConn { continue }` in iteration loop
3. Updated call at line 102: `s.broadcastExcept(content, conn)`

**Verification:** `go build` succeeds, grep confirms exclusion logic present.

**Files changed:** main.go (lines 102, 116-133)