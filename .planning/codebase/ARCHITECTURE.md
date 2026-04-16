# Architecture

**Analysis Date:** 2026-04-16

## Pattern Overview

**Overall:** Modular monolith with event-driven messaging bus

**Key Characteristics:**
- Gateway-centric orchestration pattern
- Message bus for decoupled communication between layers
- Plugin-style provider system supporting 30+ LLM APIs
- Tool registry with dynamic registration and TTL-based promotion
- Multi-agent support with routing (light/heavy model tiering)
- Channel abstraction for multi-platform messaging (currently pico only)

## Layers

### Gateway Layer
- Purpose: Runtime orchestration, service lifecycle management, hot reload
- Location: `pkg/gateway/gateway.go`
- Contains: Service coordination, config reload, PID management, listener binding
- Depends on: agent, channels, config, providers, tools, cron, heartbeat
- Used by: CLI command `gateway` subcommand

### Agent Loop Layer
- Purpose: Core AI conversation loop, tool execution, context management
- Location: `pkg/agent/loop.go`, `pkg/agent/instance.go`
- Contains: AgentLoop struct, turn processing, tool orchestration, steering queue
- Depends on: providers, tools, bus, config, session, memory, hooks
- Used by: Gateway layer

### Provider Layer
- Purpose: LLM API abstraction and fallback handling
- Location: `pkg/providers/`
- Contains: Protocol implementations (OpenAI, Anthropic, Gemini, Azure, Bedrock, etc.), fallback chain, rate limiting
- Depends on: config, protocol types
- Used by: Agent loop via `LLMProvider` interface

### Channel Layer
- Purpose: Messaging platform abstraction (inbound/outbound routing)
- Location: `pkg/channels/`
- Contains: Channel interface, Manager, dynamic mux for HTTP routing
- Depends on: bus, config, media
- Used by: Gateway layer for HTTP server setup

### Tool Layer
- Purpose: Action execution (filesystem, shell, web search, MCP, etc.)
- Location: `pkg/tools/`
- Contains: Tool interface, registry, concrete tool implementations
- Depends on: providers (for tool definitions), config, media
- Used by: Agent loop via registry

### Bus Layer
- Purpose: Decoupled message passing between components
- Location: `pkg/bus/bus.go`
- Contains: MessageBus with inbound/outbound/media channels, streaming delegate
- Depends on: None (pure infrastructure)
- Used by: Agent loop, Channel manager, services

### Session/Memory Layer
- Purpose: Conversation persistence and history management
- Location: `pkg/session/`, `pkg/memory/`
- Contains: SessionStore interface, JSONL backend, session key management
- Depends on: providers (for message types)
- Used by: Agent instances for context building

## Data Flow

### Inbound Message Flow

1. External platform sends HTTP request to Channel Manager
2. Channel Manager normalizes message and publishes `InboundMessage` to bus
3. Agent Loop receives from bus `InboundChan()`, creates dispatch request
4. Router classifies message complexity (light vs. heavy model)
5. Context Manager builds conversation history from session/memory
6. Provider receives messages + tools, returns LLM response
7. If tool calls present, Agent Loop executes tools iteratively
8. Response published as `OutboundMessage` to bus
9. Channel Manager receives from bus, sends to platform

### Hot Reload Flow

1. Config file watcher detects modification (2s polling)
2. Gateway loads new config, validates
3. Stops all services (cron, heartbeat, channels, media)
4. Creates new provider with new model config
5. Agent loop reloads provider and config atomically
6. Restart all services with new configuration
7. Log level updated if changed

### Tool Execution Flow

1. LLM returns tool call in response
2. Agent loop extracts tool name and arguments
3. Tool registry lookup by name
4. Tool.Execute() called with context and args
5. Result converted to ToolResult struct
6. Published as assistant message + tool result to context
7. Next LLM call includes tool result for continuation

## Key Abstractions

### LLMProvider Interface
- Purpose: Abstract all LLM API interactions
- Examples: `pkg/providers/types.go`
- Pattern: Interface with Chat() and GetDefaultModel() methods
- Extensions: StreamingProvider, ThinkingCapable, NativeSearchCapable

### Tool Interface
- Purpose: Abstract executable actions for AI
- Examples: `pkg/tools/base.go`
- Pattern: Interface with Name(), Definition(), Execute() methods
- Registry pattern with TTL-based promotion for hidden tools

### Channel Interface
- Purpose: Abstract messaging platform operations
- Examples: `pkg/channels/base.go`
- Pattern: Interface with Send(), Get(), Edit(), React() methods
- Manager pattern with per-channel workers and rate limiting

### MessageBus
- Purpose: Decouple components via typed channels
- Examples: `pkg/bus/bus.go`
- Pattern: Multi-channel async queue (inbound, outbound, media, audio)
- StreamDelegate for streaming-capable channels

### AgentInstance
- Purpose: Fully configured agent with workspace, tools, sessions
- Examples: `pkg/agent/instance.go`
- Pattern: Composite struct holding all agent dependencies
- Pre-computed candidates for fallback chain

### ContextManager
- Purpose: Build and manage conversation context windows
- Examples: `pkg/agent/context_manager.go`
- Pattern: Strategy pattern (legacy, budget, seahorse)
- Budget-based context trimming for large conversations

## Entry Points

### CLI Entry Point
- Location: `cmd/picoclaw/main.go`
- Triggers: User command execution
- Responsibilities: Cobra command tree, banner display, timezone setup

### Gateway Entry Point
- Location: `pkg/gateway/gateway.go` → `Run()`
- Triggers: `picoclaw gateway -E` or web launcher
- Responsibilities: Service orchestration, config loading, hot reload

### Web Launcher Entry Point
- Location: `web/backend/main.go`
- Triggers: `picoclaw-launcher` binary or desktop entry
- Responsibilities: Spawn gateway subprocess, WebSocket proxy, dashboard API

### Agent Loop Entry Point
- Location: `pkg/agent/loop.go` → `Run(ctx)`
- Triggers: Gateway startup
- Responsibilities: Consume inbound bus, process turns, execute tools

## Error Handling

**Strategy:** Fallback chain with classification

**Patterns:**
- `FailoverError` wraps provider errors with classification (auth, rate_limit, timeout, etc.)
- `FallbackChain` manages cooldown and rate limiting across candidates
- Non-retriable errors (format, context_overflow) skip fallback
- Graceful degradation with startup blocked provider when no model configured

## Cross-Cutting Concerns

**Logging:** Structured logging via `pkg/logger/` with file rotation and panic capture
**Validation:** Config validation in `pkg/config/config.go` with model list checking
**Authentication:** OAuth/token auth via `pkg/auth/` for Claude CLI and Codex providers
**Rate Limiting:** Per-provider RPM limits via `pkg/providers/ratelimiter.go`
**Sensitive Data:** Filtering of API keys/credentials before LLM via `pkg/config/config.go` → `FilterSensitiveData()`
**Isolation:** Subprocess isolation config via `pkg/isolation/` for tool execution

## Special Patterns

### SubAgent Coordination
- Location: `pkg/agent/subturn.go`
- Purpose: Delegate subtasks to specialized agents
- Pattern: SubTurn struct with parent context, result collection

### Steering Queue
- Location: `pkg/agent/steering.go`
- Purpose: Inject messages between tool calls
- Pattern: Queue-based message injection for refactor/agent guidance

### Hook System
- Location: `pkg/agent/hooks.go`, `pkg/agent/hook_process.go`
- Purpose: Event observation and interception
- Pattern: Observer/Interceptor hooks with priority ordering
- External process hooks via stdin/stdout transport

### MCP Integration
- Location: `pkg/mcp/manager.go`, `pkg/tools/mcp_tool.go`
- Purpose: Model Context Protocol for external tool servers
- Pattern: Process-based transport, isolated command handling

---

*Architecture analysis: 2026-04-16*