# Daggerbot - AI Agent Instructions

## Build & Test Commands

- **Run tests**: `task test` or `go test ./...` (single package: `go test ./bot/commands`, single test: `go test -run TestName ./package`)
- **Lint/format**: `task check` (runs trunk check + formatters) or `trunk fmt`
- **Build**: `go build -o daggerbot .` or `task build`
- **Run locally**: `go run main.go` (requires `.env` with `DISCORD_AUTH_TOKEN` and optional `DISCORD_BOT_PREFIX`)
- **Docker build**: `docker compose up --build`

## Code Style

### Module & Imports

- Module path: `github.com/nerdwerx/daggerbot` (use for internal imports)
- Group imports with blank lines: stdlib → third-party (e.g., `github.com/bwmarrin/discordgo`) → local (e.g., `github.com/nerdwerx/daggerbot/config`)
- Enforced by `gofmt` and `golangci-lint` via trunk

### Types & Naming

- Exported: PascalCase (`Command`, `MessageSend`, `Handler`)
- Unexported: camelCase (`admin`, `handler`, `startup`)
- Use descriptive names; avoid single-letter vars except loops/closures

### Error Handling

- Check errors immediately: `if err != nil { return err }`
- Wrap with context: `fmt.Errorf("failed to X: %w", err)` or `fmt.Errorf("error X: %w", err)`
- `log.Fatal` only in `main()`/`init()`; return errors elsewhere
- Log errors before returning in handler/command functions

### Project Structure

- Commands: register via `RegisterCommand(NewCommand(...))` in package `init()` functions
- Handler signature: `func(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error`
- Config: use `config.Debug` and `config.Verbose` for conditional logging
- Message length: validate ≤2000 chars (Discord limit) in `MessageSend`/`MessagePrivateSend`
