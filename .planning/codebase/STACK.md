# Technology Stack

**Analysis Date:** 2026-05-11

## Languages

**Primary:**
- Go 1.25.10 - Core runtime, CLI, gateway, agent engine, all backend logic (`go.mod`, `cmd/`, `pkg/`, `web/backend/`)

**Secondary:**
- TypeScript ~5.9.3 - Web UI frontend (`web/frontend/`, `package.json`)
- Shell (POSIX sh) - Docker entrypoint, scripts (`docker/entrypoint.sh`, `scripts/`)

## Runtime

**Environment:**
- Go 1.25+ (toolchain locked to `local` in `Makefile`)
- Multi-platform: Linux, macOS, Darwin, Windows, FreeBSD, NetBSD, Android (GOOS/GOARCH)

**Package Manager:**
- Go modules (`go.mod`, `go.sum`)
- Lockfile: present (`go.sum`)
- pnpm 10.33.0 (Web UI frontend, `web/frontend/package.json`)
- pnpm lockfile: present (`web/frontend/pnpm-lock.yaml`)

## Frameworks

**Core:**
- `github.com/spf13/cobra` v1.10.2 - CLI framework (`cmd/picoclaw/main.go:33`)
- `github.com/spf13/pflag` v1.0.10 - CLI flag parsing

**Frontend:**
- React 19.2.5 - UI library (`web/frontend/package.json`)
- Vite 8.0.10 - Build tool (`web/frontend/vite.config.ts`)
- Tailwind CSS 4.2.4 - CSS framework
- Radix UI (via `shadcn`) - Headless component primitives (`web/frontend/components.json`)
- TanStack Router 1.169+ - Client-side routing
- TanStack React Query 5.99+ - Server state management

**Testing:**
- `github.com/stretchr/testify` v1.11.1 - Go test assertions
- Go standard `testing` package - Test runner

**Build/Dev:**
- `golangci-lint` v2 - Linting and formatting (`.golangci.yaml`)
- `goreleaser` v2 - Multi-platform release automation (`.goreleaser.yaml`)
- GNUMake - Build orchestration (`Makefile`)
- `eslint` 10.2+ - TypeScript linting
- `prettier` 3.8+ - Code formatting
- `typescript` 5.9.3 - TypeScript compiler
- `go:embed` - Static asset embedding (`web/backend/embed.go:16`)

## Key Dependencies

**Critical (Direct):**
- `modernc.org/sqlite` v1.48.2 - Embedded SQLite database (pure Go, no CGo) — used for session state, password store, channel state
- `github.com/rs/zerolog` v1.35.1 - Structured JSON logging (`pkg/logger/`)
- `github.com/google/uuid` v1.6.0 - UUID generation
- `github.com/caarlos0/env/v11` v11.4.0 - Environment variable parsing (`pkg/config/envkeys.go`)
- `gopkg.in/yaml.v3` v3.0.1 - YAML parsing for `.security.yml` (`pkg/config/security.go`)
- `github.com/gorilla/websocket` v1.5.3 - WebSocket client/server for chat channels
- `github.com/creack/pty` v1.1.24 - Pseudoterminal for subprocess exec (`pkg/tools/shell*.go`)
- `github.com/charmbracelet/lipgloss` v1.1.0 - Terminal UI styling (`cmd/picoclaw/internal/cliui/`)
- `github.com/muesli/termenv` v0.16.0 - Terminal color support
- `github.com/gomarkdown/markdown` - Markdown rendering

**Infrastructure:**
- `github.com/minio/selfupdate` v0.6.0 - Self-updating binary (`pkg/updater/`)
- `golang.org/x/oauth2` v0.36.0 - OAuth2 client (`pkg/auth/`)
- `golang.org/x/term` v0.42.0 - Terminal handling
- `golang.org/x/crypto` v0.50.0 - Cryptographic primitives
- `golang.org/x/sync` v0.20.0 - Concurrency primitives

**LLM/AI SDKs:**
- `github.com/openai/openai-go/v3` v3.22.0 - OpenAI API client (`pkg/providers/openai_compat/`)
- `github.com/anthropics/anthropic-sdk-go` v1.26.0 - Anthropic API client
- `github.com/aws/aws-sdk-go-v2/service/bedrockruntime` v1.50.6 - AWS Bedrock (conditional build tag `bedrock`)
- `github.com/google/jsonschema-go` v0.4.3 - JSON Schema for tool definitions
- `github.com/github/copilot-sdk/go` v0.2.0 - GitHub Copilot integration and dev tooling
- `github.com/modelcontextprotocol/go-sdk` v1.5.0 - MCP protocol implementation (`pkg/mcp/`)

**Chat Platform SDKs (19+ channels):**
- `github.com/mymmrac/telego` v1.8.0 - Telegram
- `github.com/bwmarrin/discordgo` (forked) - Discord
- `github.com/slack-go/slack` v0.17.3 - Slack
- `go.mau.fi/whatsmeow` - WhatsApp native
- `github.com/larksuite/oapi-sdk-go/v3` v3.6.1 - Feishu/Lark
- `github.com/open-dingtalk/dingtalk-stream-sdk-go` v0.9.1 - DingTalk
- `github.com/line/line-bot-sdk-go/v8` v8.19.0 - LINE
- `github.com/SevereCloud/vksdk/v3` v3.3.1 - VK
- `github.com/ergochat/irc-go` v0.6.0 - IRC
- `github.com/tencent-connect/botgo` v0.2.1 - QQ
- `maunium.net/go/mautrix` v0.27.0 - Matrix
- `github.com/pion/webrtc/v3` v3.3.6 - WebRTC

**Desktop UI:**
- `fyne.io/systray` v1.12.1 - System tray icon (`web/backend/systray.go`)

## Configuration

**Environment:**
- Config file: `~/.picoclaw/config.json` (JSON format) — overridable via `PICOCLAW_CONFIG` env var
- Secrets file: `~/.picoclaw/.security.yml` (YAML format) — sensitive data separated from config
- Home directory: `~/.picoclaw/` — overridable via `PICOCLAW_HOME` env var
- Key env vars: `PICOCLAW_HOME`, `PICOCLAW_CONFIG`, `PICOCLAW_BINARY`, `PICOCLAW_GATEWAY_HOST`, `PICOCLAW_BUILTIN_SKILLS`, `TZ`
- Provider API keys: loaded from `.security.yml` or env vars (e.g., `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, `TELEGRAM_BOT_TOKEN`)

**Build:**
- `Makefile` — primary build with LDFLAGS for version injection (`pkg/config/version.go`)
- `.goreleaser.yaml` — multi-platform release (Linux, macOS, Windows, FreeBSD, NetBSD; amd64/arm64/riscv64/loong64/arm/mipsle/s390x)
- Build tags: `goolm`, `stdjson` (default), `bedrock` (AWS Bedrock), `whatsapp_native` (native WhatsApp), `integration` (integration tests)
- CGo: disabled by default (`CGO_ENABLED=0`) — pure Go binaries

## Platform Requirements

**Development:**
- Go 1.25+
- Node.js 22+ and pnpm 10.33.0+ (for Web UI / launcher builds)
- `make`, `git`
- `golangci-lint` (for linting)

**Production:**
- Single static binary (pure Go, no CGo)
- Deploy on any Linux kernel (x86_64, ARM64, ARMv7, RISC-V, LoongArch, MIPS)
- macOS (arm64), Windows (amd64), FreeBSD (amd64), NetBSD (amd64/arm64)
- Android via APK package or Termux
- Docker: multi-stage Alpine-based images (<50MB)

---

*Stack analysis: 2026-05-11*
