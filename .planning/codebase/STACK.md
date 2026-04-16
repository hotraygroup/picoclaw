# Technology Stack

**Analysis Date:** 2026-04-16

## Languages

**Primary:**
- Go 1.25.9 - Core backend, agent loop, providers, tools, gateway server
- TypeScript 5.9 - Web frontend (React SPA)

**Secondary:**
- YAML - Security configuration (`.security.yml`)
- JSON - Main configuration files

## Runtime

**Environment:**
- Go runtime (statically compiled, CGO_ENABLED=0 for most builds)
- Node.js 20.19+ or 22.13+ (frontend development)

**Package Manager:**
- Go modules (go.mod/go.sum) - Backend dependencies
- pnpm 10.33.0 - Frontend dependencies (web frontend monorepo)

## Frameworks

**Core Backend:**
- Cobra 1.10.2 - CLI framework for command structure
- Gorilla WebSocket 1.5.3 - WebSocket server (gateway, pico channel)
- rs/zerolog 1.35.0 - Structured logging
- spf13/pflag 1.0.10 - Flag parsing (CLI)

**Web Frontend:**
- React 19.2.5 - UI framework
- Vite 8.0.8 - Build tool, dev server
- TanStack Router 1.167.0 - Client-side routing
- TanStack React Query 5.97.0 - Data fetching, caching
- TailwindCSS 4.2.2 - Styling
- shadcn/ui 4.2.0 - UI component library
- Jotai 2.19.1 - State management (atomic state)

**Testing:**
- stretchr/testify 1.11.1 - Go assertion library (unit tests)
- Vitest - Frontend testing (via Vite)

**Build/Dev:**
- golangci-lint - Go linting/formatting
- goreleaser - Release automation
- Make - Build orchestration

## Key Dependencies

**LLM Provider SDKs:**
- github.com/anthropics/anthropic-sdk-go v1.26.0 - Anthropic API client
- github.com/openai/openai-go/v3 v3.22.0 - OpenAI API client
- github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.50.4 - AWS Bedrock
- github.com/google/generative-ai (via antigravity) - Gemini
- github.com/modelcontextprotocol/go-sdk v1.5.0 - MCP protocol

**Messaging Platform SDKs:**
- github.com/mymmrac/telego v1.8.0 - Telegram Bot API
- github.com/bwmarrin/discordgo v0.29.0 - Discord Bot API
- go.mau.fi/whatsmeow - WhatsApp native protocol
- github.com/slack-go/slack v0.17.3 - Slack API
- maunium.net/go/mautrix v0.26.4 - Matrix protocol
- github.com/larksuite/oapi-sdk-go/v3 v3.5.3 - Feishu/Lark
- github.com/open-dingtalk/dingtalk-stream-sdk-go v0.9.1 - DingTalk
- github.com/SevereCloud/vksdk/v3 v3.3.1 - VK (Antigravity)
- github.com/ergochat/irc-go v0.6.0 - IRC protocol
- github.com/tencent-connect/botgo v0.2.1 - QQ Bot API

**Infrastructure:**
- modernc.org/sqlite v1.48.2 - Pure Go SQLite (session storage, auth)
- github.com/gorilla/websocket v1.5.3 - WebSocket handling
- github.com/minio/selfupdate v0.6.0 - Self-update mechanism
- fyne.io/systray v1.12.0 - System tray integration (launcher)
- github.com/adhocore/gronx v1.19.6 - Cron expression parsing (scheduled tasks)

**UI/Terminal:**
- github.com/rivo/tview v0.42.0 - TUI components
- github.com/gdamore/tcell/v2 v2.13.8 - Terminal cell handling
- github.com/charmbracelet/lipgloss v1.1.0 - Terminal styling
- github.com/ergochat/readline v0.1.3 - Readline support

**Media/Processing:**
- github.com/pion/webrtc/v3 v3.3.6 - WebRTC for voice
- github.com/pion/rtp v1.10.1 - RTP protocol
- github.com/h2non/filetype v1.1.3 - File type detection
- github.com/gomarkdown/markdown - Markdown rendering

## Build Tags

Default tags: `goolm,stdjson`

Special build tags:
- `whatsapp_native` - WhatsApp native support (larger binary)
- `bedrock` - AWS Bedrock provider support

## Configuration

**Environment:**
- Main config: `~/.picoclaw/config.json` (or `PICOCLAW_CONFIG` env)
- Security config: `~/.picoclaw/.security.yml` (API keys, tokens)
- Launcher config: `~/.picoclaw/launcher-config.json`

**Key Environment Variables:**
- `PICOCLAW_HOME` - Base data directory
- `PICOCLAW_CONFIG` - Config file path override
- `PICOCLAW_BINARY` - Binary path (for launcher subprocess)
- `PICOCLAW_GATEWAY_HOST` - Gateway host override
- `PICOCLAW_LAUNCHER_HOST` - Launcher bind host override
- `PICOCLAW_LAUNCHER_TOKEN` - Stable auth token

## Platform Requirements

**Development:**
- Go 1.25+
- Node.js 20.19+ / 22.13+ (for web frontend)
- pnpm 10.33.0+
- golangci-lint

**Production:**
- Cross-compiled for: Linux ARM, Linux ARM64, Linux MIPS32 LE, Android ARM64, Windows, macOS
- Binary sizes: <10MB (optimized for $10 hardware)
- Supports Raspberry Pi Zero 2 W, MIPS devices (Ingenic X2600)

---

*Stack analysis: 2026-04-16*