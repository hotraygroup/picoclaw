# Testing Patterns

**Analysis Date:** 2026-05-11

## Overview

The codebase has a well-established Go testing ecosystem with 284 test files. The web frontend (`web/frontend/`) has **no test files** — no vitest, jest, or other test runner is configured.

---

## Go Testing

### Test Framework

**Runner:**
- Standard `go test` via Makefile
- CI: GitHub Actions `pr.yml` runs `go test -tags goolm,stdjson ./...`

**Assertion Library:**
- `github.com/stretchr/testify` v1.11.1
  - `assert` package — for non-fatal assertions
  - `require` package — for fatal assertions (stops test immediately)

**Run Commands:**
```bash
make test                                          # Run all Go tests (excludes web/)
go test -tags goolm,stdjson ./...                  # Direct run with build tags
go test -tags goolm,stdjson ./pkg/credential/...   # Run specific package
```

### Test File Organization

**Location:**
- Co-located with source: `pkg/credential/credential.go` → `pkg/credential/credential_test.go`
- CLI tests co-located: `cmd/picoclaw/internal/auth/login.go` → `cmd/picoclaw/internal/auth/login_test.go`

**Naming:**
- `*_test.go` suffix (standard Go convention)
- Integration tests: `*_integration_test.go` (e.g., `pkg/config/migration_integration_test.go`)

**Package naming pattern:**
- Most tests use the same package as source: `package credential` (tests in `pkg/credential/store_test.go`)
- Some use external test package: `package credential_test` (e.g., `pkg/credential/credential_test.go`) — enforces testing only through public API

### Test Structure

**Table-driven tests** are the dominant pattern:

```go
// From pkg/logger/logger_test.go
func TestParseLevelValid(t *testing.T) {
    tests := []struct {
        input string
        want  LogLevel
    }{
        {"debug", DEBUG},
        {"DEBUG", DEBUG},
        {"Debug", DEBUG},
        {"info", INFO},
        // ...
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got, ok := ParseLevel(tt.input)
            if !ok {
                t.Fatalf("ParseLevel(%q) returned ok=false, want true", tt.input)
            }
            if got != tt.want {
                t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

**Test helper pattern:**
```go
// From pkg/agent/hooks_test.go
func newHookTestLoop(
    t *testing.T,
    provider providers.LLMProvider,
) (*AgentLoop, *AgentInstance, func()) {
    t.Helper()
    // setup...
    return al, agent, func() {
        al.Close()
        _ = os.RemoveAll(tmpDir)
    }
}
```

**Cleanup pattern:** Return a `func()` cleanup from helper, or use `t.TempDir()` for automatic tmpdir cleanup.

### Test Function Naming

- `Test{FunctionName}` — direct function test
- `Test{FunctionName}_{Scenario}` — scenario-specific: `TestResolve_PlainKey`, `TestResolve_FileKey_Success`
- `Test{FunctionName}_{EdgeCase}` — edge case: `TestResolve_FileKey_NotFound`, `TestResolve_EncKey_WrongPassphrase`
- `Test{MethodName}{Scenario}` — for methods: `TestAgentModelConfig_UnmarshalString`

### Assertion Patterns

**With testify (`assert` / `require`):**
```go
// From cmd/picoclaw/main_test.go
require.NotNil(t, cmd)
assert.Equal(t, "picoclaw", cmd.Use)
assert.Equal(t, short, cmd.Short)
assert.True(t, longHas)
assert.Len(t, subcommands, len(allowedCommands))
assert.False(t, subcmd.Hidden)
```

**Standard Go (without testify):**
```go
// From pkg/credential/credential_test.go
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
if got != "sk-plaintext-key" {
    t.Fatalf("got %q, want %q", got, "sk-plaintext-key")
}
```

**Mixed usage:** Some test files use plain `t.Fatal`/`t.Error`, others use `assert`/`require`. The `cmd/picoclaw/` tests lean heavily on testify, while `pkg/credential/` uses direct Go assertions.

### Mocking

**No dedicated mocking framework** detected (no testify/mock, no gomock). Mocks are hand-written:

```go
// From pkg/agent/mock_provider_test.go
type mockProvider struct{}

func (m *mockProvider) Chat(
    ctx context.Context,
    messages []providers.Message,
    tools []providers.ToolDefinition,
    model string,
    opts map[string]any,
) (*providers.LLMResponse, error) {
    return &providers.LLMResponse{
        Content:   "Mock response",
        ToolCalls: []providers.ToolCall{},
    }, nil
}

func (m *mockProvider) GetDefaultModel() string {
    return "mock-model"
}
```

**Pattern:** Define a struct implementing the interface directly in `_test.go` files. No interface generation tools used.

### Fixtures and Test Data

**Environment variables via `t.Setenv()`:**
```go
t.Setenv("PICOCLAW_SSH_KEY_PATH", sshKeyPath)
t.Setenv("PICOCLAW_KEY_PASSPHRASE", passphrase)
t.Setenv("HOME", "/tmp/home")
```

**Temporary directories:**
```go
dir := t.TempDir()  // auto-cleaned after test
tmpDir, err := os.MkdirTemp("", "agent-hooks-*")  // manual, with explicit cleanup
```

**Test config construction:**
```go
cfg := &config.Config{
    Agents: config.AgentsConfig{
        Defaults: config.AgentDefaults{
            Workspace:         tmpDir,
            ModelName:         "test-model",
            MaxTokens:         4096,
            MaxToolIterations: 10,
        },
    },
}
```

**Security-aware setup helpers:**
```go
func mustSetupSSHKey(t *testing.T) {
    t.Helper()
    keyPath := filepath.Join(t.TempDir(), "picoclaw_ed25519.key")
    if err := credential.GenerateSSHKey(keyPath); err != nil {
        t.Fatalf("mustSetupSSHKey: %v", err)
    }
    t.Setenv("PICOCLAW_SSH_KEY_PATH", keyPath)
}
```

### Coverage

**No coverage target enforced** in CI or Makefile. The `make test` command does not include `-cover` flag.

**To run with coverage manually:**
```bash
go test -tags goolm,stdjson -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Types

**Unit Tests:**
- Standard table-driven tests in `*_test.go` files
- Test individual functions and methods
- Package-internal or external test packages

**Integration Tests:**
- Suffix `_integration_test.go` (e.g., `pkg/config/migration_integration_test.go`, `pkg/config/security_integration_test.go`)
- Test cross-component interactions (config loading, migration, security)

**E2E Tests:**
- Not used in Go code

**Docker Tests:**
```bash
make docker-test   # Runs scripts/test-docker-mcp.sh — tests MCP tools in Docker
```

### Test Concurrency

- Individual tests use `t.Run()` subtests for parallelism when table-driven
- No explicit `t.Parallel()` usage detected broadly
- Some linter exclusions (e.g., `paralleltest` is disabled in `.golangci.yaml`)

### CI Test Configuration

**PR checks** (from `.github/workflows/pr.yml`):
```yaml
test:
  name: Tests
  runs-on: ubuntu-latest
  steps:
    - uses: actions/setup-go@v6
      with:
        go-version-file: go.mod
    - run: go generate ./...
    - run: go test -tags goolm,stdjson ./...
```

**Lint job** runs in parallel with tests:
```yaml
lint:
  - uses: golangci/golangci-lint-action@v9
    with:
      version: v2.10.1
      args: --build-tags=goolm,stdjson
```

**Security check** also runs in CI:
```yaml
vuln_check:
  - run: go install golang.org/x/vuln/cmd/govulncheck@v1.1.4
  - run: govulncheck -C . -format text ./...
```

### Test Level Parameters (Logger Tests)

When testing logger/comparable functions, tests capture and restore state:
```go
func TestLogLevelFiltering(t *testing.T) {
    initialLevel := GetLevel()
    defer SetLevel(initialLevel)    // Restore after test

    SetLevel(WARN)
    // ... test logic ...
}
```

### Common Testing Patterns Summary

| Pattern | File reference |
|---------|---------------|
| Table-driven tests | `pkg/logger/logger_test.go` |
| Test helper with `t.Helper()` | `pkg/agent/hooks_test.go:26` |
| Hand-written mock struct | `pkg/agent/mock_provider_test.go` |
| `t.Setenv()` for env config | `cmd/picoclaw/main_test.go` |
| `t.TempDir()` for temp dirs | `pkg/credential/credential_test.go` |
| testify `assert`/`require` | `cmd/picoclaw/main_test.go` |
| Security setup helper | `pkg/config/config_test.go:20` |
| Cleanup via returned `func()` | `pkg/agent/hooks_test.go:50` |
| Sentry-level state save/restore | `pkg/logger/logger_test.go:17-18` |

---

## Frontend Testing

**No test files exist** in `web/frontend/src/`. No vitest, jest, or testing-library configuration is present. The `package.json` scripts do not include a `test` command.

**Adding frontend tests would require:**
- Install `vitest`, `@testing-library/react`, `@testing-library/jest-dom`
- Add a `test` script to `package.json`
- Create test files following the co-location convention (e.g., `src/components/chat/chat-composer.test.tsx`)

---

*Testing analysis: 2026-05-11*
