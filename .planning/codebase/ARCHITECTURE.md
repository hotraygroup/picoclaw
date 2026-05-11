<!-- refreshed: 2026-05-11 -->
# Architecture

**Analysis Date:** 2026-05-11

## System Overview

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Presentation Layer                                  │
├────────────────────────────┬────────────────────────────────────────────────┤
│   CLI (`cmd/picoclaw/`)    │   Web Console (`web/`)                         │
│   Cobra commands, TUI      │   React SPA + Go API server                    │
│   (`cmd/picoclaw/main.go`) │   (`web/frontend/`, `web/backend/`)            │
└────────────┬───────────────┴────────────────────┬───────────────────────────┘
             │                                    │
             ▼                                    ▼ (WebSocket via Pico Ch.)
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Gateway / Orchestration Layer                         │
│   `pkg/gateway/gateway.go`                                                   │
│   Lifecycle: init agent → start channels → run event loop → shutdown        │
└─────────────────────────────────────────────────────────────────────────────┘
      │                  │                │                │
      ▼                  ▼                ▼                ▼
┌──────────┐  ┌──────────────────┐  ┌──────────┐  ┌──────────────────┐
│ Channels │  │   Agent Loop     │  │ Cron Svc │  │ Heartbeat /      │
│ Manager  │  │   `pkg/agent/`   │  │ `cron/`  │  │ Health / Devices │
│`channels/`│  │                  │  │          │  └──────────────────┘
└────┬─────┘  └────────┬─────────┘  └──────────┘
     │                 │
     │    ┌────────────┼───────────────┐
     │    ▼            ▼               ▼
     │ ┌──────┐  ┌───────────┐  ┌──────────┐
     │ │Tools │  │ Providers │  │  Context  │
     │ │`tools`│ │`providers`│  │`agent/con-│
     │ └──────┘  └───────────┘  │ text*.go` │
     │                         └───────────┘
     ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  Storage / State Layer                                                       │
│  `pkg/memory/`  `pkg/session/`  `pkg/state/`  `pkg/seahorse/`               │
│  `pkg/skills/`  `pkg/migrate/`  `pkg/config/`                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility | File |
|-----------|----------------|------|
| CLI Entry | Cobra command tree, banner, TZ setup | `cmd/picoclaw/main.go` |
| Gateway | Orchestrates runtime services lifecycle | `pkg/gateway/gateway.go` |
| AgentLoop | Core AI turn processing engine | `pkg/agent/agent.go` |
| Pipeline | Turn execution phases (setup→LLM→finalize→execute) | `pkg/agent/pipeline.go` |
| ContextManager | Token budget, context window management | `pkg/agent/context_manager.go` |
| HookManager | Pre/post processing hooks on agent turns | `pkg/agent/hooks.go` |
| MessageBus | Internal pub/sub message system | `pkg/bus/bus.go` |
| ChannelManager | Multi-platform channel lifecycle, routing, rate limiting | `pkg/channels/manager.go` |
| ProviderFactory | LLM client instantiation (OpenAI, Anthropic, Bedrock, etc.) | `pkg/providers/factory.go` |
| ToolRegistry | Agent tool discovery and invocation | `pkg/tools/registry.go` |
| CommandRegistry | Built-in slash commands (/help, /clear, /switch, etc.) | `pkg/commands/registry.go` |
| Config | Central configuration, versioning, migration, security | `pkg/config/config.go` |
| Web Backend | API server, gateway process control, dashboard auth | `web/backend/main.go` |
| Web Frontend | React SPA for chat, config, logs, agent management | `web/frontend/src/main.tsx` |

## Pattern Overview

**Overall:** Modular Monolith with Plugin Architecture

**Key Characteristics:**
- Single Go binary with compile-time channel registration via `init()` side-effect imports
- Cobra CLI root with subcommands for each operational mode (agent, gateway, onboard, etc.)
- Component-based architecture with clearly separated concerns in `pkg/`
- Adapter pattern for channels (17+ messaging platforms) and providers (6+ LLM backends)
- Pipeline pattern for turn processing: setup → context assembly → LLM call → response finalization
- Event-driven via internal event bus (`pkg/events/`) for runtime observability
- JSON configuration (`config.example.json`) with multi-version migration support

## Layers

**Presentation / Entry:**
- Purpose: User-facing interfaces and command dispatch
- Location: `cmd/picoclaw/main.go` (CLI), `web/backend/` (API server), `web/frontend/` (SPA)
- Contains: Cobra commands, HTTP/SSE APIs, React components and routes
- Depends on: All `pkg/` modules
- Used by: End users, browser clients

**Gateway / Orchestration:**
- Purpose: Service lifecycle management, component wiring, startup sequencing
- Location: `pkg/gateway/gateway.go`
- Contains: Service initialization, event loop, graceful shutdown, panic recovery
- Depends on: agent, channels, cron, heartbeat, health, devices, providers, state, config, bus, events, logger
- Used by: `cmd/picoclaw` gateway subcommand

**Agent / Processing:**
- Purpose: Core AI turn processing, context assembly, LLM interaction, response generation
- Location: `pkg/agent/`
- Contains: AgentLoop, Pipeline, ContextManager, HookManager, TurnCoordinator, Steering, SubTurn
- Depends on: bus, channels (interfaces), providers, commands, routing, session, state, config, memory, seahorse, tools, skills, mcp, media, logger, events
- Used by: Gateway

**Channel Adapters:**
- Purpose: Multi-platform messaging protocol adapters
- Location: `pkg/channels/` and `pkg/channels/{platform}/`
- Contains: Channel interface, Manager, 17+ platform implementations (Discord, Telegram, Slack, QQ, WeChat, Feishu, DingTalk, LINE, VK, WhatsApp, IRC, Matrix, OneBot, MQTT, Teams, MaixCam, Pico)
- Depends on: bus, config, events, logger, media, utils
- Used by: Gateway (via Manager)

**Provider Adapters:**
- Purpose: LLM provider abstraction and fallback chains
- Location: `pkg/providers/`
- Contains: OpenAI, Anthropic, Azure, AWS Bedrock, OAuth, CLI, HTTP API providers; factory, fallback, rate limiting
- Depends on: config, logger
- Used by: Agent

**Tools / Capabilities:**
- Purpose: Agent executable capabilities (shell, fs, subagents, cron, search)
- Location: `pkg/tools/`
- Contains: Shell executor, file system tools, subagent spawning, cron scheduling, search, platform-specific implementations
- Depends on: config, bus, logger, isolation
- Used by: Agent

**Storage / State:**
- Purpose: Persistent state and data management
- Location: `pkg/memory/` (JSONL), `pkg/session/` (session store), `pkg/state/` (runtime state), `pkg/seahorse/` (FTS5 long-term context), `pkg/config/` (configuration)
- Contains: Memory store, session allocator, state manager, seahorse context engine, config struct and migration
- Depends on: SQLite (modernc.org/sqlite), filesystem
- Used by: Agent, Gateway

## Data Flow

### Primary Request Path (Chat Messaging)

1. Inbound message arrives via channel platform API → Channel adapter translates to `bus.InboundMessage` (`pkg/channels/manager.go`)
2. Message published to AgentLoop via bus (`pkg/bus/bus.go` → `pkg/agent/agent.go`)
3. AgentLoop routes to correct agent instance via `routing.Router` (`pkg/routing/router.go`)
4. For each turn: Pipeline setup (`pkg/agent/pipeline_setup.go`) — context assembly, hook pre-processing
5. LLM call via provider fallback chain (`pkg/agent/pipeline_llm.go` → `pkg/providers/factory.go`)
6. Response finalization: post-processing, tool call execution (`pkg/agent/pipeline_finalize.go`)
7. Outbound message published to bus → Channel Manager sends via platform adapter (`pkg/channels/manager.go`)
8. Additional processing: hook post-mount, memory persistence, session update

### Web Console Path (WebSocket Chat)

1. Browser connects via Pico Channel WebSocket (`web/frontend/src/features/chat/websocket.ts`)
2. Messages flow through Pico Channel (`pkg/channels/pico/`) as a standard channel
3. API operations (config, models, logs, skills) go through REST API (`web/backend/api/router.go`)
4. Gateway status/health uses SSE via `/api/gateway/events`

### Startup Sequence

1. `cmd/picoclaw/main.go` → Cobra dispatches to gateway subcommand (`cmd/picoclaw/internal/gateway/`)
2. Gateway loads config (`pkg/config/config.go`), runs migration if needed
3. Initializes logger, event bus, state manager, media store
4. Creates AgentLoop, registers channels, starts channel listeners
5. Starts cron service, heartbeat, health server
6. Enters event loop: listens for inbound messages, OS signals, reload requests

**State Management:**
- Runtime state via `pkg/state/state.go` (JSON file)
- Session state via `pkg/session/manager.go` (allocates per-chat sessions)
- Memory via `pkg/memory/store.go` (JSONL append-only)
- Long-term context via `pkg/seahorse/store.go` (SQLite FTS5)
- Agent turn state via `pkg/agent/turn_state.go` and `sync.Map` fields in AgentLoop

## Key Abstractions

**MessageBus (`pkg/bus/`):**
- Purpose: Internal publish/subscribe for inbound/outbound messages between channels and agent
- Examples: `pkg/bus/bus.go`, `pkg/bus/types.go`, `pkg/bus/events.go`
- Pattern: Channels publish inbound, AgentLoop subscribes; AgentLoop publishes outbound, channels receive via ChannelManager

**Channel Interface (implied by `pkg/channels/`):**
- Purpose: Abstract messaging platform adapter
- Examples: `pkg/channels/telegram/`, `pkg/channels/discord/`, `pkg/channels/slack/`
- Pattern: Each channel registers via `init()` blank import, implements common lifecycle methods (Start, Stop, Name)

**Pipeline (Agent Turn Processing):**
- Purpose: Structured turn execution with hook points
- Examples: `pkg/agent/pipeline.go`, `pkg/agent/pipeline_setup.go`, `pkg/agent/pipeline_llm.go`, `pkg/agent/pipeline_finalize.go`, `pkg/agent/pipeline_execute.go`
- Pattern: Phase-based execution with pipeline wiring, hooks at each phase boundary

**FallbackChain (`pkg/providers/fallback.go`):**
- Purpose: Multi-provider fallback on LLM failures with cooldown management
- Examples: `pkg/providers/fallback.go`, `pkg/providers/cooldown.go`
- Pattern: Iterates configured models, applies cooldown on errors, tries next provider

**Config (`pkg/config/config.go`):**
- Purpose: Single source of truth for all configuration with version-based migration
- Examples: `pkg/config/config.go` (1596 lines), `pkg/config/migration.go`, `pkg/config/security.go`
- Pattern: JSON file with schema versioning, migration from v1→v2→v3, encrypted credential storage

## Entry Points

**CLI (Cobra):**
- Location: `cmd/picoclaw/main.go` → `NewPicoclawCommand()`
- Triggers: User invokes `picoclaw [subcommand]`
- Responsibilities: Banner display, timezone setup, subcommand dispatch

**Gateway:**
- Location: `pkg/gateway/gateway.go`
- Triggers: `picoclaw gateway` subcommand
- Responsibilities: Full runtime service orchestration

**Agent (standalone):**
- Location: `cmd/picoclaw/internal/agent/`
- Triggers: `picoclaw agent` subcommand
- Responsibilities: Direct agent interaction without full gateway

**Web Backend:**
- Location: `web/backend/main.go`
- Triggers: Separate `picoclaw-web` binary launch
- Responsibilities: HTTP API server, frontend dist serving, gateway process management

**Web Frontend:**
- Location: `web/frontend/src/main.tsx` → `web/frontend/src/routeTree.gen.ts`
- Triggers: Browser navigation to `http://localhost:{port}`
- Responsibilities: SPA routing, chat UI, configuration UI, logs viewer

## Architectural Constraints

- **Threading:** Go goroutine-based concurrent model; AgentLoop uses `sync.Map` for per-session turn deduplication; worker semaphore (`workerSem`) limits concurrent turn processing; Channel manager uses per-channel goroutines with bounded queues
- **Global state:** Config singleton via `pkg/config/` package; session allocator via `pkg/session/`; state manager via `pkg/state/`; module-level `var` in `pkg/config/config.go` for version/build info; channel registry via `pkg/channels/registry.go` `init()` side-effects
- **Circular imports:** Not detected — package dependency graph is acyclic; `pkg/agent/interfaces/` exists to break channel↔agent dependency loop
- **Platform constraints:** Conditional compilation via `*_unix.go` / `*_windows.go` suffixes in `pkg/tools/`, `pkg/isolation/`, `pkg/logger/`, `web/backend/`; CGO optional but used for systray on Windows

## Anti-Patterns

### init() Side-Effect Channel Registration

**What happens:** Channel packages (e.g., `pkg/channels/telegram/`) register themselves in a global registry via `init()` functions. Gateway imports them with blank identifiers (`_ "github.com/sipeed/picoclaw/pkg/channels/telegram"`) in `pkg/gateway/gateway.go`.
**Why it's wrong:** Makes channel set compile-time only; adding/removing channels requires code changes and recompilation; import order controls registration order implicitly.
**Do this instead:** This is intentional for a single-binary distribution model, but for dynamic loading consider a configuration-driven plugin approach. Document which files contain the blank imports (`pkg/gateway/gateway.go` lines 22–39).

### Large Config File

**What happens:** `pkg/config/config.go` is 1596 lines and handles config struct, parsing, validation, migration, security, events, channels, models, and diagnostics.
**Why it's wrong:** Single file becomes hard to navigate; changes risk unintended side effects; coupling between config domains.
**Do this instead:** Split into `config_core.go`, `config_models.go`, `config_channels.go`, `config_security.go` following existing partial split pattern used for `config_channel.go` and `config_struct.go`.

### AgentLoop as God Object

**What happens:** `AgentLoop` struct in `pkg/agent/agent.go` holds 20+ fields including bus, config, registry, state, hooks, fallback, channelManager, mediaStore, transcriber, cmdRegistry, mcp, hookRuntime, steering, and various maps.
**Why it's wrong:** High coupling; hard to test individual concerns; lifetime management of all sub-components tied to single struct.
**Do this instead:** Continue extracting Pipeline pattern; consider AgentRuntime composed of sub-services each with their own lifecycle methods.

## Error Handling

**Strategy:** Go idiomatic — errors returned as values, wrapped with context using `fmt.Errorf("...: %w", err)`.

**Patterns:**
- `pkg/channels/errors.go` defines sentinel errors for channel operations
- Provider errors classified in `pkg/providers/error_classifier.go` for fallback decisions
- Gateway panics captured via `pkg/logger/panic.go` and written to `gateway_panic.log`
- CLI errors rendered via `cliui.FormatCLIError` for user-facing display

## Cross-Cutting Concerns

**Logging:** `pkg/logger/` wraps `rs/zerolog`; supports file + console dual output; configurable log levels; 3rd-party log routing via `logger_3rd_party.go`

**Validation:** Config validation at load time in `pkg/config/config.go`; model config validation in `pkg/config/model_config_test.go`; tool argument validation in `pkg/tools/validate.go`

**Authentication:**
- Dashboard auth via `web/backend/dashboardauth/` and `web/backend/middleware/`
- OAuth token management via `pkg/auth/oauth.go` with PKCE support (`pkg/auth/pkce.go`)
- API token auth for web backend via `web/backend/api/auth.go`
- Channel-specific auth managed per-channel (e.g., Telegram bot token, Discord token)

**Security:** Credential encryption via `pkg/credential/store.go`; secure config storage via `pkg/config/security.go`; input sanitization in channel adapters; file system isolation via `pkg/isolation/runtime.go`

---

*Architecture analysis: 2026-05-11*
