# Codebase Structure

**Analysis Date:** 2026-04-16

## Directory Layout

```
picoclaw/
├── cmd/                    # CLI entry points
│   ├── picoclaw/           # Main CLI binary
│   ├── picoclaw-launcher-tui/ # TUI launcher variant
│   └── membench/           # Memory benchmarking tool
├── pkg/                    # Core library packages
│   ├── agent/              # AI agent loop and orchestration
│   ├── providers/          # LLM provider implementations
│   ├── channels/           # Messaging channel abstraction
│   ├── tools/              # Tool execution system
│   ├── bus/                # Message bus infrastructure
│   ├── config/             # Configuration system
│   ├── gateway/            # Gateway runtime orchestration
│   ├── session/            # Session management
│   ├── memory/             # Conversation storage (JSONL)
│   ├── mcp/                # Model Context Protocol
│   ├── skills/             # Skills registry and installer
│   ├── routing/            # Light/heavy model routing
│   ├── commands/           # CLI command definitions
│   ├── state/              # State management
│   ├── audio/              # ASR/TTS integration
│   ├── media/              # Media file handling
│   ├── isolation/          # Subprocess isolation
│   ├── auth/               # OAuth/token authentication
│   └── ...                 # Other utility packages
├── web/                    # Web launcher monorepo
│   ├── backend/            # Go HTTP server (WebSocket proxy)
│   └── frontend/           # React SPA (Vite + TanStack Router)
├── config/                 # Example configuration files
├── docs/                   # Documentation (multi-language)
├── scripts/                # Build and release scripts
├── docker/                 # Docker configuration
├── assets/                 # Static assets
└── examples/               # Example projects (echo server)
```

## Directory Purposes

### cmd/picoclaw/
- Purpose: Main CLI entry point and subcommand implementations
- Contains: Cobra command definitions, internal subcommand packages
- Key files: `main.go`, `internal/gateway/`, `internal/agent/`, `internal/auth/`

### pkg/agent/
- Purpose: Core AI agent orchestration
- Contains: AgentLoop, AgentInstance, context management, hooks, steering, subturns
- Key files: `loop.go` (4474 lines), `instance.go`, `context.go`, `hooks.go`, `steering.go`, `subturn.go`

### pkg/providers/
- Purpose: LLM API implementations
- Contains: Provider implementations for 30+ APIs, fallback chain, rate limiting
- Key files: `types.go`, `factory_provider.go`, `fallback.go`, `ratelimiter.go`
- Subdirectories: `anthropic/`, `openai_compat/`, `azure/`, `bedrock/`, `protocoltypes/`

### pkg/tools/
- Purpose: Executable tools for AI actions
- Contains: Tool interface, registry, concrete implementations
- Key files: `base.go`, `registry.go`, `filesystem.go`, `shell.go`, `web.go`, `mcp_tool.go`
- Tools: read_file, write_file, edit_file, shell, web_search, cron, spawn, subagent, skills_*

### pkg/channels/
- Purpose: Messaging platform abstraction
- Contains: Channel interface, Manager, pico channel implementation
- Key files: `base.go`, `manager.go`, `interfaces.go`
- Subdirectory: `pico/` - HTTP-based pico protocol

### pkg/bus/
- Purpose: Decoupled message passing infrastructure
- Contains: MessageBus with typed channels
- Key files: `bus.go`, `types.go`, `inbound_context.go`, `outbound_context.go`

### pkg/config/
- Purpose: Configuration loading and validation
- Contains: Config struct, model list management, security config, sensitive data filtering
- Key files: `config.go` (1383 lines), `security.go`, `model_list.go`

### pkg/gateway/
- Purpose: Runtime orchestration and service lifecycle
- Contains: Gateway Run function, service setup, hot reload
- Key files: `gateway.go` (807 lines), `listen.go`

### pkg/session/
- Purpose: Session key management and allocation
- Contains: SessionStore interface, allocator, key parsing
- Key files: `manager.go`, `allocator.go`, `key.go`, `scope.go`

### pkg/memory/
- Purpose: Conversation history persistence
- Contains: Store interface, JSONL backend, migration
- Key files: `store.go`, `jsonl.go`, `migration.go`

### pkg/skills/
- Purpose: Skills registry and installation
- Contains: GitHub/ClawHub registry, installer, loader
- Key files: `registry.go`, `installer.go`, `loader.go`, `github_registry.go`

### pkg/mcp/
- Purpose: Model Context Protocol integration
- Contains: MCP manager, isolated transport
- Key files: `manager.go`, `isolated_command_transport.go`

### web/backend/
- Purpose: Web launcher HTTP server
- Contains: REST APIs, WebSocket proxy, auth, launcher config
- Key files: `main.go`, `api/`, `middleware/`, `model/`

### web/frontend/
- Purpose: Web launcher SPA
- Contains: React 19 + Vite + TanStack Router
- Key files: `src/main.tsx`, `src/routes/`, `src/components/`, `src/features/`

### docs/
- Purpose: Multi-language documentation
- Contains: User guides, API docs, configuration references
- Subdirectories: `hooks/`, `zh/`, `vi/`, `pt-br/`, `ja/`, `fr/`, `my/`

## Key File Locations

### Entry Points
- `cmd/picoclaw/main.go`: CLI entry point
- `web/backend/main.go`: Web launcher entry point
- `pkg/gateway/gateway.go`: Gateway runtime entry (Run function)

### Configuration
- `pkg/config/config.go`: Main config struct and loading
- `pkg/config/security.go`: Security config (API keys, tokens)
- `pkg/config/model_list.go`: Model-centric provider config

### Core Logic
- `pkg/agent/loop.go`: Agent loop (4474 lines - core processing)
- `pkg/agent/instance.go`: Agent instance creation
- `pkg/providers/factory_provider.go`: Provider factory

### Interface Definitions
- `pkg/providers/types.go`: LLMProvider interface
- `pkg/tools/base.go`: Tool interface
- `pkg/channels/interfaces.go`: Channel interface
- `pkg/memory/store.go`: Memory Store interface

### Testing
- `pkg/agent/loop_test.go`: Agent loop tests (126243 lines)
- `pkg/providers/factory_provider_test.go`: Provider tests
- `pkg/tools/*_test.go`: Tool tests

## Naming Conventions

### Files
- Go source: `*.go` (snake_case)
- Tests: `*_test.go` (co-located with source)
- Build-specific: `*_unix.go`, `*_windows.go`, `*_linux.go`

### Packages
- All lowercase, single word preferred
- Utility packages: `utils`, `logger`, `constants`
- Domain packages: `agent`, `providers`, `tools`, `channels`

### Structs/Types
- PascalCase for public types
- camelCase for private fields
- Interface suffix: `Provider`, `Store`, `Manager`, `Registry`

## Where to Add New Code

### New Tool
- Implementation: `pkg/tools/<tool_name>.go`
- Register: Call `toolsRegistry.Register()` in `pkg/agent/instance.go`
- Config: Add schema in `pkg/config/config.go` under `ToolsConfig`

### New Provider
- Implementation: `pkg/providers/<provider_name>/` or use OpenAI-compatible pattern
- Factory: Add case in `pkg/providers/factory_provider.go` switch
- Protocol meta: Add entry in `protocolMetaByName` map

### New Channel
- Implementation: `pkg/channels/<channel_name>/`
- Interface: Implement `Channel` interface from `pkg/channels/interfaces.go`
- Config: Add channel type and settings in `pkg/config/config.go`

### New Feature (Agent-level)
- Core logic: `pkg/agent/<feature>.go`
- Integration: Connect via `AgentLoop` struct in `loop.go`

### New CLI Command
- Implementation: `cmd/picoclaw/internal/<command>/`
- Register: Add to Cobra tree in `cmd/picoclaw/main.go`

### Utilities
- Shared helpers: `pkg/utils/`
- Constants: `pkg/constants/`

## Special Directories

### .planning/
- Purpose: GSD workflow artifacts
- Contains: PROJECT.md, ROADMAP.md, phase directories
- Generated: Yes (by GSD commands)
- Committed: Yes (workflow state)

### web/backend/dist/
- Purpose: Embedded frontend build output
- Contains: Compiled frontend assets
- Generated: Yes (by `make build-frontend`)
- Committed: Yes (embedded in binary)

### config/
- Purpose: Example configuration templates
- Contains: JSON config examples, security.yml template
- Generated: No
- Committed: Yes

### pkg/providers/protocoltypes/
- Purpose: Shared protocol type definitions
- Contains: Message, ToolCall, LLMResponse structs
- Used across: All providers for consistent API

## Build Variants

### Platform-specific Files
- `*_unix.go`: Unix/Linux-specific implementations
- `*_windows.go`: Windows-specific implementations
- `*_linux.go`: Linux-only features (e.g., I2C, SPI)

### Build Tags
- `goolm`: Default Go OL-m implementation
- `stdjson`: Standard JSON encoding
- `whatsapp_native`: Native WhatsApp support (larger binary)
- `bedrock`: AWS Bedrock provider support

---

*Structure analysis: 2026-04-16*