package execgo

import (
	"context"
	"strings"
	"testing"
	"time"

	errs "github.com/pleme-io/errors-go"
)

func TestRunStdout(t *testing.T) {
	res, err := New("sh", WithArgs("-c", "echo hi")).Run(context.Background())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(string(res.Stdout)) != "hi" {
		t.Fatalf("stdout %q", res.Stdout)
	}
	if res.ExitCode != 0 {
		t.Fatalf("exit %d", res.ExitCode)
	}
}

func TestNonZeroExit(t *testing.T) {
	res, err := New("sh", WithArgs("-c", "echo boom >&2; exit 3")).Run(context.Background())
	if err == nil {
		t.Fatal("want error")
	}
	if res.ExitCode != 3 {
		t.Fatalf("exit %d", res.ExitCode)
	}
	if errs.CodeOf(err) != "exec_nonzero" {
		t.Fatalf("code %q", errs.CodeOf(err))
	}
	if !strings.Contains(string(res.Stderr), "boom") {
		t.Fatalf("stderr %q", res.Stderr)
	}
}

func TestRunJSON(t *testing.T) {
	type out struct {
		A int `json:"a"`
	}
	got, err := RunJSON[out](context.Background(), New("sh", WithArgs("-c", "echo '{\"a\":7}'")))
	if err != nil {
		t.Fatalf("runjson: %v", err)
	}
	if got.A != 7 {
		t.Fatalf("a=%d", got.A)
	}
}

func TestTimeout(t *testing.T) {
	_, err := New("sh", WithArgs("-c", "sleep 5"), WithTimeout(50*time.Millisecond)).Run(context.Background())
	if err == nil {
		t.Fatal("want timeout error")
	}
	if errs.CodeOf(err) != "exec_timeout" {
		t.Fatalf("code %q", errs.CodeOf(err))
	}
}

func TestEnv(t *testing.T) {
	res, err := New("sh", WithArgs("-c", "echo $FOO"), WithEnv("FOO=bar")).Run(context.Background())
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if strings.TrimSpace(string(res.Stdout)) != "bar" {
		t.Fatalf("env %q", res.Stdout)
	}
}
