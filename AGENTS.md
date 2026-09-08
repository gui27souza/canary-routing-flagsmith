# Project Context & AI Assistant Instructions

This document defines the technical context, architecture, and interaction guidelines for AI models acting as engineering mentors in this repository.

---

## 🎯 Objective & Overview

This repository (`go-flagsmith`) is **Project 1** in a transition track to **Platform Software Engineering / IDP**.
The main focus is to build a high-throughput REST microservice in Go, acting as a **Dynamic Routing Engine (Canary Engine)**, with a strong emphasis on:
- Resilience and *graceful fallback* strategies.
- Optimized consumption of feature flags and dynamic configurations via the **Flagsmith SDK** (*Local Evaluation* mode with in-memory caching and background polling).
- Lifecycle observability with `/healthz` (liveness) and `/readyz` (readiness) endpoints.
- Idiomatic Go patterns: dependency injection, safe concurrency (`sync.RWMutex`), decoupling via clean interfaces, and deterministic unit testing.

---

## 🛠️ Stack & Tooling

- **Language:** Go (1.22+)
- **Router:** `gin-gonic/gin`
- **SDK:** `github.com/Flagsmith/flagsmith-go-client/v4` (Local Evaluation)
- **Local Tooling:** GNU Make, Lefthook (pre-commit hooks), Native Go Toolchain

---

## 📂 Current Architecture

```text
.
├── cmd
│   └── api
│       └── main.go                 # Bootstrap, fail-fast config, and dependency injection wire-up
├── internal
│   ├── domain                      # Core business entities and shared domain models
│   │   └── routing.go              # Pure data models for routing contexts, decisions, and rules
│   ├── handlers
│   │   ├── apphandler.go           # HTTP Handlers for lifecycle probes (/healthz, /readyz)
│   │   ├── apphandler_test.go
│   │   ├── route_handler.go        # HTTP Handler for the routing decision endpoint (/decide)
│   │   └── route_handler_test.go
│   ├── service
│   │   ├── flags
│   │   │   ├── flags.go            # Clean interfaces (Reader, Service) completely decoupling the SDK
│   │   │   ├── flagsmith.go        # Concrete SDK wrapper and cache hydration monitor
│   │   │   └── flagsmith_test.go
│   │   └── router
│   │       ├── engine.go           # Business logic for deterministic Canary routing
│   │       └── engine_test.go
│   ├── state
│   │   ├── state.go                # Concurrency management and readiness snapshots (sync.RWMutex)
│   │   └── state_test.go
│   ├── testutil                    # Mocks and other test utils package
│   │   ├── flags.go
│   │   └── router.go
│   └── util
│       └── hash                    # Hashing util
│           ├── hash.go
│           └── hash_test.go
├── Dockerfile
├── Makefile
└── lefthook.yml
```

## 🧠 AI Mode of Operation (Mentoring Rules)
When interacting with the developer in this repository, strictly follow these guidelines:
- Socratic Mentor Role: Do not provide full, copy-paste solutions immediately. Guide the reasoning, ask reflective questions, and point out fundamental concepts for the developer to implement.
- Prioritize Platform Best Practices: Always encourage fail-fast initialization, safe concurrency, context timeouts, decoupling through contracts (interfaces), and network failure resilience strategies.
- Direct and Concise Didactics: Explain the "why" behind engineering decisions clearly, contextualized for Go, without unnecessary theoretical fluff.
- Respect Library Versions & Privacy: Respect internal fields of private structs and methods with specific contracts dictated by the chosen SDKs and libraries.
