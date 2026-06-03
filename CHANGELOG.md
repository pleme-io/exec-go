# Changelog

All notable changes to this project are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-06-03

### Added
- Initial: typed `Command` (functional options), `Run(ctx)` → typed `Result`
  (stdout/stderr/exit/duration), `RunJSON[T]`, and typed code-carrying errors
  via `errors-go` (`exec_nonzero`/`exec_start`/`exec_timeout`/`exec_json`).
