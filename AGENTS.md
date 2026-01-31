# AGENTS.md

This file provides guidance to Qoder (qoder.com) when working with code in this repository.

## Project Information

**Framework**: Ursa (github.com/loveuer/ursa)  
**Version**: 1.0.0  
**Description**: A fast, simple, and production-ready web framework for Go

## Development Commands

- **Build**: `go build`
- **Test all**: `go test ./...`
- **Run single test**: `go test -v -run <TestName>`
- **Lint**: `golangci-lint run` (standard Go linting, though no config file is present)
- **CLI Tool**: The management tool is located in `ursatool/ursactl`. To run it: `go run ursatool/ursactl/main.go`

## Code Architecture

### Core Framework
- **App (`app.go`)**: The central engine. It manages the `http.Server` lifecycle, configuration, and the radix tree routers for different HTTP methods. It uses a `sync.Pool` for `Ctx` objects to minimize memory allocations.
- **Ctx (`ctx.go`)**: The request/response context. It encapsulates `http.Request` and `http.ResponseWriter`, providing a rich API for parameter extraction, body parsing (`BodyParser`, `QueryParser`), and response sending (`JSON`, `HTML`, `SSEvent`).
- **RouterGroup (`routergroup.go`)**: Implements the `IRouter` interface. It handles route registration and middleware composition using a simple slice-based handler chain.
- **Radix Tree (`tree.go`)**: Efficient route matching using a radix tree implementation.

### Request Pipeline
1. `App.ServeHTTP` is the entry point for all requests.
2. A `Ctx` is retrieved from the pool and reset with the current request/writer.
3. `handleHTTPRequest` finds the matching route in the radix tree for the current method.
4. If found, the handler chain (middleware + final handler) is assigned to `Ctx.handlers`.
5. `Ctx.Next()` is called to execute the chain sequentially.

### Built-in Middlewares
- **NewCORS/NewCORSWithConfig**: CORS (Cross-Origin Resource Sharing) support
- **NewSecure/NewSecureWithConfig**: Security headers (XSS, Content-Type, Frame Options, HSTS, CSP)
- **NewRequestID/NewRequestIDWithConfig**: Request tracking with UUID generation
- **NewLogger**: HTTP request logging
- **NewRecover**: Panic recovery with stack trace

### Internal Packages
- **internal/schema**: Handles the binding and decoding of request data (queries, forms, JSON) into Go structs.
- **internal/bytesconv**: Optimized zero-allocation conversions between strings and byte slices.
- **internal/sse**: Server-Sent Events implementation.

### CLI Tool (ursactl)
Located in `ursatool/ursactl`, it uses `spf13/cobra` for command-line interaction. It provides utilities for project management and scaffolding.
