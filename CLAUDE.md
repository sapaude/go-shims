# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-shims is a Go utility library providing common helper functions for Go projects. The repository is organized as a monorepo with three main independent modules:

1. **shim** - Core utility functions (string processing, number conversion, file operations, path handling, crypto, LLM helpers, etc.)
2. **x/log** - Logging library based on logrus
3. **infra/kafkax** - Kafka integration utilities with consumer group support

Each module has its own `go.mod` file and can be imported independently.

## Repository Structure

The repository uses a **single version strategy** where all modules share the same version number tracked in the `version` file at the root. When any module changes, the entire repository version is bumped.

### Module Paths

- `github.com/sapaude/go-shims/shim` - Core utility module
- `github.com/sapaude/go-shims/x/log` - Logging module
- `github.com/sapaude/go-shims/infra/kafkax` - Kafka infrastructure module

### Key Directories

- `shim/` - Core utility functions with comprehensive test coverage
- `x/log/` - Logging library with hooks and context support
- `infra/kafkax/` - Kafka consumer group implementation with example usage in `infra/kafkax/example/`

## Build and Test Commands

### Running Tests

```bash
# Test a specific module
cd shim && go test ./...
cd x/log && go test ./...
cd infra/kafkax && go test ./...

# Run a single test
cd shim && go test -run TestSpecificFunction

# Run tests with verbose output
cd shim && go test -v ./...
```

### Version Management

The repository uses a Makefile for version management:

```bash
# Bump patch version, commit, and create a git tag
make tag

# Push tags and main branch to remote
make push

# Bump version and push in one command
make update
```

The version file tracks the current version (e.g., `v0.2.8`). Running `make tag` automatically increments the patch version.

## Architecture Notes

### Shim Module

The shim module provides utility functions organized by domain:

- **strings.go** - String manipulation, JSON conversion, parsing, hashing, stream generation
- **number.go** - Number conversion utilities
- **files.go** - File operation helpers
- **path.go** - Path handling utilities
- **crypto.go** - Cryptographic functions
- **money.go** - Money/decimal handling using shopspring/decimal
- **random.go** - Random generation utilities
- **network.go** - Network-related helpers
- **llm.go** - LLM (Large Language Model) integration helpers
- **data_copy.go** - Data copying utilities
- **order.go** - Ordering/sorting helpers
- **time.go** - Time-related utilities

All functions are well-tested with corresponding `*_test.go` files.

### Kafkax Module

The kafkax module wraps IBM/sarama for simplified Kafka consumer group usage. Key files:

- `kafka_config.go` - Kafka configuration setup
- `consumer_group.go` - Consumer group implementation
- `example/` - Working examples showing how to use the library

### X/Log Module

Structured logging built on logrus with:

- `log.go` - Main logging interface
- `logger.go` - Logger implementation
- `config.go` - Configuration options
- `context.go` - Context integration
- `hook.go` - Custom hooks support

## Development Guidelines

- Each module is independent with its own go.mod
- Changes to any module trigger a version bump for the entire repository
- All new functions should include comprehensive tests
- The repository uses Go 1.24.1
- Dependencies are managed per-module but version bumps are synchronized

## Testing Kafka Integration

See examples in `infra/kafkax/example/` for consumer implementation patterns. The module requires external Kafka setup for integration testing.
