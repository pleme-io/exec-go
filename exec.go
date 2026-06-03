// Package execgo is the fleet's one typed wrapper for running external
// commands — so no tool scatters raw os/exec.Command with ad-hoc output
// handling. A Command is built with functional options, Run captures a typed
// Result (stdout/stderr/exit/duration), and a non-zero exit becomes a typed,
// severity- and exit-code-carrying error via errors-go.
package execgo

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"time"

	errs "github.com/pleme-io/errors-go"
)

// Command is a typed, options-built external command.
type Command struct {
	Name    string
	Args    []string
	Env     []string // appended to os.Environ() unless EnvOnly is set
	EnvOnly bool
	Dir     string
	Stdin   []byte
	Timeout time.Duration // 0 = no timeout (ctx still governs)
}

// Option configures a Command.
type Option func(*Command)

// WithArgs sets the command arguments.
func WithArgs(args ...string) Option { return func(c *Command) { c.Args = args } }

// WithEnv appends KEY=VALUE entries to the inherited environment.
func WithEnv(kv ...string) Option { return func(c *Command) { c.Env = append(c.Env, kv...) } }

// WithEnvOnly replaces (does not append to) the environment.
func WithEnvOnly(kv ...string) Option {
	return func(c *Command) { c.Env = kv; c.EnvOnly = true }
}

// WithDir sets the working directory.
func WithDir(dir string) Option { return func(c *Command) { c.Dir = dir } }

// WithStdin supplies stdin bytes.
func WithStdin(b []byte) Option { return func(c *Command) { c.Stdin = b } }

// WithTimeout bounds the run (in addition to the caller's context).
func WithTimeout(d time.Duration) Option { return func(c *Command) { c.Timeout = d } }

// New builds a Command.
func New(name string, opts ...Option) *Command {
	c := &Command{Name: name}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Result is the typed outcome of a run.
type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
}

// Run executes the command, capturing output. A non-zero exit yields a typed
// *errs error (Severity Error, code "exec_nonzero", carrying the Result via
// CodeOf-friendly context); ctx cancellation / timeout yields code "exec_timeout".
func (c *Command) Run(ctx context.Context) (Result, error) {
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, c.Name, c.Args...)
	if c.EnvOnly {
		cmd.Env = c.Env
	} else if len(c.Env) > 0 {
		cmd.Env = append(os.Environ(), c.Env...)
	}
	cmd.Dir = c.Dir
	if c.Stdin != nil {
		cmd.Stdin = bytes.NewReader(c.Stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	res := Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), Duration: time.Since(start)}

	if ctxErr := ctx.Err(); ctxErr != nil {
		return res, errs.Wrap(ctxErr, "execgo: "+c.Name+" timed out / cancelled", errs.WithCode("exec_timeout"))
	}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			res.ExitCode = ee.ExitCode()
			return res, errs.New("execgo: "+c.Name+" exited "+itoa(res.ExitCode)+": "+truncate(res.Stderr), errs.WithCode("exec_nonzero"))
		}
		return res, errs.Wrap(err, "execgo: "+c.Name+" failed to start", errs.WithCode("exec_start"))
	}
	return res, nil
}

// RunJSON runs the command and decodes its stdout JSON into T.
func RunJSON[T any](ctx context.Context, c *Command) (T, error) {
	var out T
	res, err := c.Run(ctx)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(res.Stdout, &out); err != nil {
		return out, errs.Wrap(err, "execgo: "+c.Name+" stdout is not valid JSON", errs.WithCode("exec_json"))
	}
	return out, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

func truncate(b []byte) string {
	const max = 200
	if len(b) > max {
		return string(b[:max]) + "…"
	}
	return string(b)
}
