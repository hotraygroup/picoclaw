# Coding Conventions

**Analysis Date:** 2026-05-11

## Overview

The codebase has two distinct language domains:
- **Go** — the core backend/CLI (`cmd/`, `pkg/`) with 284 test files
- **TypeScript/React** — the web frontend console (`web/frontend/`)

---

## Go Conventions

### Naming Patterns

**Files:**
- Snake_case: `store.go`, `http_client_test.go`, `agent_utils_test.go`
- Platform-specific suffixes: `panic_unix.go`, `panic_win.go`
- Test files: `*_test.go` co-located with source

**Functions:**
- Exported: PascalCase — `NewSecureStore()`, `GetConfigPath()`, `LoadConfig()`
- Unexported: camelCase — `formatFieldValue()`, `syncCliUIColor()`, `authLoginCmd()`
- Test functions: `TestXxx` — `TestResolve_PlainKey`, `TestLogLevelFiltering`
- Constructor pattern: `NewXxx(...)` returns `*Xxx`, e.g., `NewAuthCommand()`, `NewResolver(dir)`

**Variables:**
- Exported: PascalCase — `Logo`, `Component`, `ErrBusClosed`
- Unexported: camelCase — `currentLevel`, `consoleWriter`, `initMu`
- Constants: CamelCase or ALL_CAPS for legacy — `colorBlue`, `minWidthFancy`, `DEBUG`, `INFO`

**Types/Structs:**
- PascalCase: `SecureStore`, `AgentModelConfig`, `installedSkillOriginMeta`
- JSON struct tags use snake_case: `json:"origin_kind,omitempty"`

**Package Names:**
- Lowercase, single word: `credential`, `logger`, `config`, `skills`, `cliui`
- No underscores in package directory names

### Import Organization

Go imports are organized by `gci` into three groups (enforced by golangci-lint):
1. **Standard library:** `"os"`, `"fmt"`, `"sync"`
2. **Third-party:** `"github.com/rs/zerolog"`, `"github.com/spf13/cobra"`
3. **Local module:** `"github.com/sipeed/picoclaw/pkg/config"`

Example from `cmd/picoclaw/main.go`:
```go
import (
    "fmt"
    "os"
    "time"

    "github.com/spf13/cobra"

    "github.com/sipeed/picoclaw/cmd/picoclaw/internal"
    "github.com/sipeed/picoclaw/cmd/picoclaw/internal/agent"
    "github.com/sipeed/picoclaw/cmd/picoclaw/internal/auth"
    // ...
)
```

### Code Style (golangci-lint v2)

**Formatting (`.golangci.yaml`):**
- Formatters: `gci`, `gofmt` (with simplify), `gofumpt`, `goimports`, `golines`
- Max line length: **120** characters (enforced by `golines` and `lll` linter)
- Tab width: 4 spaces
- `gofmt` rewrite: `interface{}` → `any`, `a[b:len(a)]` → `a[b:]`

**Currently disabled linters (targeted for gradual re-enablement):**
- `errcheck`, `errorlint`, `gosec`, `staticcheck`, `revive`, `funlen`, `gocognit`, `gocyclo` — these are temporarily disabled with a note: "should fix and enable step by step"
- Pragmatic exclusions: `gosmopolitan` (legitimately uses CJK text), `testpackage` (uses internal test package pattern)

**nolint usage:**
- Inline for legitimate exceptions: `//nolint:gosmopolitan // We intentionally set local timezone from TZ env`

### Error Handling

**Error wrapping with `%w`:**
```go
// From pkg/credential/credential.go
var ErrPassphraseRequired = errors.New("credential: enc:// passphrase required")
var ErrDecryptionFailed = errors.New("credential: enc:// decryption failed (wrong passphrase or SSH key?)")

// Wrapping sentinel errors
return "", fmt.Errorf("%w: %w", ErrDecryptionFailed, err)

// From cmd/picoclaw/internal/cron/add.go
return fmt.Errorf("error adding job: %w", err)

// From pkg/logger/logger.go
return fmt.Errorf("failed to create log directory: %w", err)
```

**Sentinel errors (exported, package-level):**
- `pkg/credential/credential.go`: `ErrPassphraseRequired`, `ErrDecryptionFailed`
- `pkg/bus/bus.go`: `ErrBusClosed`
- `pkg/channels/errutil.go`: `ErrRateLimit`, `ErrTemporary`, `ErrSendFailed`

**Error return pattern:**
- Functions return `(result, error)` — last return value is error
- CLI commands use `RunE` (returns error) rather than `Run`
- Errors are often checked inline: `if err != nil { return err }`

### Logging

**Framework:** Custom wrapper around `github.com/rs/zerolog` v1.35.1

**Package:** `pkg/logger`

**Level functions (from `pkg/logger/logger.go`):**
```
Debug(message)     Info(message)     Warn(message)     Error(message)     Fatal(message)
DebugC(comp, msg)  InfoC(comp, msg)  WarnC(comp, msg)  ErrorC(comp, msg)  FatalC(comp, msg)
Debugf(fmt, ...)   Infof(fmt, ...)   Warnf(fmt, ...)   Errorf(fmt, ...)   Fatalf(fmt, ...)
DebugF(msg, flds)  InfoF(msg, flds)  WarnF(msg, flds)  ErrorF(msg, flds)  FatalF(msg, flds)
DebugCF(c, m, f)   InfoCF(c, m, f)   WarnCF(c, m, f)   ErrorCF(c, m, f)   FatalCF(c, m, f)
```

**Naming convention:** Suffix patterns for convenience functions:
- `C` = with Component name
- `f` = with format string (sprintf-like)
- `F` = with Fields map (`map[string]any`)
- `CF` = with both Component and Fields

**Pattern when logging:**
```go
logger.Info("Starting picoclaw gateway")
logger.ErrorC("config", "Failed to load config")
logger.ErrorF("Failed to send", map[string]any{"error": err, "channel": channelName})
```

**Output:** Console (with color if TTY) and optional file (`PICOCLAW_LOG_FILE` env var)

### Comments / Documentation

- **Package comments:** Standard Go convention — first line of `cliui/cliui.go`: `// Package cliui renders human-oriented CLI output...`
- **Function comments:** Document exported functions — `// GetPicoclawHome returns the picoclaw home directory.`
- **Deprecation:** `// Deprecated: Use pkg/config.FormatVersion instead`
- **Inline comments:** Used for clarifying non-obvious logic — `// Round-robin: 依次尝试不同的 DNS 服务器`

### Module Design

**Exports pattern:**
- Exported types/functions in package root files
- Internal details in `helpers.go` files (unexported helper functions)
- `internal/` directory used for application-specific CLI commands (`cmd/picoclaw/internal/`)

**Package structure:**
```
pkg/
  {feature}/          # Feature-specific package
    {feature}.go      # Core logic
    {feature}_test.go # Tests co-located
    helpers.go        # Internal helpers (sometimes)
```

### CLI Patterns

**Framework:** `github.com/spf13/cobra` v1.10.2

**Command construction:**
- Each subcommand gets its own directory under `internal/`
- Constructor function `NewXxxCommand()` returns `*cobra.Command`
- Subcommand files: `command.go` (registers subcommands), `helpers.go` (logic), `{action}.go` (subcommand definition)
- Example: `internal/cron/` has `command.go`, `add.go`, `list.go`, `remove.go`, `enable.go`, `disable.go`, `helpers.go`

**Argument pattern:**
```go
func NewAuthCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "auth",
        Short: "Manage authentication (login, logout, status)",
        RunE: func(cmd *cobra.Command, _ []string) error {
            return cmd.Help()
        },
    }
    cmd.AddCommand(newLoginCommand(), newLogoutCommand(), ...)
    return cmd
}
```

### Key Go Dependencies

- **CLI:** `github.com/spf13/cobra`, `github.com/spf13/pflag`
- **Logging:** `github.com/rs/zerolog`
- **UI:** `github.com/charmbracelet/lipgloss`, `github.com/muesli/termenv`
- **SQLite:** `modernc.org/sqlite` (CGo-free)
- **Config:** `github.com/caarlos0/env/v11`
- **Testing:** `github.com/stretchr/testify` (assert + require)
- **YAML:** `gopkg.in/yaml.v3`

---

## TypeScript/React Conventions

### Naming Patterns

**Files:**
- kebab-case: `chat-composer.tsx`, `ansi-log-line.tsx`, `use-theme.ts`
- Route files in `src/routes/` use `.tsx` extension
- UI components in `src/components/ui/` use shadcn convention (kebab-case)

**Functions/Components:**
- React components: PascalCase — `ChatComposer`, `AppSidebar`, `ContextUsageRing`
- Hooks: camelCase with `use` prefix — `useTheme`, `useMobile`, `useGateway`
- Utility functions: camelCase — `cn()`, `formatFieldValue()`
- API modules: camelCase file name, default export — `api/gateway.ts`, `api/sessions.ts`

**Types/Interfaces:**
- PascalCase: `ChatComposerProps`, `ChatInputDisabledReason`, `ContextUsage`
- Union types: PascalCase — `ChatInputDisabledReason`
- Type imports use `import type`: `import type { ChatAttachment, ContextUsage } from "@/store/chat"`

**Variables:**
- camelCase: `disabledMessage`, `canInput`, `contextUsage`
- Constants in kebab-case files: no special convention detected

### Import Organization

**Prettier plugin `@trivago/prettier-plugin-sort-imports` enforces:**
1. Built-in modules: `"path"`, `"react"`
2. Third-party modules: `"@tabler/icons-react"`, `"react-textarea-autosize"`
3. `@/` path alias (local source): `@/components/chat/context-usage-ring`
4. Relative imports: `"./models"`, `"../store/chat"`

**Path alias:** `@/` → `src/` (configured in both `tsconfig.json` and `vite.config.ts`)

Example from `src/components/chat/chat-composer.tsx`:
```typescript
import { IconArrowUp, IconPhotoPlus, IconX } from "@tabler/icons-react"
import type { KeyboardEvent } from "react"
import { useTranslation } from "react-i18next"
import TextareaAutosize from "react-textarea-autosize"

import { ContextUsageRing } from "@/components/chat/context-usage-ring"
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import type { ChatAttachment, ContextUsage } from "@/store/chat"
```

### Code Style (Frontend)

**Formatting (`web/frontend/prettier.config.js`):**
- No semicolons (`semi: false`)
- Print width: 80 characters
- Tab width: 2 spaces
- Import sorting enabled with separation between groups

**Linting (`web/frontend/eslint.config.js`):**
- Base: `@eslint/js` recommended + `typescript-eslint` recommended
- React: `eslint-plugin-react-hooks` (flat recommended), `eslint-plugin-react-refresh` (Vite)
- Prettier integration: `eslint-config-prettier`
- Global ignores: `dist/`, `src/components/ui/` (shadcn generated), `src/routeTree.gen.ts`
- Route files (`src/routes/**`): exempt from `react-refresh/only-export-components` (TanStack Router false positives)

### Component Patterns

**Functional components with destructured props:**
```tsx
export function ChatComposer({
  input,
  attachments,
  onInputChange,
  onAddImages,
  onRemoveAttachment,
  onSend,
  onContextDetail,
  inputDisabledReason,
  canSend,
  contextUsage,
}: ChatComposerProps) {
  // ...
}
```

**Utility function:** `cn()` — combines `clsx` + `tailwind-merge` for conditional class names
```typescript
// From src/lib/utils.ts
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

### State Management

- **Jotai** (`jotai`) for global state — `src/store/gateway.ts`, `src/store/chat.ts`
- **TanStack React Query** (`@tanstack/react-query`) for server state
- **TanStack Router** for routing — `src/routes/`

### Internationalization

- `react-i18next` + `i18next` — `src/i18n/index.ts`
- Translation keys use dot notation: `t("chat.placeholder")`, `t("chat.sendHint")`

### Key Frontend Dependencies

- **Framework:** React 19.2.5, Vite 8
- **UI:** shadcn/ui (Radix UI primitives), Tailwind CSS v4
- **Routing:** @tanstack/react-router v1.169
- **Icons:** @tabler/icons-react
- **Markdown:** react-markdown, rehype-highlight, remark-gfm
- **Date:** dayjs
- **Toasts:** sonner

---

## Configuration

**Go:**
- Config file: `.golangci.yaml` (golangci-lint v2)
- Build tags: `goolm,stdjson` (set in `Makefile`)
- Format: `gofmt -s`, `gofumpt`, `gci`, `golines`
- Line length: 120

**Frontend:**
- ESLint: `web/frontend/eslint.config.js` (flat config, ESLint 10)
- Prettier: `web/frontend/prettier.config.js` (v3.8.3)
- TypeScript: `web/frontend/tsconfig.json` (project references)
- Package manager: pnpm 10.33.0 (enforced in `package.json` `packageManager` field)

**Format commands:**
```bash
# Go
make fmt      # runs golangci-lint fmt
make fix      # golangci-lint --fix

# Frontend
cd web/frontend && pnpm check   # prettier --write && eslint --fix
cd web/frontend && pnpm format  # prettier --check
cd web/frontend && pnpm lint    # eslint
```

---

*Convention analysis: 2026-05-11*
