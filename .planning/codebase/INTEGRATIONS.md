# External Integrations

**Analysis Date:** 2026-04-16

## APIs & External Services

### LLM Providers (30+ Supported)

**Direct API Providers:**
- OpenAI (`openai/`) - GPT models, API key auth
- Anthropic (`anthropic/`, `anthropic-messages/`) - Claude models, API key auth
- Google Gemini (`gemini/`, `antigravity/`) - Gemini models, OAuth or API key
- DeepSeek (`deepseek/`) - DeepSeek models, API key
- Groq (`groq/`) - Fast inference, also Whisper transcription
- Cerebras (`cerebras/`) - Fast inference
- Mistral (`mistral/`) - Mistral models
- Moonshot (`moonshot/`) - Kimi models (Chinese)
- Zhipu AI (`zhipu/`) - GLM models (Chinese)
- Qwen/Alibaba (`qwen/`, `qwen-intl/`, `coding-plan/`) - Qwen models
- NVIDIA NIM (`nvidia/`) - NVIDIA inference
- Venice AI (`venice/`) - Uncensored models
- Minimax (`minimax/`) - Minimax models
- LongCat (`longcat/`) - LongCat models
- ModelScope (`modelscope/`) - ModelScope inference (Chinese)
- Xiaomi MiMo (`mimo/`) - MiMo models
- VolcEngine/Doubao (`volcengine/`) - ByteDance models
- Vivgrid (`vivgrid/`) - Vivgrid models
- Avian (`avian/`) - Avian models

**Aggregator/Proxy Providers:**
- OpenRouter (`openrouter/`) - Multi-model aggregator
- LiteLLM Proxy (`litellm/`) - Local proxy
- VLLM (`vllm/`) - Local inference server

**Cloud Providers:**
- Azure OpenAI (`azure/`) - Azure-hosted GPT, deployment-based
- AWS Bedrock (`bedrock/`) - AWS-hosted Claude, Titan, etc.

**Local Providers:**
- Ollama (`ollama/`) - Local inference, no API key needed
- LM Studio (`lmstudio/`) - Local inference, optional key

**CLI/Built-in Providers:**
- GitHub Copilot (`github-copilot/`) - Local gRPC, OAuth via VS Code
- Claude CLI (`claude-cli/`) - Anthropic OAuth flow
- Codex CLI (`codex-cli/`) - OpenAI OAuth flow

**SDK/Client Files:**
- OpenAI/compatible: `pkg/providers/openai_compat/provider.go`
- Anthropic: `pkg/providers/anthropic/provider.go`
- Anthropic Messages: `pkg/providers/anthropic_messages/provider.go`
- Azure: `pkg/providers/azure/provider.go`
- Bedrock: `pkg/providers/bedrock/provider.go`
- Gemini/Antigravity: `pkg/providers/antigravity_provider.go`
- Claude CLI: `pkg/providers/claude_cli_provider.go`
- Codex CLI: `pkg/providers/codex_cli_provider.go`

### Web Search Providers

**Search APIs:**
- DuckDuckGo - Free, no API key (HTML scraping)
- Brave Search - API key required (`brave.api_key`)
- Tavily - API key required (`tavily.api_key`)
- Perplexity - LLM-based search, API key (`perplexity.api_key`)
- Sogou - Free, Chinese search (HTML scraping)
- SearXNG - Self-hosted meta search (`searxng.base_url`)
- GLM Search (Zhipu) - API key (`glm_search.api_key`)
- Baidu Search - API key (`baidu_search.api_key`)

**Config Location:**
- `config/tools.web.*` in `config/config.example.json`
- Implementation: `pkg/tools/web.go`, `pkg/tools/search_tool.go`

## Data Storage

### Databases:
- SQLite (pure Go via `modernc.org/sqlite`)
  - Session storage: `pkg/session/jsonl_backend.go`
  - Auth store: `web/backend/dashboardauth/store.go`
  - Memory storage: `pkg/memory/jsonl.go`

### File Storage:
- JSONL files for conversation memory: `~/.picoclaw/memory/`
- Workspace directory: `~/.picoclaw/workspace/`
- Skills directory: `~/.picoclaw/workspace/skills/`

### Caching:
- In-memory caching for model discovery (TTL configurable)
- Search cache: `tools.skills.search_cache` with TTL

## Messaging Channels (Chat Apps)

### Supported Platforms:

| Channel | Protocol | Auth Method | Files |
|---------|----------|-------------|-------|
| Telegram | Bot API HTTP | Bot token | `pkg/channels/pico/` + SDK integration |
| Discord | WebSocket | Bot token | Uses `bwmarrin/discordgo` |
| WhatsApp | Native/Bridge | QR scan | Uses `whatsmeow` |
| Slack | Socket Mode | Bot + App token | Uses `slack-go/slack` |
| Matrix | WebSocket | Access token | Uses `mautrix` |
| Feishu/Lark | Stream API | App ID + Secret | Uses `larksuite/oapi-sdk-go` |
| DingTalk | Stream API | Client ID + Secret | Uses `dingtalk-stream-sdk-go` |
| QQ | Official Bot API | App ID + Secret | Uses `tencent-connect/botgo` |
| LINE | HTTP Webhook | Channel secret + token | Requires HTTPS |
| WeCom | WebSocket AI Bot | Bot ID + Secret | `pkg/channels/pico/` |
| IRC | TCP + TLS | Nick + Password | Uses `ergochat/irc-go` |
| OneBot | WebSocket | Access token | NapCat/Go-CQHTTP compatible |
| VK | Antigravity | OAuth | Uses `vksdk` |
| MaixCam | Custom | Custom | Hardware integration |
| Pico | WebSocket | Token | Native protocol, `pkg/channels/pico/` |

**Channel Config Location:**
- `config/channels.*` in `pkg/config/config_channel.go`

## Authentication & Identity

### Dashboard Authentication:
- Token-based auth (generated on first run)
- Optional password store (bcrypt via SQLite)
- Session cookies with signing key
- Files: `web/backend/api/auth.go`, `web/backend/dashboardauth/store.go`

### OAuth Flows:
- Anthropic OAuth (`claude-cli` provider)
- OpenAI OAuth (`codex-cli` provider)
- Google OAuth (Antigravity/Gemini)
- WeChat QR login (`weixin` flow)
- WeCom QR login (`wecom` flow)

### Provider Auth Methods:
- `api_key` - Standard API key
- `oauth` - OAuth token refresh
- `token` - Static bearer token

**Auth Implementation:**
- Token sources: `pkg/providers/common/token_source.go`
- Credential storage: Encrypted in `.security.yml`

## Monitoring & Observability

### Error Tracking:
- None (native Go error handling)

### Logs:
- rs/zerolog structured logging
- File logging: `~/.picoclaw/logs/launcher.log`
- Panic logs: `~/.picoclaw/logs/launcher_panic.log`
- Gateway log level configurable (`gateway.log_level`)

### Health:
- Heartbeat system (`heartbeat.enabled`, `heartbeat.interval`)

## Web APIs Exposed

### Gateway Server (Default: localhost:18790):
- WebSocket endpoint: `/pico/ws` - Chat traffic proxy
- HTTP webhook endpoints for LINE, etc.

### Web Launcher Backend (Default: localhost:18800):
- Config API: `/api/config/*` - CRUD operations
- Session API: `/api/session/*` - History, management
- OAuth API: `/api/oauth/*` - Provider login flows
- Gateway API: `/api/gateway/*` - Process lifecycle
- Skills API: `/api/skills/*` - Registry operations
- Tools API: `/api/tools/*` - Tool configuration
- Channel API: `/api/channels/*` - Channel catalog
- Auth API: `/api/auth/*` - Dashboard login

**Router Definition:**
- `web/backend/api/router.go` - Route registration
- Individual handlers in `web/backend/api/*.go`

## Model Context Protocol (MCP)

### MCP Integration:
- SDK: `github.com/modelcontextprotocol/go-sdk v1.5.0`
- Implementation: `pkg/tools/mcp_tool.go`

### Pre-configured MCP Servers:
- context7 - HTTP MCP server
- filesystem - Local file operations (`npx @modelcontextprotocol/server-filesystem`)
- github - GitHub API (`npx @modelcontextprotocol/server-github`)
- brave-search - Brave search integration
- postgres - PostgreSQL database
- slack - Slack API integration

**MCP Config Location:**
- `config/tools.mcp.servers.*` in `config/config.example.json`

## Skills Registry

### External Registries:
- ClawHub (`https://clawhub.ai`) - Official skills registry
- GitHub - Skills from GitHub repos

**Registry Config:**
- `config/tools.skills.registries.*` in `config/config.example.json`
- Implementation: `pkg/tools/skills_install.go`, `pkg/tools/skills_search.go`

## CI/CD & Deployment

### Hosting:
- Self-hosted (local binary)
- Docker support (`docker/Dockerfile`)
- Cross-platform builds via goreleaser

### CI Pipeline:
- GitHub Actions (`.github/workflows/`)
- Pre-commit checks: deps + fmt + vet + test

### Self-Update:
- Built-in update mechanism (`github.com/minio/selfupdate`)
- Update endpoint: `web/backend/api/update.go`

## Environment Configuration

### Required env vars (Runtime):
- `PICOCLAW_HOME` - Base directory (optional, defaults to ~/.picoclaw)
- `PICOCLAW_CONFIG` - Config path (optional)
- Provider API keys (in `.security.yml` or config.json)

### Secrets Location:
- `~/.picoclaw/.security.yml` - Encrypted API keys, tokens
- Sensitive data filtering prevents exposure in logs
- Config versioning: v0 (legacy) → v1+ (secrets in .security.yml)

## Webhooks & Callbacks

### Incoming:
- LINE webhook: `/webhook/line` (requires HTTPS)
- Gateway HTTP server for callback-based channels

### Outgoing:
- None (agent sends messages via channels, not callbacks)

---

*Integration audit: 2026-04-16*