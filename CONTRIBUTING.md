# Contributing to Alchemist

Alchemist is built and maintained by our team, and contributions from everyone
on it are welcome — including (especially!) first-time ones.

## Getting set up

```
git clone https://github.com/TheGrimmClub/alchemist.git
cd alchemist
go mod tidy
go run . --help
```

## Workflow

We use Alchemist to build Alchemist. The expected flow for any change:

1. `alchemist start` — name what you're about to do
2. Make your change
3. `gofmt -w .` and `go vet ./...`
4. `alchemist brew` — commit with a clear message
5. Open a pull request

Keep pull requests small and focused on one task.

## Project layout

```
main.go                 entry point
cmd/                    one file per command (start, brew, bottle, ...)
internal/state/         persists the current task
internal/gitutil/       thin wrappers around git
internal/editor/        opens the commit message in an editor
internal/prompt/        small stdin prompt helpers
```

## Adding a command

1. Create `cmd/yourcommand.go` defining a `*cobra.Command`.
2. Register it in the `init()` of `cmd/root.go`.
3. Add a short entry to the command list in `README.md`.

## Style

- Run `gofmt` before committing.
- Keep user-facing messages short and encouraging — beginners read them.
- Prefer adding helpers to `internal/` over duplicating git calls.

## Reporting issues

Open a GitHub issue with what you expected, what happened, your OS, and the
output of `alchemist --version` (once versioning lands). Reproduction steps
help a lot.
