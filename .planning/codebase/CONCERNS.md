# Codebase Concerns

**Analysis Date:** 2026-04-16

## Tech Debt

**Agent Loop Complexity:**
- Issue: `pkg/agent/loop.go` is 4474 lines - extremely large single file containing core agent logic
- Files: `pkg/agent/loop.go`
- Impact: Difficult to maintain, understand, and modify; risk of bugs when changes are made
- Fix approach: Split into smaller modules (turn handling, context management, tool orchestration, event handling)

**Seahorse Store Size:**
- Issue: `pkg/seahorse/store.go` is 1593 lines - large SQLite-based storage implementation
- Files: `pkg/seahorse/store.go`
- Impact: Complex summary/conversation storage logic difficult to navigate
- Fix approach: Extract operations into separate files by domain (conversations, summaries, messages)

**Shell Tool Extensive Patterns:**
- Issue: `pkg/tools/shell.go` has 48+ default deny patterns for dangerous commands
- Files: `pkg/tools/shell.go:50-98`
- Impact: Complexity in maintaining deny patterns; potential bypasses if patterns not comprehensive
- Fix approach: Consider allowlist approach for common safe commands; document pattern rationale

**Test Coverage by File Count:**
- Issue: 225 test files vs 349 source files (~64% coverage by file count)
- Files: Global
- Impact: Some packages may have insufficient test coverage; areas without tests could have hidden bugs
- Fix approach: Review test coverage per package; add tests for under-tested critical paths

## Known Bugs

**Windows Isolation Limitation:**
- Symptoms: Windows subprocess isolation does not support `expose_paths` filesystem rules
- Files: `pkg/isolation/runtime.go:300-305`, `pkg/isolation/platform_windows.go`
- Trigger: Using `expose_paths` in config on Windows
- Workaround: Windows uses restricted-token approach only; no filesystem isolation available

**Login Rate Limiting Per-IP Only:**
- Symptoms: Rate limiting (10 attempts/minute) is per-IP only, not per-user
- Files: `web/backend/api/auth_login_limiter.go:12-14`
- Trigger: Multiple users behind same NAT/proxy can exhaust shared rate limit
- Workaround: Consider session-based limiting for multi-user scenarios

## Security Considerations

**Credential Storage:**
- Risk: API keys stored in `.security.yml` file with 0600 permissions; plaintext in memory after loading
- Files: `pkg/config/security.go`, `pkg/credential/credential.go`
- Current mitigation: 
  - Separation from main config (safe to share config.json)
  - Optional AES-256-GCM encryption with SSH key + passphrase two-factor
  - Sensitive data filtering prevents LLM seeing its own credentials
  - `enc://` format for encrypted keys
- Recommendations:
  - Always use `enc://` format for production
  - Never commit `.security.yml` to version control (add to `.gitignore`)
  - Consider hardware key storage for high-security environments

**Shell Command Execution:**
- Risk: AI agent can execute arbitrary shell commands through tools
- Files: `pkg/tools/shell.go:50-115`
- Current mitigation:
  - 48+ deny patterns blocking dangerous commands (rm -rf, dd, format, sudo, etc.)
  - Command substitution blocked ($(), `${}`, backticks)
  - Pipe to shell blocked (`| sh`, `| bash`)
  - Workspace restriction option
  - Custom allow/deny patterns configurable
  - Timeout support
- Recommendations:
  - Enable `restrict_to_workspace` by default
  - Review custom deny patterns for edge cases
  - Consider sandboxing for production deployments

**Filesystem Path Traversal:**
- Risk: Path traversal through symlinks, relative paths, or absolute paths
- Files: `pkg/tools/filesystem.go:27-85`
- Current mitigation:
  - Symlink resolution with containment check
  - Workspace boundary enforcement
  - Allowed path patterns for exceptions
  - `filepath.EvalSymlinks` before containment check
- Recommendations:
  - Test symlink edge cases thoroughly
  - Validate path patterns with regex carefully

**Authentication Token Handling:**
- Risk: Session cookies and tokens stored in memory
- Files: `web/backend/api/auth.go:150-153`
- Current mitigation:
  - Constant-time comparison for token validation (`subtle.ConstantTimeCompare`)
  - bcrypt password hashing when PasswordStore enabled
  - Login rate limiting (10 attempts/minute per IP)
  - Secure cookie detection based on request
- Recommendations:
  - Ensure secure cookies in production (HTTPS)
  - Consider token expiration/rotation
  - Add audit logging for authentication events

**AI-Generated Code Scrutiny:**
- Risk: Project embraces AI-assisted contributions; AI-generated code may have subtle bugs
- Files: `CLAUDE.md` (explicit guidance)
- Current mitigation: Documentation requires extra scrutiny for path traversal, injection, credential exposure
- Recommendations:
  - Manual review required for all security-sensitive code
  - Automated linting with security-focused rules
  - Consider formal code review process for AI-generated changes

**Isolation Subprocess Security:**
- Risk: Child processes may escape isolation or access unintended resources
- Files: `pkg/isolation/runtime.go`, `pkg/isolation/platform_linux.go`, `pkg/isolation/platform_windows.go`
- Current mitigation:
  - Linux: bubblewrap-style filesystem isolation with mount rules
  - Windows: restricted token approach (limited)
  - User environment redirection (HOME, TMP, etc.)
  - Expose path validation
- Recommendations:
  - Linux isolation preferred for production
  - Test isolation boundaries thoroughly
  - Document what isolation does NOT protect against

## Performance Bottlenecks

**Memory Footprint:**
- Problem: Target is <10MB RAM for minimal hardware ($10 devices)
- Files: `docs/hardware-compatibility.md:117-123`
- Cause: Design constraint for embedded devices (Raspberry Pi Zero, routers)
- Improvement path: Monitor memory usage; optimize JSONL storage; consider streaming for large responses

**JSONL Session Storage:**
- Problem: JSONL files grow unbounded; no automatic cleanup
- Files: `pkg/session/jsonl_backend.go`, `pkg/memory/`
- Cause: Fire-and-forget write pattern for session history
- Improvement path: Add session expiration/cleanup; implement compaction; consider SQLite for structured queries

**Credential Decryption Overhead:**
- Problem: Each `enc://` key requires ~2ms decryption at startup
- Files: `pkg/credential/credential.go:89-95`
- Cause: AES-256-GCM + HKDF-SHA256 key derivation
- Improvement path: Acceptable for startup (<2ms per key); cache decrypted values in memory

**Large File Compilation:**
- Problem: `loop.go` (4474 lines) and `store.go` (1593 lines) increase compilation time
- Files: `pkg/agent/loop.go`, `pkg/seahorse/store.go`
- Cause: Single-file architecture for complex modules
- Improvement path: Split into smaller files; improve Go's parallel compilation efficiency

## Fragile Areas

**Provider Fallback Chain:**
- Files: `pkg/providers/fallback.go`, `pkg/providers/cooldown.go`
- Why fragile: Complex interaction between rate limiting, cooldown tracking, and fallback logic
- Safe modification: Test with multiple providers; verify cooldown doesn't block valid requests
- Test coverage: Has tests but edge cases (concurrent cooldown, rate limit overflow) need coverage

**Config Version Migration:**
- Files: `pkg/config/config.go:25-26`, `pkg/config/migration_integration_test.go`, `pkg/config/migration.go`
- Why fragile: Config versioning (v0 → v3) requires careful migration for existing users
- Safe modification: Add migration tests for each version; preserve backward compatibility
- Test coverage: Integration tests exist; edge cases (partial migration, corrupted config) need attention

**Sensitive Data Replacer:**
- Files: `pkg/config/security.go:205-244`
- Why fragile: Reflection-based traversal of config struct; depends on field types
- Safe modification: Test with new config types; verify SecureString/SecureStrings handling
- Test coverage: `security_test.go` exists; reflection edge cases need coverage

## Scaling Limits

**Session Management:**
- Current capacity: In-memory sessions with optional JSONL persistence
- Limit: Memory usage grows with active sessions; JSONL files grow without bounds
- Scaling path: 
  - Implement session expiration/cleanup
  - Consider SQLite backend for sessions (already used for seahorse summaries)
  - Add pagination for session listing

**Provider Rate Limiting:**
- Current capacity: Per-model RPM limits via token bucket
- Limit: Token bucket per model; no global rate limit across providers
- Scaling path: Add global rate limit for API quota management; implement request prioritization

**Web Tool API Key Pool:**
- Current capacity: API key rotation for failover
- Limit: Round-robin rotation; no smart selection based on quota remaining
- Scaling path: Implement quota-aware key selection; add per-key usage tracking

## Dependencies at Risk

**Go Version Requirement:**
- Risk: Requires Go 1.25+ (very recent)
- Impact: Users with older Go installations cannot build
- Migration plan: Document as requirement; provide pre-built binaries

**External SDK Dependencies:**
- Risk: Heavy dependency on external SDKs (Anthropic, OpenAI, AWS Bedrock, Discord, Slack, etc.)
- Impact: API changes in upstream SDKs may break functionality
- Migration plan: Pin versions; monitor SDK changelogs; implement version compatibility checks

**WhatsApp Native Build Tag:**
- Risk: `whatsapp_native` build tag creates larger binary with whatsmeow library
- Impact: Binary size increase for users needing WhatsApp support
- Migration plan: Default build excludes WhatsApp; document build tag usage

**SQLite Dependency:**
- Risk: Uses `modernc.org/sqlite` (pure Go SQLite) for seahorse summaries
- Impact: No external SQLite library needed; pure Go implementation
- Migration plan: Modernc SQLite is stable; monitor for updates

## Missing Critical Features

**No Global Rate Limit:**
- Problem: No global request rate limit across all providers
- Blocks: Users with limited total API quota cannot manage across multiple providers
- Priority: Medium

**Session Cleanup:**
- Problem: No automatic session expiration or cleanup
- Blocks: Long-running instances accumulate unbounded session history
- Priority: High

**Windows Full Isolation:**
- Problem: Windows subprocess isolation limited (no filesystem isolation)
- Blocks: Windows users cannot fully sandbox child processes
- Priority: Medium (Windows not primary target)

## Test Coverage Gaps

**Shell Tool Edge Cases:**
- What's not tested: Bypass scenarios for deny patterns; concurrent session management
- Files: `pkg/tools/shell_test.go`
- Risk: Untested bypass patterns could allow dangerous command execution
- Priority: High

**Isolation Platform Tests:**
- What's not tested: Windows isolation behavior; Linux mount rule edge cases
- Files: `pkg/isolation/platform_linux_test.go`, `pkg/isolation/runtime_test.go`
- Risk: Platform-specific isolation bugs may go undetected
- Priority: Medium

**Concurrent Fallback Tests:**
- What's not tested: Concurrent requests hitting cooldown/rate limit simultaneously
- Files: `pkg/providers/fallback_test.go`
- Risk: Race conditions in fallback chain may cause unexpected behavior
- Priority: Medium

**Credential Resolution Edge Cases:**
- What's not tested: Malformed `file://` paths; symlink escaping; missing SSH key scenarios
- Files: `pkg/credential/credential_test.go`
- Risk: Invalid input handling may cause crashes or security issues
- Priority: High

---

*Concerns audit: 2026-04-16*