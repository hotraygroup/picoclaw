# Codebase Concerns

**Analysis Date:** 2026-05-11

## Critical: Massive Set of Disabled Linters

The `.golangci.yaml` file explicitly disables a large number of critical linters with a TODO comment (line 30: `"TODO: Disabled, because they are failing at the moment, we should fix them and enable (step by step)"`). This means every build currently skips essential static analysis.

**Disabled linters with high impact:**
- **`errcheck`** — Unchecked error returns go undetected across the entire codebase
- **`gosec`** — No security scanning whatsoever (insecure practices, hardcoded credentials, path traversal)
- **`staticcheck`** — No static analysis of correctness, performance, or style
- **`revive`** — No coding style or best-practice enforcement (including context-as-first-arg checks)
- **`forcetypeassert`** — Unsafe type assertions go undetected
- **`errorlint`** — Error wrapping/comparison anti-patterns not caught
- **`contextcheck`** — Context propagation failures not detected
- **`gocritic`** — No opinionated code review comments
- **`ineffassign`** — Dead assignments silently accumulate
- **`unparam`** — Unused function parameters not flagged
- **`exhaustive`** — Missing switch cases in enum-like types not caught
- **`funlen`** — No function length enforcement (several files exceed 1000+ lines)
- **`gocyclo`/`gocognit`** — Cyclomatic/cognitive complexity not measured
- **`testifylint`** — Test assertion anti-patterns not caught
- **`modernize`** — Outdated Go patterns and idioms not flagged

**Files:** `.golangci.yaml` (lines 30-63)

**Impact:** Code reviews cannot rely on automated linting. Bugs that static analysis would catch are slipping through. Security reviews are completely manual.

**Fix approach:** Re-enable linters one at a time, starting with the highest-impact ones (`errcheck`, `staticcheck`, `gosec`). Fix violations incrementally. The linter config already has tuned settings — the linters just need to be un-disabled and violations fixed.

---

## Large / Complex Files (Maintainability Risk)

Several files exceed reasonable size/complexity thresholds, making them difficult to understand, test, and modify:

| File | Lines | Concern |
|------|-------|---------|
| `pkg/tools/integration/web.go` | 2116 | Web scraping/fetching logic — massive single file with many branches |
| `pkg/channels/manager.go` | 1664 | Channel lifecycle management — central dispatch with many edge cases |
| `pkg/seahorse/store.go` | 1642 | Message store for SeaHorse compaction — complex SQLite operations |
| `pkg/config/config.go` | 1596 | Configuration loading — schema, parsing, validation all in one file |
| `pkg/channels/matrix/matrix.go` | 1437 | Matrix protocol integration — long file with dense logic |
| `pkg/channels/telegram/telegram.go` | 1416 | Telegram integration — heavy message handling |
| `pkg/tools/fs/filesystem.go` | 1250 | Filesystem tool — many operation implementations |
| `pkg/channels/feishu/feishu_64.go` | 1239 | Feishu (Lark) integration |
| `pkg/migrate/sources/openclaw/openclaw_config.go` | 1228 | Migration config parsing |
| `pkg/channels/weixin/media.go` | 1157 | WeChat media handling |
| `pkg/tools/shell.go` | 1141 | Shell execution tool — command parsing, PTY management, I/O |
| `pkg/channels/onebot/onebot.go` | 1128 | OneBot protocol integration |
| `pkg/agent/context.go` | 1120 | Context truncation — many compaction strategies in one file |
| `web/backend/api/gateway.go` | 1403 | Gateway API endpoints — large handler file |
| `web/backend/api/skills.go` | 1108 | Skills API — many competing features in one file |

**Impact:** These files are hard to review, test, and modify. A change to one behavior risks breaking another. New contributors struggle to understand them.

**Fix approach:** Split by responsibility. For example, `web.go` could be split into `web_fetch.go`, `web_scrape.go`, `web_markdown.go`. `config.go` could separate parsing from validation from defaults.

---

## Test Coverage Gaps

### Agent Package — Heavily Under-tested

Many core files in `pkg/agent/` lack corresponding test files. These are the core agent loop files that govern the AI agent's behavior:

**Files missing tests:**
- `pkg/agent/adapters/channelmanager.go`
- `pkg/agent/adapters/messagebus.go`
- `pkg/agent/agent_command.go` — command parsing with no tests
- `pkg/agent/agent_event.go` — event generation untested
- `pkg/agent/agent_init.go` — initialization untested
- `pkg/agent/agent_inject.go` — injection logic untested
- `pkg/agent/agent_media.go` — media handling untested
- `pkg/agent/agent_message.go` — message processing untested
- `pkg/agent/agent_options.go` — option handling untested
- `pkg/agent/agent_outbound.go` — outbound message pipeline untested
- `pkg/agent/agent_steering.go` — steering/abort mechanism untested
- `pkg/agent/agent_stop.go` — agent shutdown untested
- `pkg/agent/agent_transcribe.go` — transcription pipeline untested
- `pkg/agent/context_legacy.go` — legacy context manager (deprecated path)
- `pkg/agent/context_usage.go` — usage tracking untested
- `pkg/agent/event_payloads.go` — event structures untested
- `pkg/agent/events.go` — event definitions untested
- `pkg/agent/llm_media.go` — LLM media conversion untested
- `pkg/agent/memory.go` — memory management untested
- `pkg/agent/model_resolution.go` — model resolution untested
- `pkg/agent/pipeline.go` — pipeline struct with `Steering any // TODO` (line 23)
- `pkg/agent/pipeline_execute.go` — pipeline execution with unreachable code risk
- `pkg/agent/pipeline_finalize.go` — finalization untested
- `pkg/agent/pipeline_llm.go` — LLM call pipeline untested
- `pkg/agent/pipeline_setup.go` — pipeline setup untested
- `pkg/agent/prompt_contributors.go` — contributor prompt generation untested
- `pkg/agent/prompt_turn.go` — turn prompt generation untested
- `pkg/agent/turn_context.go` — turn context untested
- `pkg/agent/turn_state.go` — turn state management untested

**Risk:** The core agent behavior (how it processes messages, calls LLMs, manages context windows) lacks test coverage. Regressions in these areas will manifest as production bugs, not test failures.

**Priority:** High — `pipeline_execute.go`, `agent_message.go`, `agent_command.go` are on the critical path for every user interaction.

### Skipped Test

- `web/backend/api/config_test.go` line 757: `t.Skip("TODO: fix this test")` — A test is intentionally skipped, masking a known issue.

---

## Unresolved TODOs in Production Code

**Agent / Pipeline:**
- `pkg/agent/agent.go:271` — Media cleanup disabled: `// TODO: Re-enable media cleanup after inbound media is properly consumed by the agent.` → Media files may accumulate indefinitely.
- `pkg/agent/pipeline.go:23` — `Steering any // TODO: *Steering` — The Steering field is typed as `any` instead of a proper type, losing type safety.
- `pkg/tools/subagent.go:175` — `// TODO(eventbus): once subagents are modeled as child turns` — Subagent event emission is incomplete, event bus consumers won't receive subagent lifecycle events.

**Tools / Integrations:**
- `pkg/tools/integration/mcp_tool.go:452` — `// TODO: Add lifecycle cleanup/retention for MCP artifact files.` → MCP artifact files (`*.txt` in `.artifacts/mcp/`) have no cleanup, causing disk accumulation.
- `pkg/providers/cli/github_copilot_provider.go:29` — `// TODO: Implement stdio mode for GitHub Copilot provider` → stdio mode is not implemented; using `grpc` mode returns an error.

**Events:**
- `pkg/events/subscription.go:233` — `// TODO: replace this with keyed executors when runtime events need per-scope ordering` → Keyed event handling is a stub/no-op.

**Logger:**
- `pkg/logger/logger.go:55` — `TimeFormat: "15:04:05", // TODO: make it configurable???` → Time format is hardcoded.

**ASR (Audio Speech Recognition):**
- `pkg/audio/asr/asr.go:42` — `// TODO: Further restrict this by modelID` → Model filtering is incomplete.

**Providers:**
- `pkg/providers/factory_test.go:104` — `// TODO: Test custom APIBase when createClaudeAuthProvider supports it` — Missing test coverage.
- `pkg/providers/factory_test.go:108` — `// TODO: This test requires openai protocol to support auth_method: "oauth"` — Missing test coverage for OAuth authentication.

---

## Deprecated APIs Still In Use

**Files with deprecated code:**
- `pkg/config/legacy_bindings.go:54-81` — Legacy bindings config is deprecated but still produces migration warnings. Migrated in-memory but not persisted.
- `pkg/config/config.go:786` — `// Deprecated: use registries.github instead.` — Old field still in config schema.
- `pkg/channels/voice_capabilities.go:16` — `// Deprecated: Channels should implement VoiceCapabilityProvider instead.`
- `pkg/agent/steering.go:447` — `// Deprecated: Use HardAbort(sessionKey) for session-safe aborts.` — Old abort path remains.
- `cmd/picoclaw/internal/helpers.go:37-49` — Three deprecated helper functions (`FormatVersion`, `FormatBuildInfo`, `GetVersion`) that duplicate `pkg/config/*` functions.

**Impact:** Deprecated code paths remain in production, risking confusion and inconsistent behavior.

**Fix approach:** Remove or replace deprecated usages, then remove the deprecated code. Add `// Deprecated:` godoc comments for any remaining transitional code.

---

## Panic in Init Functions

Two init functions use `panic()` which will crash the entire process on startup failure:

- `pkg/agent/context_seahorse.go:280` — `panic(fmt.Sprintf("register seahorse context manager: %v", err))`
- `pkg/agent/context_seahorse_unsupported.go:18` — `panic(fmt.Sprintf("register seahorse context manager: %v", err))`

**Impact:** A configuration error or missing dependency during registration will kill the entire process rather than reporting a graceful error.

**Fix approach:** Replace `panic` with error propagation pattern — register via a function returning error, or log fatal with a clear message. Init-time panics prevent the Go runtime from cleanly shutting down goroutines and file handles.

---

## context.Background() in Production Code

`context.Background()` is used extensively outside of tests, preventing proper context cancellation propagation:

**Production uses (examples):**
- `pkg/updater/updater.go:38` — HTTP request without timeout context
- `pkg/events/subscription.go:243` — Event handling falls back to Background context when nil is passed
- `pkg/agent/runtime_event_logger.go:50` — Subscription without timeout
- `pkg/agent/hooks.go:221` — Channel subscription without timeout
- `pkg/agent/pipeline_execute.go:386,412` — Internal contexts with hardcoded timeouts rather than deriving from parent
- `pkg/agent/subturn.go:336` — Sub-turn context created from Background instead of parent
- `pkg/channels/matrix/matrix.go:464,1358,1370` — Multiple Background usages in Matrix channel

**Impact:** Cancellation signals from parent contexts are lost. Shutdown may hang while waiting for operations with no deadline.

**Fix approach:** Pass contexts through the call chain. Use `context.WithTimeout(parent, duration)` instead of `context.WithTimeout(context.Background(), duration)`. Where contexts aren't available, add them as function parameters.

---

## Widespread `return nil, nil` Patterns

112 instances of `return nil, nil` found in the codebase. This pattern can mask errors:

**Files with significant usage:**
- `pkg/seahorse/short_compaction.go` — 6 occurrences in compaction logic
- `pkg/tools/integration/web.go` — 4 occurrences in web fetch/scrape
- `pkg/channels/wecom/wecom.go` — 5 occurrences in message handling
- `pkg/seahorse/short_engine.go` — 3 occurrences
- `pkg/channels/mqtt/mqtt.go` — 3 occurrences
- Many channel implementations — frequent pattern for "ignore this message" cases

**Risk:** Callers that check only the error (e.g., `if err != nil { return err }`) will miss cases where the function returns nil, nil but didn't actually produce a meaningful result. This is especially risky for message handling where silent message drops can occur.

**Fix approach:** Replace `return nil, nil` with either a sentinel error (`ErrSkipped`, `ErrNotApplicable`) or a non-nil zero-value result with a nil error. Ensure callers handle all return value states.

---

## Security Concerns

### Gosec Disabled
As noted above, `gosec` is disabled (` .golangci.yaml:46`). The codebase has no automated security scanning. Manual security review is required for all changes touching external inputs.

### Credential Management
- `pkg/credential/credential.go` — Implements AES-256-GCM encryption for API keys stored as `enc://` values. Well-designed with HKDF key derivation requiring an SSH key. However, the `PassphraseProvider` defaults to reading `PICOCLAW_KEY_PASSPHRASE` from the environment, which can leak into process lists or child processes.
- `web/backend/launcherconfig/password_store.go` — Uses bcrypt with cost 12 for dashboard passwords. Reasonable for its purpose.

### Unsafe Package Usage
`unsafe.Pointer` is used in several files for hardware/OS interaction:
- `pkg/tools/hardware/spi_linux.go` — SPI ioctl calls
- `pkg/tools/hardware/i2c_linux.go` — I2C ioctl calls
- `pkg/tools/hardware/serial_windows.go` — Windows COM port operations
- `pkg/isolation/platform_windows.go` — Windows job object/process isolation

These are legitimate low-level uses, but without `gosec` they won't be flagged if new unsafe code appears in less-appropriate places.

### Environment Variable Exposure
68 instances of `os.Getenv()`/`os.Setenv()` exist. Environment variables can leak secrets into child processes and are visible in `/proc/*/environ` on Linux. The roadmap (ROADMAP.md:36) calls for replacing hardcoded API keys with OAuth flows, which would reduce this surface.

### SSH Key Path Validation
`pkg/credential/credential.go:252-279` — The `allowedSSHKeyPath` function restricts SSH key access to specific directories. This is a good defense, but only applies within the credential package. Other packages reading SSH keys may not have the same protection.

---

## Performance Concerns

### Large File Sizes Suggest Monolithic Design
Several files exceed 1000 lines (see table above). This correlates with functions that do too many things and are hard to optimize or profile.

### Goroutine Spawning Without Clear Lifecycle
124+ goroutines are spawned, many in channel implementations (`pkg/channels/`). Without proper lifecycle management, goroutine leaks can accumulate over the lifetime of a running agent.

### MCP Artifact Leak
`pkg/tools/integration/mcp_tool.go:452` — MCP artifact files have no cleanup mechanism. On long-running instances, the `.artifacts/mcp/` directory will grow unbounded.

### Media Cleanup Disabled
`pkg/agent/agent.go:271` — Media cleanup is commented out. Inbound media files will accumulate on disk.

---

## Fragile Areas

### Shell Execution (`pkg/tools/shell.go`)
- 1141 lines handling shell command execution, PTY management, session multiplexing, and safety guards.
- The command safety guard (line 306) relies on regex pattern matching — could be bypassed with creative command syntax.
- Multiple `_ =` ignored error returns for process kill operations (lines 397, 402).
- Global mutable state: `globalSessionManager` (line 27) with a separate mutex.

### Configuration Loading (`pkg/config/config.go`)
- 1596-line monolithic config file handling parsing, validation, migration, and defaults.
- Changes to config schema are high-risk and hard to test in isolation.

### Channel Manager (`pkg/channels/manager.go`)
- 1664 lines managing all channel lifecycle (start, stop, reconnect, message routing).
- A bug here affects ALL messaging channels (Telegram, Discord, Matrix, WeChat, Slack, etc.).

### SeaHorse Compaction (`pkg/seahorse/short_compaction.go`)
- 898 lines of context window compaction logic.
- Several `return nil, nil` patterns in parsing/decision code.
- Directly affects the quality of LLM responses by controlling what context the AI sees.

### Agent Core (`pkg/agent/agent.go`)
- 640 lines implementing the main agent loop.
- Manages message bus, runtime events, tool execution, LLM interaction, and session state.
- Missing tests for many of its sub-components (see test coverage gaps).

---

## Matrix Channel Documentation Gap

Matrix channel README files in multiple languages have a `## 4. TODO` section:
- `docs/channels/matrix/README.md:75`
- `docs/channels/matrix/README.ja.md:63`
- `docs/channels/matrix/README.zh.md:74`
- `docs/channels/matrix/README.vi.md:63`
- `docs/channels/matrix/README.fr.md:63`
- `docs/channels/matrix/README.pt-br.md:63`

These TODO sections indicate incomplete documentation for Matrix channel features.

---

## Dependency: Discord Fork

`go.mod:143` replaces the official `github.com/bwmarrin/discordgo` with a fork:
```
replace github.com/bwmarrin/discordgo => github.com/yeongaori/discordgo-fork v0.0.0-20260319072544-e8e546f5d532
```
**Risk:** The fork may not track upstream security fixes. If the fork is abandoned, migrating back upstream would require reverting whatever changes the fork provides.

---

## Roadmap Acknowledged Tech Debt

The ROADMAP.md explicitly calls out several areas of intended improvement that represent current gaps:

- **Security hardening** (section 2): Prompt injection defense, tool abuse prevention, SSRF protection, filesystem sandbox, context isolation, and privacy redaction are all planned but not yet implemented.
- **Crypto upgrade** (section 2): Plan to adopt `ChaCha20-Poly1305` for secret storage (currently using AES-256-GCM).
- **OAuth 2.0 flow** (section 2): Deprecate hardcoded API keys in CLI — planned but not done.
- **Provider architecture** (section 3): Refactor from vendor-based to protocol-based classification — in progress.

---

*Concerns audit: 2026-05-11*
