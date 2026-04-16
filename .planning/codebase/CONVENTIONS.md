# Coding Conventions

**Analysis Date:** 2026-04-16

## Naming Patterns

**Files:**
- Go files: lowercase with underscores (`filesystem.go`, `loop_test.go`)
- Test files: source file name + `_test.go` suffix
- Benchmark files: source file name + `_bench_test.go` suffix (e.g., `short_bench_test.go`)
- Frontend files: lowercase with underscores or kebab-case for directories

**Functions:**
- Public functions: PascalCase (e.g., `NewToolRegistry`, `ValidatePathWithAllowPaths`)
- Private functions: camelCase (e.g., `validatePathWithAllowPaths`, `isAllowedPath`)
- Constructors: `New<TypeName>` pattern (e.g., `NewToolRegistry`, `NewReadFileBytesTool`)
- Test functions: `Test<FunctionName>_<Scenario>` (e.g., `TestFilesystemTool_ReadFile_Success`)
- Benchmark functions: `Benchmark<FunctionName>_<Scenario>` (e.g., `BenchmarkIngest_SingleMessage`)

**Variables:**
- Local: camelCase
- Constants: PascalCase for exported, camelCase for private
- Interface types: often suffixed with `Aware` (e.g., `mediaStoreAware`)
- Mock types: `mock<Type>` (e.g., `mockProvider`)

**Types:**
- Structs: PascalCase (e.g., `ToolRegistry`, `AgentLoop`, `processOptions`)
- Interfaces: PascalCase, often named by capability (e.g., `LLMProvider`, `StreamingProvider`, `ThinkingCapable`)
- Type aliases: use `type <Name> = <Source>` for re-exports (e.g., in `pkg/providers/types.go`)
- Error types: `<Reason>Error` pattern (e.g., `FailoverError`)
- Config structs: `<Feature>Config` (e.g., `RoutingConfig`, `SubTurnConfig`)

**Packages:**
- Short, lowercase, single-word names (e.g., `agent`, `tools`, `config`, `logger`)
- Subdirectories for logical grouping (e.g., `audio/asr`, `audio/tts`)

## Code Style

**Formatting:**
- Tool: golangci-lint fmt (`make fmt`)
- Line length: 120 characters maximum
- Tab width: 4 spaces (configurable in `.golangci.yaml`)
- Use `any` instead of `interface{}`
- Simplify slice expressions: `a[b:]` instead of `a[b:len(a)]`

**Linting:**
- Tool: golangci-lint (`make lint`)
- Config: `.golangci.yaml`
- Auto-fix: `make fix` runs `golangci-lint run --fix`
- Key enabled linters (via `default: all`):
  - gofmt, gofumpt, goimports, gci (formatting)
  - golines (line length management)
  - Many others (see `.golangci.yaml` for full list)
- Key disabled linters: depguard, exhaustruct, gochecknoglobals, testpackage, varnamelen, wrapcheck

**Frontend Formatting:**
- Tool: Prettier (via `pnpm check` or `pnpm format`)
- Config: `web/frontend/prettier.config.js`
- Settings: semi=false, printWidth=80, tabWidth=2
- Import sorting: `@trivago/prettier-plugin-sort-imports`
- Order: builtin → third-party → `@/` → relative

**Frontend Linting:**
- Tool: ESLint (via `pnpm lint`)
- Config: `web/frontend/eslint.config.js`
- Plugins: react-hooks, react-refresh, typescript-eslint
- Ignores: `dist`, `src/components/ui`, `src/routeTree.gen.ts`

## Import Organization

**Go:**
- Managed by gci formatter with custom order:
  1. Standard library
  2. Default (third-party)
  3. Local module (`github.com/sipeed/picoclaw/...`)
- Example from `pkg/agent/loop.go`:
```go
import (
    "context"
    "encoding/json"
    "errors"
    // ... standard library ...
    
    "github.com/sipeed/picoclaw/pkg/audio/asr"
    "github.com/sipeed/picoclaw/pkg/audio/tts"
    // ... local packages ...
)
```

**Frontend:**
- Order defined in `prettier.config.js`:
  1. `<BUILTIN_MODULES>` (Node builtins)
  2. `<THIRD_PARTY_MODULES>` (npm packages)
  3. `^@/` (project aliases)
  4. `^[./]` (relative imports)

**Path Aliases:**
- Frontend: `@/*` maps to `./src/*` (defined in `tsconfig.json`)

## Error Handling

**Patterns:**
- Wrap errors with context: `fmt.Errorf("failed to ...: %w", err)`
- Custom error types with `Error()` and `Unwrap()` methods (e.g., `FailoverError`)
- Error classification via typed constants (e.g., `FailoverReason` enum)
- Return early with descriptive error messages

**Example from `pkg/tools/filesystem.go`:**
```go
if workspace == "" {
    return path, fmt.Errorf("workspace is not defined")
}
absWorkspace, err := filepath.Abs(workspace)
if err != nil {
    return "", fmt.Errorf("failed to resolve workspace path: %w", err)
}
```

**Error Types:**
- Custom errors implement `error` interface and often `Unwrap()` for error chain
- Example: `pkg/providers/types.go` defines `FailoverError` with reason classification

## Logging

**Framework:** zerolog (configured in `pkg/logger/logger.go`)

**Log Levels:**
- DEBUG, INFO, WARN, ERROR, FATAL
- Default: INFO
- Configurable via `SetLevel()` or `PICOCLAW_LOG_LEVEL` env var

**Logging Patterns:**
- Use component-specific logging: `logger.DebugCF("tools", "message", fields)`
- Field-based logging with `map[string]any`
- Console output with timestamps and caller info
- Optional file logging via `PICOCLAW_LOG_FILE`

**API:**
```go
// Basic logging
logger.Info("message")
logger.Errorf("formatted %s", arg)

// With component
logger.InfoC("tools", "registered tool")
logger.ErrorCF("gateway", "connection failed", map[string]any{"port": 8080})

// With fields
logger.DebugF("processing", map[string]any{"count": 10, "status": "active"})
```

**Component Convention:**
- Pass component name as first argument to `*C` variants
- Common components: `tools`, `gateway`, `agent`, `seahorse`

## Comments

**When to Comment:**
- Public API: document exported types, functions, constants
- Complex logic: explain non-obvious algorithms
- Config structs: document each field's purpose and defaults
- Build tags: document special compile conditions

**JSDoc/TSDoc:**
- Frontend uses `/** @type {import('...').Config} */` for type hints in JS config files

**Go Comments:**
- Package-level doc comments for major packages
- Struct field comments aligned after field definition
- Example from `pkg/config/config.go`:
```go
type ModelConfig struct {
    // Required fields
    ModelName string `json:"model_name"` // User-facing alias for the model
    Model     string `json:"model"`      // Protocol/model-identifier
}
```

**Test Comments:**
- Each test function has a descriptive comment explaining what it verifies
- Example: `// TestFilesystemTool_ReadFile_Success verifies successful file reading`

## Function Design

**Size:**
- Line limit: 120 (enforced by `.golangci.yaml` `funlen.lines: 120`)
- Statement limit: 40 (enforced by `funlen.statements: 40`)
- Complexity limit: 25 (gocognit.min-complexity)

**Parameters:**
- Argument limit: 7 (revive.argument-limit)
- Result limit: 3 (revive.function-result-limit)
- Use structs for options when many parameters needed (e.g., `processOptions`, `AssembleInput`)

**Return Values:**
- Return error as last value when applicable
- Use named returns sparingly (disabled linter: `nonamedreturns`)
- Helper functions use `t.Helper()` for better error trace

## Module Design

**Exports:**
- Export interfaces, not implementations (e.g., `LLMProvider` interface)
- Factory functions return interfaces where appropriate
- Private implementation details stay unexported

**Barrel Files:**
- Type alias re-exports for cleaner imports
- Example: `pkg/providers/types.go` re-exports types from `protocoltypes`
```go
type (
    ToolCall               = protocoltypes.ToolCall
    FunctionCall           = protocoltypes.FunctionCall
    // ...
)
```

**Package Organization:**
- One package per directory
- Related functionality grouped in subdirectories (e.g., `audio/asr`, `audio/tts`)
- Test files co-located with source

---

*Convention analysis: 2026-04-16*