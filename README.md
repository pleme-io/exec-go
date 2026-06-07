# exec-go

The fleet's one typed wrapper for running external commands — so no tool
scatters raw `os/exec.Command` with ad-hoc output and error handling.

## What

A `Command` built with functional options (`WithArgs`/`WithEnv`/`WithDir`/
`WithStdin`/`WithTimeout`), a `Run(ctx)` that captures a typed `Result`
(stdout/stderr/exit/duration), and `RunJSON[T]` that decodes stdout into a typed
value. A non-zero exit, a failed start, a timeout, and bad JSON each become a
typed, code-carrying error via `errors-go` (`exec_nonzero` / `exec_start` /
`exec_timeout` / `exec_json`).

## Why

Wrapping external CLIs (kubectl, helm, git, …) recurs across
tools. One typed builder means uniform output capture, uniform timeout handling,
and uniform exit-code → `errs.Exit` mapping — never a hand-rolled `exec.Command`
+ buffer + `if err != nil` again.

## Install

```
go get github.com/pleme-io/exec-go
```

## Usage

```go
res, err := execgo.New("kubectl", execgo.WithArgs("get", "ns", "-o", "json"),
    execgo.WithTimeout(10*time.Second)).Run(ctx)
if err != nil { return errs.Exit(err) } // exit_* code → process exit code

list, err := execgo.RunJSON[NamespaceList](ctx,
    execgo.New("kubectl", execgo.WithArgs("get", "ns", "-o", "json")))
```

## Configuration

None — a pure library. Callers that read the command/args/env from config use
`shikumi-go` and pass them via the options.

## Release

Pull-model (Go modules): an annotated `vX.Y.Z` tag is the release; pkg.go.dev
indexes it. See the GSDS module delivery FSM.
