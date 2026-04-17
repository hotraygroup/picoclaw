# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PicoClaw is an ultra-lightweight personal AI assistant written in Go. It runs on $10 hardware with <10MB RAM and supports 30+ LLM providers. The project was substantially developed with AI assistance and embraces AI-assisted contributions.

## Build Commands

```bash
make build              # Build core binary for current platform (runs go generate first)
make build-launcher     # Build WebUI launcher binary
make build-launcher-tui # Build TUI launcher binary
make build-all          # Build core binaries for all platforms
make install            # Install to ~/.local/bin

# Platform-specific builds
make build-linux-arm    # Raspberry Pi Zero 2 W 32-bit
make build-linux-arm64  # Raspberry Pi Zero 2 W 64-bit
make build-linux-mipsle # MIPS32 LE (Ingenic X2600)
make build-android-arm64

# Android bundle
make build-android-bundle  # Creates universal zip with JNI libs
```

## Testing & Linting

```bash
make test               # Run all tests (Go + web)
make check              # Full pre-commit: deps + fmt + vet + test
make fmt                # Format code (golangci-lint fmt)
make lint               # Run linters
make fix                # Auto-fix linting issues
make vet                # Static analysis

# Run single test
go test -run TestName -v ./pkg/session/

# Run benchmarks
go test -bench=. -benchmem -run='^$' ./...
```

## Web UI Development

```bash
cd web
make dev                # Full launcher dev: builds picoclaw, starts Go backend + Vite frontend
make dev-frontend       # Vite dev server only
make dev-backend        # Go backend only
make build              # Build standalone launcher binary
make build-frontend     # Build embeddable frontend into backend/dist
make test               # Backend Go tests + frontend lint
make lint               # go vet + pnpm check
```

Frontend prerequisites: Node.js 22+ and pnpm 10.33.0+

## Architecture

### Core Packages

- **pkg/agent/**: Agent loop (`loop.go`), context management, hooks system, steering, subturn coordination
- **pkg/providers/**: LLM provider implementations (OpenAI, Anthropic, Gemini, Azure, Bedrock, Ollama, etc.)
- **pkg/channels/**: Messaging channels (pico only)
- **pkg/tools/**: Built-in tools (filesystem, shell, web search, cron, MCP, skills, spawn/subagent)
- **pkg/config/**: Configuration system with security.yml for sensitive data separation
- **pkg/gateway/**: Gateway server that orchestrates channels, agents, and services
- **pkg/session/**: Session management and allocation
- **pkg/memory/**: JSONL-based memory storage

### Provider Format

Providers use `protocol/model` format in config:
- `openai/gpt-5.4`, `anthropic/claude-opus-4-6`, `gemini/gemini-3-flash`
- `ollama/llama3.1:8b` (local, no API key needed)
- `azure/gpt-deployment-name`, `bedrock/claude-3-opus`

### Tools System

Tools in `pkg/tools/` include:
- `filesystem.go`: File operations with sandbox
- `shell.go`: Command execution with timeout/policy controls
- `web.go`: Web search (DuckDuckGo, Tavily, Brave, SearXNG, Baidu, Perplexity)
- `cron.go`: Scheduled tasks/reminders
- `mcp_tool.go`: Model Context Protocol integration
- `skills_install.go`, `skills_search.go`: Skills registry management
- `spawn.go`, `subagent.go`: Async task spawning and sub-agent coordination

### Configuration

- Main config: `~/.picoclaw/config.json` (or `PICOCLAW_CONFIG` env)
- Sensitive data: `~/.picoclaw/.security.yml` (API keys, tokens)
- Launcher config: `~/.picoclaw/launcher-config.json`

Config versioning: v0 (legacy, sensitive in config) → v1+ (sensitive in .security.yml). Migration handled automatically.

### Web Launcher Architecture

`web/` is a monorepo:
- `backend/`: Go HTTP server with REST APIs, auth, WebSocket proxy to gateway
- `frontend/`: Vite + React 19 + TanStack Router SPA

The launcher manages `picoclaw gateway -E` as a subprocess. Chat traffic proxies through `/pico/ws`.

## Key Patterns

### Adding a New Provider

1. Create package in `pkg/providers/`
2. Implement `LLMProvider` interface (defined in `pkg/providers/types.go`)
3. Add protocol metadata in `pkg/providers/factory_provider.go` (if using OpenAI-compatible format)
4. Update provider creation logic if non-standard protocol

### Adding a New Tool

1. Add tool definition in `pkg/tools/`
2. Register in `pkg/tools/registry.go`
3. Add config schema in `pkg/config/config.go` under `tools`
4. Update `docs/tools_configuration.md`

## Build Tags

Default build tags: `goolm,stdjson`

Optional build tags (add to reduce binary size):
- `voice`: WebRTC/audio support (TTS, ASR, voice agent) - adds ~1MB and pion/webrtc dependency
- `selfupdate`: Self-update command - adds minio/selfupdate dependency
- `systray`: System tray (web launcher) - adds fyne.io/systray dependency
- `whatsapp_native`: WhatsApp native support via whatsmeow (larger binary)
- `bedrock`: AWS Bedrock provider support

Minimal build example:
```bash
go build -tags 'goolm,stdjson' ./cmd/picoclaw  # ~25MB
go build -tags 'goolm,stdjson,voice,selfupdate,systray' ./cmd/picoclaw  # ~26MB (full)
```

## Environment Variables

Key environment variables (see `pkg/config/envkeys.go`):
- `PICOCLAW_CONFIG`: Config file path
- `PICOCLAW_BINARY`: Path to main picoclaw binary (for launcher)
- `PICOCLAW_LAUNCHER_TOKEN`: Stable launcher auth token
- `PICOCLAW_GATEWAY_HOST`: Gateway host override (default 127.0.0.1)
- `PICOCLAW_LOG_LEVEL`: Gateway log verbosity

## Security Considerations

- AI-generated code requires extra scrutiny for path traversal, injection, credential exposure
- Config versioning separates sensitive data into `.security.yml`
- Sensitive data filtering prevents API keys from appearing in logs
- See `docs/security_configuration.md` and `docs/sensitive_data_filtering.md`

## Documentation Structure

Key docs in `docs/`:
- `providers.md`: Provider configuration details
- `chat-apps.md`: Channel setup guides
- `configuration.md`: Full config reference
- `tools_configuration.md`: Tool enable/disable and policies
- `cron.md`: Scheduled tasks
- `hooks/README.md`: Hook system (observers, interceptors, approval hooks)
- `steering.md`: Message injection between tool calls
- `subturn.md`: Subagent coordination

## Go Version

Requires Go 1.25+ (see `go.mod`). Line length limit: 120 chars.