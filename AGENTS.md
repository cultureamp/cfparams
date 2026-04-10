# AGENTS.md

This file provides guidance to AI agents when working with code in this repository.

## Project Overview

**cfparams** is a Go CLI tool that wrangles AWS CloudFormation parameters. It reads expected parameters from CloudFormation templates and converts them into JSON format suitable for `aws cloudformation create-stack` and `update-stack` commands. The tool handles parameter values from multiple sources (CLI args, YAML files, AWS SSM Parameter Store) and supports CloudFormation's "UsePreviousValue" semantics for stack updates.

## Common Commands

### Building
```bash
# Build binary (embeds version from git tags via ldflags)
make cfparams

# Install to $GOPATH/bin
make install
```

### Testing
```bash
# Run all tests
make test

# Run specific test
go test -run TestDeployScenario

# Run tests with verbose output
go test -v ./...
```

### Release
```bash
# Cross-compile for darwin and linux (note: GoReleaser handles actual releases)
make release
```

### Local Development
```bash
# Run without installing
go run . --template=example/cfn.yaml --parameters=example/parameters.yaml

# Test against example files
go run . --template=example/cfn.yaml --accept-defaults --no-previous Recipient=world
```

## Architecture

### Core Processing Pipeline

The tool follows a pipeline architecture:

1. **Template Parsing** (`template.go`): Extracts parameter specifications from CloudFormation YAML
2. **Parameter Collection** (`parameters.go`): Gathers values from YAML file, CLI args, and optionally AWS Parameter Store
3. **Validation**: Ensures no unexpected parameters are provided
4. **JSON Generation**: Produces CloudFormation-compatible JSON with conditional structure (ParameterValue vs UsePreviousValue)

### Key Components

- **`main.go`**: CLI entry point, flag parsing, orchestration
- **`parameters.go`**: Core parameter processing logic and JSON marshaling
- **`tags.go`**: Stack tags processing (simpler than parameters)
- **`template.go`**: CloudFormation template YAML parsing
- **`parameterstore/store.go`**: AWS SSM Parameter Store client with test faking support

### Custom JSON Marshaling

The `ParameterItem` type implements custom `MarshalJSON()` to produce different JSON structures:
- Parameters with values: `{"ParameterKey": "X", "ParameterValue": "Y"}`
- Parameters using previous values: `{"ParameterKey": "X", "UsePreviousValue": true}`

This conditional marshaling is central to how the tool minimizes verbosity in CloudFormation updates.

### AWS Integration

The `parameterstore` package supports a fake implementation (`parameterstore.Fake()`) to allow testing without AWS credentials. Custom YAML tags (`!ParameterStore`) enable declarative parameter references in YAML files.

### Testing Strategy

Tests in `main_test.go` use:
- Real-world scenario naming (LaunchScenario, DeployScenario)
- Embedded CloudFormation templates as test fixtures
- Helper functions (`mustGetJson`, `mustGetParameterItems`) to reduce boilerplate
- Faked AWS integration for deterministic testing

## Release Process

Releases are managed via GoReleaser and automated through Buildkite:

1. Buildkite pipeline (`.buildkite/pipeline.yaml`) requires manual approval
2. Chinmina issues short-lived GitHub tokens scoped to the release profile (no AWS or long-lived secrets)
3. `git-cliff` calculates the next version from commit history and creates an annotated git tag
4. GoReleaser builds for multiple platforms (darwin/linux/windows on amd64/arm64)
5. Updates Homebrew tap (`cultureamp/homebrew-tap`) via a separate org-scoped token
6. Generates checksums and changelog

Version information is injected at build time via `ldflags` from git tags.

## Dependencies

- **aws-sdk-go-v2**: AWS SSM Parameter Store integration
- **sanathkr/go-yaml**: YAML parsing with custom tag support (enables `!ParameterStore` syntax)
- **testify**: Testing assertions and mocking
