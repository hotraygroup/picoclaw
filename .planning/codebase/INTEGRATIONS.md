# External Integrations

**Analysis Date:** 2026-05-11

## APIs & External Services

### LLM Providers (30+ protocols via abstraction layer):

PicoClaw uses a pluggable provider system at `pkg/providers/`. Each provider is specified in `model_list` using `protocol/model` format in `config.json`.

**Cloud Provider APIs:**
- **OpenAI** — GPT-5.4, GPT-4o, o3, etc.
  - Protocol: `openai/`
  - SDK: `github.com/openai/openai-go/v3`
  - Auth: `api_key` in `.security.yml` or `OPENAI_API_KEY` env var
  - Config: `~/.picoclaw/config.json` → `model_list[]`

- **Anthropic** — Claude Opus 4.6, Sonnet 4.6
  - Protocol: `anthropic/` or `anthropic-messages/`
  - SDK: `github.com/anthropics/anthropic-sdk-go`
  - Auth: `api_key` in `.security.yml` or `ANTHROPIC_API_KEY` env var

- **Google Gemini** — Gemini 3 Flash, 2.5 Pro
  - Protocol: `gemini/`
  - SDK: Custom HTTP implementation at `pkg/providers/httpapi/gemini_provider.go`
  - Auth: `api_key` in `.security.yml` or `GEMINI_API_KEY` env var

- **AWS Bedrock** — Claude, Llama, Mistral on AWS (conditional build)
  - Protocol: `bedrock/`
  - SDK: `github.com/aws/aws-sdk-go-v2/service/bedrockruntime`
  - Build tag: `bedrock` (not in default build)
  - Auth: AWS credentials (SDK default chain) or `AWS_REGION` env var
  - Files: `pkg/providers/bedrock/` (tagged `//go:build bedrock`)

- **Azure OpenAI** — Enterprise Azure deployment
  - Protocol: `azure/`
  - SDK: OpenAI-compatible via `pkg/providers/openai_compat/`
  - Auth: `api_key` in `.security.yml`

- **DeepSeek** — DeepSeek-V3, DeepSeek-R1
  - Protocol: `deepseek/`
  - SDK: OpenAI-compatible
  - Auth: `api_key` in `.security.yml`

- **Zhipu (GLM)** — GLM-4.7, GLM-5
  - Protocol: `zhipu/`
  - Auth: `api_key` in `.security.yml` or `ZHIPU_API_KEY` env var

- **OpenRouter** — 200+ models, unified API
  - Protocol: `openrouter/`
  - Auth: `api_key` in `.security.yml` or `OPENROUTER_API_KEY` env var

- **Volcengine (Doubao/Ark)**, **Qwen (DashScope)**, **Groq**, **Moonshot (Kimi)**, **Minimax**, **Mistral**, **NVIDIA NIM**, **Cerebras**, **Novita AI**, **Xiaomi MiMo**
  - All via OpenAI-compatible protocol adapter at `pkg/providers/openai_compat/`

**Local/On-Premise Providers:**
- **Ollama** — Local models, OpenAI-compatible
  - Protocol: `ollama/`
  - Auth: not required

- **vLLM** — Self-hosted OpenAI-compatible
  - Protocol: `vllm/`

- **LiteLLM** — Proxy for 100+ providers
  - Protocol: `litellm/`

- **LM Studio** — Local desktop deployment
  - Protocol: `lmstudio/`

**OAuth Providers:**
- **GitHub Copilot** — via device code OAuth flow
  - Protocol: `github-copilot/`
  - SDK: `github.com/github/copilot-sdk/go`
  - Auth: OAuth (device code), no API key needed
  - Files: `pkg/providers/oauth/`

- **Google Antigravity** — Google Cloud AI
  - Protocol: `antigravity/`
  - Auth: OAuth (OAuth2 via `golang.org/x/oauth2`)
  - Files: `pkg/providers/oauth/`

- **Claude CLI / Codex CLI** — CLI-based providers
  - Files: `pkg/providers/cli/`

### Chat Platform Channels (19+):

All channels are registered in `pkg/gateway/gateway.go:22-39` via blank imports. Active channels are enabled via `channels.*` config sections in `config.json`.

| Channel | Protocol | SDK | Config Key |
|---------|----------|-----|------------|
| **Telegram** | Long polling | `go.mau.fi/telego` | `channels.telegram` |
| **Discord** | WebSocket | `github.com/yeongaori/discordgo-fork` (fork) | `channels.discord` |
| **WhatsApp** | Native / Bridge WS | `go.mau.fi/whatsmeow` | `channels.whatsapp` |
| **Slack** | Socket Mode | `github.com/slack-go/slack` | `channels.slack` |
| **Matrix** | Sync API | `maunium.net/go/mautrix` | `channels.matrix` |
| **DingTalk** | Stream SDK | `github.com/open-dingtalk/dingtalk-stream-sdk-go` | `channels.dingtalk` |
| **Feishu / Lark** | WebSocket/SDK | `github.com/larksuite/oapi-sdk-go/v3` | `channels.feishu` |
| **LINE** | Webhook | `github.com/line/line-bot-sdk-go/v8` | `channels.line` |
| **QQ** | WebSocket | `github.com/tencent-connect/botgo` | `channels.qq` |
| **VK** | Long Poll | `github.com/SevereCloud/vksdk/v3` | `channels.vk` |
| **WeCom** | WebSocket | `github.com/open-dingtalk/dingtalk-stream-sdk-go` (AI Bot) | `channels.wecom` |
| **WeChat/Weixin** | iLink API (QR login) | Custom implementation | `channels.weixin` |
| **IRC** | IRC protocol | `github.com/ergochat/irc-go` | `channels.irc` |
| **OneBot** | WebSocket (OneBot v11) | Custom implementation | `channels.onebot` |
| **MQTT** | MQTT pub/sub | `github.com/eclipse/paho.mqtt.golang` | `channels.mqtt` |
| **MaixCam** | TCP socket | Custom implementation | `channels.maixcam` |
| **Pico** | Native WebSocket protocol | Custom, built-in | `channels.pico` |
| **Pico Client** | WebSocket client | Custom, built-in | `channels.pico_client` |
| **Teams Webhook** | Incoming webhook | `github.com/atc0005/go-teams-notify/v2` | `channels.teams_webhook` |

All webhook-based channels share a single Gateway HTTP server at `gateway.host:gateway.port` (default `127.0.0.1:18790`). Feishu uses SDK mode and runs its own WebSocket connection.

### Web Search Providers:

Configured in `tools.web` section. All search implementations are at `pkg/tools/search_tool.go`.

| Provider | Auth | Config Key |
|----------|------|------------|
| **DuckDuckGo** | None (free) | `tools.web.duckduckgo` |
| **Baidu Search** | API key | `tools.web.baidu_search` |
| **Tavily** | API key | `tools.web.tavily` |
| **Brave Search** | API key | `tools.web.brave` |
| **Perplexity** | API key | `tools.web.perplexity` |
| **SearXNG** | Self-hosted | `tools.web.searxng` |
| **GLM Search** | API key (Zhipu) | `tools.web.glm_search` |
| **Sogou** | Built-in scraping | `tools.web.sogou` |

### MCP (Model Context Protocol):

Native MCP client implementation at `pkg/mcp/manager.go`. Supports external MCP servers via stdio, SSE, and HTTP transports. Managed via `pkg/mcp/isolated_command_transport.go`.

- Protocol: MCP (Anthropic open standard)
- SDK: `github.com/modelcontextprotocol/go-sdk`
- Config: `tools.mcp.servers` in `config.json`
- Server types: `command` (stdio), `http` (HTTP/SSE)

### Skills Registries:

- **ClawHub** — `https://clawhub.ai` — skill marketplace
- **GitHub** — `https://github.com` — skill repos
- Config: `tools.skills.registries` in `config.json`

### WebRTC:

- `github.com/pion/webrtc/v3` v3.3.6 — Voice call / real-time media support
- `github.com/pion/rtp` v1.10.1 — RTP protocol for media streaming
- Files: `pkg/audio/`, `pkg/channels/voice_capabilities.go`

## Data Storage

**Databases:**
- **SQLite** (`modernc.org/sqlite` v1.48.2) — Pure Go embedded database
  - Files: `pkg/state/state.go`, `pkg/memory/jsonl.go`, `pkg/memory/store.go`
  - Used for: session state, JSONL memory store, password store, channel state, cron job data
  - Connection: file-based, located in `~/.picoclaw/workspace/`
  - Migration: `pkg/migrate/`, `pkg/memory/migration.go`

**File Storage:**
- Local filesystem only — `~/.picoclaw/workspace/` directory
- Media attachments: stored locally with cleanup (`tools.media_cleanup`)
- File operations: restricted to workspace with optional read/write path allowlists (`tools.allow_read_paths`, `tools.allow_write_paths`)
- Skills: loaded from `workspace/skills/` directory and built-in `skills/` directory
- Files: `pkg/media/store.go`, `pkg/tools/fs/`

**Caching:**
- Go build cache: `.cache/go-build/` (local)
- Go module cache: `.cache/go-mod/` (local)
- Skill search cache: in-memory LRU with TTL (`tools.skills.search_cache`)
- npm cache: Docker volume (`picoclaw-npm-cache`)

## Authentication & Identity

**Auth Provider:** Custom — no external auth service required

- **Dashboard Password:** bcrypt-hashed password for WebUI launcher (`web/backend/dashboardauth/`)
- **WebUI Local Auto-Login:** time-limited token for loopback access (`web/backend/middleware/`)
- **OAuth2:** Supported for GitHub Copilot (`github-copilot/`) and Google Antigravity (`antigravity/`) providers via `pkg/auth/oauth.go`
- **OAuth PKCE:** Implemented at `pkg/auth/pkce.go` for device code flows
- **Channel-level auth:** Each channel has its own token/auth mechanism (API keys, QR login, etc.)
- **Sensitive data:** Separated into `.security.yml` (YAML) with optional encryption (`pkg/config/security.go`)

**Cross-Agent Identity:**
- Concept of session/identity maintained at `pkg/identity/`
- Agent isolation via `pkg/isolation/`

## Monitoring & Observability

**Error Tracking:**
- None external — panic captures written to filesystem (`pkg/logger/`)
- Panic log: `~/.picoclaw/logs/panic.log`

**Logs:**
- Framework: `github.com/rs/zerolog` v1.35.1
- Structured JSON logging to file and/or console
- Gateway log: `~/.picoclaw/logs/gateway.log`
- Launcher log: `~/.picoclaw/logs/launcher.log`
- Log level configurable: `debug`, `info`, `warn`, `error`, `fatal` (env `PICOCLAW_LOG_LEVEL`)

**Health Checks:**
- Built-in HTTP health server: `GET /health` with JSON status response
- Configurable bearer token auth on health endpoint
- PID file tracking: `~/.picoclaw/.picoclaw.pid`
- Files: `pkg/health/server.go`

**Heartbeat:**
- Internal heartbeat system at `pkg/heartbeat/` (configurable interval, default 30s)

## CI/CD & Deployment

**Hosting:**
- Self-hosted — single binary deployment
- Docker images published to GHCR (`ghcr.io/sipeed/picoclaw`) and Docker Hub (`docker.io/sipeed/picoclaw`)
- Web download: `https://picoclaw.io/download/`
- Android APK: auto-generated in release workflow

**CI Pipeline:**
- GitHub Actions (`./github/workflows/`)
- Workflows:
  - `build.yml` — PR build checks
  - `pr.yml` — PR validation (lint, test, build)
  - `release.yml` — GoReleaser multi-platform release
  - `nightly.yml` — Nightly builds
  - `docker-build.yml` — Docker image builds
  - `stale.yml` — Issue/PR staleness
- Release: GoReleaser v2 handles cross-compilation, Docker images, packages (rpm/deb), macOS notarization

**Self-Update:**
- `github.com/minio/selfupdate` v0.6.0 — In-place binary self-update from GitHub Releases
- File: `pkg/updater/`

## Environment Configuration

**Required env vars (critical):**
- At least one LLM API key (e.g., `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`)
- Channel token for active chat platform (e.g., `TELEGRAM_BOT_TOKEN`)

**Optional env vars:**
- `PICOCLAW_HOME` — data directory (default: `~/.picoclaw`)
- `PICOCLAW_CONFIG` — config file path (default: `$PICOCLAW_HOME/config.json`)
- `PICOCLAW_GATEWAY_HOST` — gateway bind host (default: `localhost`)
- `PICOCLAW_LOG_LEVEL` — log level
- `TZ` — timezone (e.g., `Asia/Shanghai`)
- `BRAVE_SEARCH_API_KEY` — Brave web search

**Secrets location:**
- `.security.yml` alongside `config.json` in `~/.picoclaw/`
- Environment variables for provider API keys
- `.env` file supported (`.env.example` template provided)
- `docker/data/` volume for Docker deployments

## Webhooks & Callbacks

**Incoming:**
- All chat channel webhooks routed through Gateway HTTP server (`gateway.host:gateway.port`, default `127.0.0.1:18790`)
- LINE: `/webhook/line`
- Feishu: event subscription callback (WebSocket/SDK)
- Slack: Socket Mode events
- Other channels: via their native protocols (WebSocket, Long Polling, etc.)

**Outgoing:**
- Channel notification delivery via each platform's SDK
- Teams Incoming Webhook for notifications (`pkg/channels/teams_webhook/`)
- Execution results streamed via WebSocket (Pico Channel protocol)
- Cron job notifications via active channels (`pkg/cron/`)

---

*Integration audit: 2026-05-11*
