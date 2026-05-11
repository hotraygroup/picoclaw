# Codebase Structure

**Analysis Date:** 2026-05-11

## Directory Layout

```
picoclaw/
├── cmd/                         # Go binary entry points
│   ├── picoclaw/                # Main PicoClaw CLI binary
│   │   ├── main.go              # Cobra root command, banner, main()
│   │   ├── dns_noresolv.go      # DNS resolution bypass
│   │   └── internal/            # CLI subcommand implementations
│   │       ├── agent/           # "picoclaw agent" subcommand
│   │       ├── auth/            # "picoclaw auth" subcommand
│   │       ├── cliui/           # CLI rendering utilities
│   │       ├── cron/            # "picoclaw cron" subcommand
│   │       ├── gateway/         # "picoclaw gateway" subcommand
│   │       ├── mcp/             # "picoclaw mcp" subcommand
│   │       ├── migrate/         # "picoclaw migrate" subcommand
│   │       ├── model/           # "picoclaw model" subcommand
│   │       ├── onboard/         # "picoclaw onboard" subcommand
│   │       ├── skills/          # "picoclaw skills" subcommand
│   │       ├── status/          # "picoclaw status" subcommand
│   │       └── version/         # "picoclaw version" subcommand
│   └── membench/                # Memory benchmark utility
│
├── pkg/                         # Shared library packages
│   ├── agent/                   # Core AI agent processing engine
│   │   ├── agent.go             # AgentLoop struct, main turn processing
│   │   ├── pipeline*.go         # Pipeline: setup, LLM, finalize, execute phases
│   │   ├── context*.go          # ContextManager, budget, seahorse, legacy
│   │   ├── hooks*.go            # HookManager, mount, process
│   │   ├── turn_*.go            # TurnCoordinator, turn_state, subturn
│   │   ├── prompt*.go           # Prompt assembly, contributors
│   │   ├── steering*.go         # Steering queue for agent control
│   │   ├── registry*.go         # Agent registry and discovery
│   │   ├── instance*.go         # Agent instance management
│   │   ├── agent_mcp.go         # MCP integration
│   │   ├── agent_message.go     # Message processing entry
│   │   ├── agent_outbound.go    # Outbound message handling
│   │   ├── agent_media.go       # Media handling
│   │   ├── agent_transcribe.go  # Audio transcription
│   │   ├── agent_stop.go        # Agent stop logic
│   │   ├── agent_command.go     # Slash command processing
│   │   ├── llm_media.go         # LLM media utilities
│   │   ├── memory.go            # Memory management
│   │   ├── events*.go           # Agent-level event bus
│   │   ├── event_payloads.go    # Event type definitions
│   │   ├── dispatch_request.go  # Request dispatching
│   │   ├── model_resolution.go  # Model selection logic
│   │   ├── thinking.go          # Extended thinking support
│   │   ├── discovery.go         # Agent capability discovery
│   │   ├── adapters/            # Adapters bridging channel/agent interfaces
│   │   │   ├── channelmanager.go
│   │   │   └── messagebus.go
│   │   └── interfaces/          # Agent-level interface definitions
│   │       └── interfaces.go
│   ├── audio/                   # Audio processing
│   │   ├── asr/                 # Automatic Speech Recognition
│   │   ├── tts/                 # Text-to-Speech
│   │   └── ogg.go               # OGG format handling
│   ├── auth/                    # Authentication (OAuth, PKCE, token mgmt)
│   ├── bus/                     # Internal message bus (inbound/outbound)
│   ├── channels/                # Multi-platform messaging channel adapters
│   │   ├── manager.go           # Channel lifecycle manager (1664 lines)
│   │   ├── base.go              # Base channel implementation
│   │   ├── registry.go          # Channel registry
│   │   ├── interfaces.go        # Channel capability interfaces
│   │   ├── dingtalk/            # DingTalk channel
│   │   ├── discord/             # Discord channel
│   │   ├── feishu/              # Feishu (Lark) channel
│   │   ├── irc/                 # IRC channel
│   │   ├── line/                # LINE channel
│   │   ├── maixcam/             # MaixCam channel
│   │   ├── matrix/              # Matrix channel
│   │   ├── mqtt/                # MQTT channel
│   │   ├── onebot/              # OneBot channel
│   │   ├── pico/                # Pico internal WebSocket channel
│   │   ├── qq/                  # QQ channel
│   │   ├── slack/               # Slack channel
│   │   ├── teams_webhook/       # Microsoft Teams webhook
│   │   ├── telegram/            # Telegram channel
│   │   ├── vk/                  # VKontakte channel
│   │   ├── wecom/               # WeCom (Enterprise WeChat) channel
│   │   ├── weixin/              # WeChat Official Account channel
│   │   ├── whatsapp/            # WhatsApp (mautrix) channel
│   │   └── whatsapp_native/     # WhatsApp native (whatsmeow) channel
│   ├── commands/                # Built-in slash commands (/help, /clear, etc.)
│   ├── config/                  # Configuration management
│   │   ├── config.go            # Config struct, parsing, validation (1596 lines)
│   │   ├── migration.go         # Config schema migration
│   │   ├── security.go          # Encrypted credential storage
│   │   └── version.go           # Build version info
│   ├── constants/               # Shared constants
│   ├── credential/              # Credential encryption and keygen
│   ├── cron/                    # Cron scheduling service
│   ├── devices/                 # Device management
│   │   ├── service.go           # Device service
│   │   ├── sources/             # Device source adapters
│   │   └── events/              # Device events
│   ├── events/                  # Runtime event bus (pub/sub with filtering)
│   ├── fileutil/                # File system utilities
│   ├── gateway/                 # Gateway orchestrator
│   │   ├── gateway.go           # Main gateway runtime (829 lines)
│   │   ├── listen.go            # Event loop
│   │   ├── events.go            # Gateway event types
│   │   └── channel_matrix.go    # Matrix channel integration
│   ├── health/                  # Health check HTTP server
│   ├── heartbeat/               # Heartbeat service
│   ├── identity/                # Agent identity management
│   ├── isolation/               # Process/filesystem isolation (Linux/Windows)
│   ├── logger/                  # Logging (zerolog wrapper)
│   ├── mcp/                     # Model Context Protocol integration
│   ├── media/                   # Media storage (filesystem-backed)
│   ├── memory/                  # JSONL-based memory store
│   ├── migrate/                 # Data migration framework
│   ├── netbind/                 # Network address binding utilities
│   ├── pid/                     # PID file management
│   ├── providers/               # LLM provider adapters
│   │   ├── factory.go           # Provider factory
│   │   ├── fallback.go          # Multi-provider fallback chain
│   │   ├── cooldown.go          # Provider cooldown management
│   │   ├── legacy_provider.go   # Backward-compatible provider interface
│   │   ├── anthropic/           # Anthropic Claude provider
│   │   ├── anthropic_messages/  # Anthropic Messages API
│   │   ├── openai_compat/       # OpenAI-compatible providers
│   │   ├── azure/               # Azure OpenAI provider
│   │   ├── bedrock/             # AWS Bedrock provider
│   │   ├── cli/                 # CLI-based provider
│   │   ├── httpapi/             # HTTP API provider
│   │   ├── oauth/               # OAuth-based provider
│   │   └── common/              # Shared provider utilities
│   ├── routing/                 # Agent routing and message classification
│   ├── seahorse/                # Long-term context engine (SQLite FTS5)
│   ├── session/                 # Session management
│   ├── skills/                  # Installable skill modules (ClawHub/GitHub)
│   ├── state/                   # Runtime state persistence
│   ├── tokenizer/               # Token estimation
│   ├── tools/                   # Agent tools (shell, fs, subagents, cron, search)
│   │   ├── registry.go          # Tool registry
│   │   ├── shell.go             # Shell execution
│   │   ├── subagent.go          # Sub-agent spawning
│   │   ├── spawn.go             # Process spawning
│   │   ├── cron.go              # Cron integration
│   │   ├── search_tool.go       # Search capability
│   │   ├── fs/                  # File system tools
│   │   ├── hardware/            # Hardware tools
│   │   ├── integration/         # Integration tools
│   │   └── shared/              # Shared tool utilities
│   ├── updater/                 # Self-update mechanism
│   ├── utils/                   # General utilities
│   └── env.go                   # Environment variable reading
│
├── web/                         # Web Console (separate binary)
│   ├── backend/                 # Go HTTP API server
│   │   ├── main.go              # Web server entry point (699 lines)
│   │   ├── app_runtime.go       # Shutdown orchestration
│   │   ├── embed.go             # Embedded frontend dist
│   │   ├── systray.go           # System tray integration
│   │   ├── api/                 # REST API handlers
│   │   │   ├── router.go        # Route registration
│   │   │   ├── gateway.go       # Gateway control APIs
│   │   │   ├── config.go        # Config CRUD APIs
│   │   │   ├── models.go        # Model management APIs
│   │   │   ├── pico.go          # Pico Channel WebSocket proxy
│   │   │   ├── auth.go          # Authentication APIs
│   │   │   ├── oauth.go         # OAuth flow APIs
│   │   │   ├── skills.go        # Skill management APIs
│   │   │   ├── tools.go         # Tool management APIs
│   │   │   ├── channels.go      # Channel management APIs
│   │   │   ├── session.go       # Session APIs
│   │   │   └── version.go       # Version info APIs
│   │   ├── dashboardauth/       # Dashboard authentication
│   │   ├── launcherconfig/      # Launcher configuration
│   │   ├── middleware/          # HTTP middleware
│   │   ├── model/               # API data models
│   │   ├── utils/               # Backend utilities
│   │   └── dist/                # Built frontend output (committed)
│   ├── frontend/                # React SPA (TypeScript)
│   │   ├── src/
│   │   │   ├── main.tsx         # App entry point
│   │   │   ├── app-providers.tsx# QueryClient, Router, i18n providers
│   │   │   ├── routeTree.gen.ts # Auto-generated route tree (TanStack Router)
│   │   │   ├── index.css        # Global styles (Tailwind)
│   │   │   ├── api/             # API client modules
│   │   │   │   ├── http.ts      # HTTP client with auth
│   │   │   │   ├── gateway.ts   # Gateway API client
│   │   │   │   ├── pico.ts      # Pico Channel/WebSocket client
│   │   │   │   ├── channels.ts  # Channel API client
│   │   │   │   ├── models.ts    # Models API client
│   │   │   │   ├── skills.ts    # Skills API client
│   │   │   │   ├── tools.ts     # Tools API client
│   │   │   │   └── ...
│   │   │   ├── components/      # React components
│   │   │   │   ├── ui/          # shadcn/ui primitives (button, dialog, etc.)
│   │   │   │   ├── chat/        # Chat UI components
│   │   │   │   ├── agent/       # Agent management UI (hub, skills, tools)
│   │   │   │   ├── config/      # Configuration UI
│   │   │   │   ├── models/      # Model management UI
│   │   │   │   ├── credentials/ # Credential management UI
│   │   │   │   ├── channels/    # Channel configuration UI
│   │   │   │   ├── logs/        # Log viewer UI
│   │   │   │   ├── tour/        # Onboarding tour
│   │   │   │   ├── app-layout.tsx
│   │   │   │   ├── app-sidebar.tsx
│   │   │   │   └── app-header.tsx
│   │   │   ├── features/        # Domain-specific feature modules
│   │   │   │   └── chat/        # Chat feature: state machine, protocol, WebSocket
│   │   │   ├── hooks/           # React custom hooks
│   │   │   ├── i18n/            # Internationalization (i18next)
│   │   │   ├── lib/             # Utility functions
│   │   │   ├── routes/          # TanStack Router file-based routes
│   │   │   │   ├── __root.tsx   # Root layout with sidebar
│   │   │   │   ├── index.tsx    # Home / redirect
│   │   │   │   ├── chat.tsx     # Chat route
│   │   │   │   ├── logs.tsx     # Logs route
│   │   │   │   ├── config.tsx   # Config route
│   │   │   │   ├── models.tsx   # Models route
│   │   │   │   ├── credentials.tsx # Credentials route
│   │   │   │   ├── agent.tsx    # Agent layout route
│   │   │   │   ├── agent/       # Agent sub-routes
│   │   │   │   │   ├── hub.tsx      # Skill marketplace
│   │   │   │   │   ├── skills.tsx   # Installed skills
│   │   │   │   │   └── tools.tsx    # Tool configuration
│   │   │   │   └── channels/    # Channel sub-routes
│   │   │   │       └── $name.tsx    # Per-channel config
│   │   │   └── store/           # Jotai atoms for global state
│   │   │       ├── index.ts
│   │   │       ├── chat.ts
│   │   │       ├── gateway.ts
│   │   │       └── tour.ts
│   │   ├── package.json         # Node dependencies (React, TanStack, shadcn, etc.)
│   │   ├── vite.config.ts       # Vite bundler config (Tailwind, TanStack Router plugin)
│   │   ├── tsconfig.json        # TypeScript configuration
│   │   ├── eslint.config.js     # ESLint configuration
│   │   └── prettier.config.js   # Prettier configuration
│   └── Makefile                 # Web build targets
│
├── docs/                        # Documentation (multi-language)
│   ├── architecture/            # Architecture docs
│   ├── channels/                # Per-channel setup guides
│   ├── guides/                  # User guides
│   ├── migration/               # Migration guides
│   ├── operations/              # Operations docs
│   ├── project/                 # Project READMEs (multi-language)
│   ├── reference/               # Reference documentation
│   └── security/                # Security documentation
│
├── docker/                      # Dockerfiles and compose
│   ├── Dockerfile               # Standard build
│   ├── Dockerfile.full          # Full build with extra features
│   ├── Dockerfile.heavy         # Heavy build
│   ├── Dockerfile.launcher      # Web launcher build
│   ├── Dockerfile.goreleaser    # GoReleaser build
│   └── docker-compose.yml       # Docker Compose configuration
│
├── config/                      # Example/default config
│   └── config.example.json
│
├── scripts/                     # Build/utility scripts
├── assets/                      # Static assets
├── build/                       # Build output directory
├── examples/                    # Example projects (pico-echo-server)
├── workspace/                   # Agent workspace directory
│
├── go.mod                       # Go module definition (go 1.25.10)
├── go.sum                       # Go dependency checksums
├── Makefile                     # Root build system (511 lines)
├── .golangci.yaml               # Go linter configuration
├── .goreleaser.yaml             # GoReleaser configuration
├── .gitignore
├── .gitattributes
├── LICENSE                      # MIT License
├── README.md                    # Project README
├── ROADMAP.md                   # Development roadmap
└── CONTRIBUTING.md              # Contribution guidelines
```

## Directory Purposes

**`cmd/picoclaw/`:**
- Purpose: Main executable entry point; Cobra CLI command tree
- Contains: `main.go` (root command, banner, subcommand registration), `internal/` (one subdirectory per subcommand)
- Key files: `cmd/picoclaw/main.go`, `cmd/picoclaw/internal/gateway/`, `cmd/picoclaw/internal/agent/`

**`pkg/`:**
- Purpose: All shared library code; the core of the application
- Contains: ~36 package directories covering agent logic, channels, providers, tools, storage, utilities
- Key files: `pkg/agent/agent.go`, `pkg/gateway/gateway.go`, `pkg/channels/manager.go`, `pkg/providers/factory.go`

**`pkg/agent/`:**
- Purpose: Core AI agent turn processing, context management, prompt assembly, hook system
- Contains: 79 files — the most complex package in the codebase
- Key files: `agent.go` (AgentLoop), `pipeline.go` (Pipeline), `context_manager.go`, `hooks.go`, `turn_coord.go`

**`pkg/channels/`:**
- Purpose: Multi-platform messaging channel adapters with 17+ platform implementations
- Contains: 46 entries — manager, base, interfaces, events, and one subdirectory per platform
- Key files: `manager.go` (1664 lines), `interfaces.go`, `base.go`, `registry.go`

**`pkg/providers/`:**
- Purpose: LLM provider abstraction, factory, fallback chains, rate limiting
- Contains: 38 entries — factory, fallback, 6+ provider backends, common utilities
- Key files: `factory.go`, `fallback.go`, `cooldown.go`, `provider_catalog.go`

**`web/backend/`:**
- Purpose: Separate Go binary providing HTTP API and web console server; manages the main gateway process
- Contains: 21 entries — main server, API handlers (router, auth, config, models, pico proxy, etc.), middleware
- Key files: `main.go` (699 lines), `api/router.go`, `api/gateway.go`, `api/pico.go`

**`web/frontend/`:**
- Purpose: React SPA for web-based chat, configuration, and management
- Contains: TypeScript/React sources organized by feature; ~100+ component files
- Key files: `src/main.tsx`, `src/routeTree.gen.ts`, `src/routes/__root.tsx`, `src/features/chat/state.ts`

**`docs/`:**
- Purpose: Multi-language documentation organized by topic
- Contains: Architecture docs, channel setup guides, migration guides, multi-language READMEs
- Key files: `docs/architecture/README.md`, channel-specific READMEs under `docs/channels/`

**`docker/`:**
- Purpose: Container build and deployment configuration
- Contains: Multiple Dockerfiles for different build profiles, docker-compose configurations
- Key files: `Dockerfile`, `docker-compose.yml`, `Dockerfile.full`

## Key File Locations

**Entry Points:**
- `cmd/picoclaw/main.go`: CLI entry — Cobra command tree, banner display, main()
- `pkg/gateway/gateway.go`: Gateway runtime — all service orchestration
- `web/backend/main.go`: Web console server entry point
- `web/frontend/src/main.tsx`: React app mount point

**Configuration:**
- `config/config.example.json`: Reference configuration file
- `pkg/config/config.go`: Config struct definition, loading, validation (1596 lines)
- `pkg/config/migration.go`: Schema version migration (v1→v2→v3)
- `pkg/config/security.go`: Encrypted credential handling
- `pkg/config/envkeys.go`: Environment variable key definitions

**Core Logic:**
- `pkg/agent/agent.go`: AgentLoop — main AI processing loop (640 lines)
- `pkg/agent/pipeline.go`: Pipeline — structured turn execution phases
- `pkg/agent/context_manager.go`: Token budget and context window management
- `pkg/channels/manager.go`: Channel lifecycle, routing, rate limiting (1664 lines)
- `pkg/providers/factory.go`: LLM provider instantiation
- `pkg/providers/fallback.go`: Multi-provider fallback chain logic

**Testing:**
- Tests co-located with source files using `_test.go` suffix (Go convention)
- Example: `pkg/agent/agent_test.go`, `pkg/channels/manager_test.go`
- Web frontend: Jest/Vitest patterns (configured in `vite.config.ts`)

**Build:**
- `Makefile`: 511-line root build system with cross-platform support
- `web/Makefile`: Web console build targets
- `.goreleaser.yaml`: Automated release configuration
- `.golangci.yaml`: Go linter rules

## Naming Conventions

**Files:**
- Go: `snake_case.go` for most packages (e.g., `context_manager.go`, `agent_init.go`); `_test.go` suffix for tests
- TypeScript: `kebab-case.tsx` for components and utilities (e.g., `chat-page.tsx`, `use-chat-models.ts`); `PascalCase.tsx` for route files matching TanStack Router convention
- Platform-specific: `*_unix.go` / `*_windows.go` suffixes; `*_nocgo.go` for non-CGO stubs
- Sub-package init: `init.go` files for package-level initialization (e.g., `pkg/channels/pico/init.go`)

**Directories:**
- Go packages: `lowercase` single-word names (e.g., `agent`, `channels`, `providers`); sub-packages use nested directories (e.g., `channels/telegram/`)
- TypeScript: `lowercase` or `kebab-case` (e.g., `components/ui/`, `features/chat/`, `hooks/`)
- Route directories: Parameterized routes use `$paramName` pattern (e.g., `routes/channels/$name.tsx`)

**Go Exports:**
- Exported types/functions: PascalCase (e.g., `AgentLoop`, `Pipeline`, `StartGateway`)
- Unexported: camelCase (e.g., `runTurn`, `processInbound`)
- Interface naming: descriptive nouns (e.g., `MessageBus`, `ChannelManager`, `TypingCapable`)

**TypeScript:**
- Components: PascalCase files with default exports (e.g., `ChatPage.tsx`)
- Hooks: `use-*.ts` with named exports (e.g., `useChatModels`, `useGateway`)
- API modules: named exports grouping related functions
- State atoms: named exports from `store/` (e.g., `chatAtom`, `gatewayAtom`)

## Where to Add New Code

**New Go Feature (e.g., new tool capability):**
- Primary code: `pkg/tools/newfeature.go`
- Registration: Add to `pkg/tools/registry.go`
- Tests: `pkg/tools/newfeature_test.go`
- CLI exposure (if needed): `cmd/picoclaw/internal/newfeature/`

**New Channel Adapter:**
- Implementation: `pkg/channels/newplatform/newplatform.go`
- Registration: `pkg/channels/newplatform/init.go` with `init()` side-effect register
- Import registration: Add blank import to `pkg/gateway/gateway.go`
- Docs: `docs/channels/newplatform/README.md`

**New LLM Provider:**
- Implementation: `pkg/providers/newprovider/newprovider.go`
- Factory integration: Add case to `pkg/providers/factory.go`
- Common utilities: `pkg/providers/common/` if shared
- Tests: `pkg/providers/newprovider/newprovider_test.go`

**New Web Frontend Route:**
- Route definition: `web/frontend/src/routes/newpage.tsx`
- Components: `web/frontend/src/components/newpage/newpage-page.tsx`
- API client (if needed): `web/frontend/src/api/newfeature.ts`
- Store (if global state): `web/frontend/src/store/newfeature.ts`

**New Utilities:**
- Shared helpers: `pkg/utils/` for small utilities
- Shared frontend lib: `web/frontend/src/lib/` for frontend utilities

## Special Directories

**`build/`:**
- Purpose: Compiled binary output
- Generated: Yes (by `make build`)
- Committed: No (in `.gitignore`)

**`web/backend/dist/`:**
- Purpose: Built frontend assets embedded by Go backend
- Generated: Yes (by `npm run build:backend`)
- Committed: Yes (embedded for standalone distribution)

**`web/frontend/src/routeTree.gen.ts`:**
- Purpose: Auto-generated route tree from TanStack Router file-based routing
- Generated: Yes (by Vite plugin)
- Committed: Yes (checked in for build reproducibility)

**`web/frontend/node_modules/`:**
- Purpose: Node.js dependencies
- Generated: Yes (by `pnpm install`)
- Committed: No (in `.gitignore`)

**`workspace/`:**
- Purpose: Agent runtime workspace for file operations
- Generated: Partially (created at runtime)
- Committed: No (in `.gitignore`)

**`.cache/`:**
- Purpose: Go build cache (`go-build`, `go-mod`)
- Generated: Yes
- Committed: No (in `.gitignore`)

**`logs/`:**
- Purpose: Runtime log output (`gateway.log`, `gateway_panic.log`, `launcher.log`)
- Generated: Yes
- Committed: No (in `.gitignore`)

---

*Structure analysis: 2026-05-11*
