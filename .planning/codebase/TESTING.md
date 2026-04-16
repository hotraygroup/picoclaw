# Testing Patterns

**Analysis Date:** 2026-04-16

## Test Framework

**Runner:**
- Go standard testing package
- testify/assert for assertions
- Config: No separate config file (standard Go test conventions)

**Assertion Library:**
- `github.com/stretchr/testify/assert` - primary assertion library
- Standard `testing` package for basic test operations

**Run Commands:**
```bash
make test               # Run all tests (Go + web frontend)
make check              # Full pre-commit: deps + fmt + vet + test
go test -run TestName -v ./pkg/session/  # Run single test with verbosity
go test -bench=. -benchmem -run='^$' ./...  # Run benchmarks with memory stats
```

## Test File Organization

**Location:**
- Test files co-located with source files
- Pattern: `<source>_test.go` in same directory
- Benchmark files: `<source>_bench_test.go`

**Naming:**
- Test functions: `Test<FunctionName>_<Scenario>`
- Example: `TestFilesystemTool_ReadFile_Success`, `TestFilesystemTool_ReadFile_NotFound`
- Benchmark functions: `Benchmark<FunctionName>_<Scenario>`
- Example: `BenchmarkIngest_SingleMessage`, `BenchmarkAssemble_MessagesOnly`

**Structure:**
```
pkg/
├── tools/
│   ├── filesystem.go
│   ├── filesystem_test.go     # Tests for filesystem.go
│   ├── registry.go
│   ├── registry_test.go       # Tests for registry.go
├── agent/
│   ├── loop.go
│   ├── loop_test.go           # Tests for loop.go
│   ├── mock_provider_test.go  # Mock implementations for tests
├── seahorse/
│   ├── store.go
│   ├── store_test.go
│   ├── short_bench_test.go    # Benchmarks
```

## Test Structure

**Suite Organization:**
```go
func TestFilesystemTool_ReadFile_Success(t *testing.T) {
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.txt")
    os.WriteFile(testFile, []byte("test content"), 0o644)

    tool := NewReadFileBytesTool("", false, MaxReadFileSize)
    ctx := context.Background()
    args := map[string]any{"path": testFile}

    result := tool.Execute(ctx, args)

    if result.IsError {
        t.Errorf("Expected success, got IsError=true: %s", result.ForLLM)
    }
}
```

**Patterns:**
- Setup: Use `t.TempDir()` for temporary directories (auto-cleaned)
- Helper functions: Mark with `t.Helper()` for better error traces
- Cleanup: Use `t.Cleanup()` for resource cleanup
- Table-driven tests: Used for testing multiple scenarios

**Helper Pattern from `pkg/agent/loop_test.go`:**
```go
func newStartedTestChannelManager(
    t *testing.T,
    msgBus *bus.MessageBus,
    store media.MediaStore,
    name string,
    ch channels.Channel,
) *channels.Manager {
    t.Helper()

    cm, err := channels.NewManager(&config.Config{}, msgBus, store)
    if err != nil {
        t.Fatalf("NewManager() error = %v", err)
    }
    cm.RegisterChannel(name, ch)
    if err := cm.StartAll(context.Background()); err != nil {
        t.Fatalf("StartAll() error = %v", err)
    }
    t.Cleanup(func() {
        if err := cm.StopAll(context.Background()); err != nil {
            t.Fatalf("StopAll() error = %v", err)
        }
    })
    return cm
}
```

## Mocking

**Framework:** Custom mock implementations (no external mock framework)

**Patterns:**
- Define mock types in test files with `_test.go` suffix
- Implement required interfaces inline
- Simple mocks return fixed values

**Example from `pkg/agent/mock_provider_test.go`:**
```go
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

**Recording Mock from `pkg/agent/loop_test.go`:**
```go
type recordingProvider struct {
    lastMessages []providers.Message
    lastModel    string
}

func (r *recordingProvider) Chat(...) (*providers.LLMResponse, error) {
    r.lastMessages = append([]providers.Message(nil), messages...)
    r.lastModel = model
    return &providers.LLMResponse{Content: "Mock response"}, nil
}
```

**What to Mock:**
- External services (LLM providers, databases)
- Network connections
- File system operations (when testing logic, not fs behavior)

**What NOT to Mock:**
- Standard library functions
- Simple utility functions
- Business logic being tested

## Fixtures and Factories

**Test Data:**
- Use `t.TempDir()` for isolated file fixtures
- Create test files inline with `os.WriteFile()`
- Define test structs inline for config tests

**Location:**
- Inline in test functions (most common)
- Helper functions for repeated setup (e.g., `newBenchStore`)

**Example from `pkg/seahorse/short_bench_test.go`:**
```go
func newBenchStore(b *testing.B) (*Store, func()) {
    b.Helper()
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        b.Fatalf("open test db: %v", err)
    }
    if err := runSchema(db); err != nil {
        db.Close()
        b.Fatalf("migration: %v", err)
    }
    return &Store{db: db}, func() { db.Close() }
}
```

## Coverage

**Requirements:** None explicitly enforced

**View Coverage:**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Types

**Unit Tests:**
- Focus on single function/module behavior
- Fast execution, no external dependencies
- Examples: `pkg/tools/filesystem_test.go`, `pkg/config/config_struct_test.go`

**Integration Tests:**
- Test interactions between components
- May use temporary databases, test servers
- Examples: `pkg/agent/loop_test.go`, `web/backend/api/*_test.go`

**E2E Tests:**
- Not used extensively in this codebase
- Web backend has API-level integration tests

## Common Patterns

**Async Testing:**
- Use context with timeout for async operations
- Example pattern:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
result, err := someAsyncOp(ctx)
```

**Error Testing:**
- Check both error presence and error content
- Example:
```go
result := tool.Execute(ctx, args)
if !result.IsError {
    t.Errorf("Expected error for missing file, got IsError=false")
}
if !strings.Contains(result.ForLLM, "failed to open file") {
    t.Errorf("Expected error message, got: %s", result.ForLLM)
}
```

**HTTP Handler Testing:**
- Use `httptest.NewRecorder` and `httptest.NewRequest`
- Example from `web/backend/api/pico_test.go`:
```go
func TestEnsurePicoChannel_FreshConfig(t *testing.T) {
    configPath := filepath.Join(t.TempDir(), "config.json")
    h := NewHandler(configPath)

    changed, err := h.EnsurePicoChannel("")
    if err != nil {
        t.Fatalf("EnsurePicoChannel() error = %v", err)
    }
    if !changed {
        t.Fatal("should report changed on a fresh config")
    }
}
```

**Benchmark Structure:**
```go
func BenchmarkIngest_SingleMessage(b *testing.B) {
    s, cleanup := newBenchStore(b)
    defer cleanup()
    ctx := context.Background()
    conv, _ := s.GetOrCreateConversation(ctx, "bench:ingest")
    convID := conv.ConversationID

    b.ResetTimer()  // Reset timer after setup
    for i := 0; i < b.N; i++ {
        _, err := s.AddMessage(ctx, convID, "user", "Test message", 15)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Web Backend Testing

**Location:**
- `web/backend/*_test.go` files co-located with source

**Test Commands:**
```bash
cd web && make test       # Backend Go tests + frontend lint
cd web/backend && go test ./...
```

**Patterns:**
- HTTP handler tests using `httptest`
- Config loading tests with temporary files
- Middleware tests for access control

## Frontend Testing

**Framework:** ESLint + Prettier for static analysis (no unit tests currently)

**Run Commands:**
```bash
cd web/frontend && pnpm lint    # ESLint check
cd web/frontend && pnpm check   # Prettier format + ESLint fix
```

**Test Scope:**
- Static analysis only (linting, formatting)
- Type checking via TypeScript compiler
- No runtime tests (React components, utilities)

## Test Directory Exclusions

**Ignored Paths:**
- `web/` directory excluded from main test runs
- Pattern in Makefile: `grep -v github.com/sipeed/picoclaw/web/`

**Platform-Specific Tests:**
- Build tag constrained tests (e.g., `//go:build linux`)
- Example: `pkg/tools/shell_timeout_unix_test.go` - Unix-only tests

---

*Testing analysis: 2026-04-16*