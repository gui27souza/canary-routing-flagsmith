# 📦 Resilient Microservice with Flagsmith & Go

## 🎯 Objective & Use Cases

Develop a REST API in Go acting as a **Dynamic Routing Engine (Canary Engine)**, consuming Feature Flags and Remote Config in a non-blocking, resilient, and secure manner for high-throughput environments.

The microservice must be able to make real-time decisions about traffic destinations (e.g., `v1-legacy` vs `v2-new`) based on user context and platform-defined rules, without the need for additional deployments.

---

## 🛠️ Technical Stack

- **Language:** Go (1.22+)
- **Router:** `gin-gonic/gin`
- **SDK:** Flagsmith Go SDK (with local evaluation cache and background polling)
- **Container:** Docker (Multi-stage build with `gcr.io/distroless/static-debian12:nonroot`)
- **Performance & Automation:** `k6`, GNU Make, Lefthook

---

## 📌 Main Deliverables & Requirements — v1 (Completed)

- [x] **SDK Integration:** Configure Flagsmith client with in-memory caching and background polling intervals (Local Evaluation mode).
- [x] **Safe Fallback Pattern:** Ensure network failures or config sync errors do not panic or drop traffic (fallback to conservative defaults).
- [x] **Observability & Health Checks:** `/healthz` (liveness) and `/readyz` (readiness) endpoints reflecting application integrity and cache hydration state.
- [x] **Dynamic Canary Engine:** Implement deterministic percentage-based canary routing via FNV-1a normalized hashing (`0-100%`) on `/decide`.
- [x] **Automated Testing:** Table-driven unit tests utilizing clean interfaces (`flags.Reader`, `DecisionMaker`) and isolated mocks (`internal/testutil`).
- [x] **Graceful Shutdown:** Handle `SIGINT` / `SIGTERM` signals with connection draining (`http.Server.Shutdown`) and background routine cancellation.
- [x] **Production Containerization:** Minimal, secure multi-stage Dockerfile running on Distroless as a non-root user (`nonroot:65532`).
- [x] **Performance Benchmarking:** Automated load testing suite via `k6` validating >2.1k RPS and sub-6ms p95 latency.

---

## 🚀 Roadmap & Platform Enhancements — v2

Planned enhancements to elevate the microservice further into an enterprise-grade Internal Developer Platform (IDP) component:

- [ ] **CI/CD Automation (GitHub Actions):**
  - Implement a pipeline checking formatting, vetting, testing with race detector (`-race`), and building the distroless image on every push/PR.
  - Automated security scanning (Trivy / Grype) for container vulnerabilities.
- [ ] **Telemetry & Metrics (Prometheus / Grafana):**
  - Expose a `/metrics` endpoint with Prometheus metrics.
  - Track metrics such as: routing decision counts partitioned by target (`canary_decisions_total{target="v1|v2"}`), flag cache sync latency, and HTTP request durations.
- [ ] **Distributed Tracing (OpenTelemetry):**
  - Inject trace propagation headers (`traceparent`) into routing responses and logs to correlate gateway requests with downstream services.
- [ ] **Kubernetes Deployment Artifacts (GitOps-Ready):**
  - Create Helm charts or Kustomize manifests defining `Deployment`, `Service`, `HPA` (Horizontal Pod Autoscaler), and `PodDisruptionBudget`.
  - Configure native Kubernetes liveness and readiness probe specs pointing to `/healthz` and `/readyz`.
- [ ] **Rate Limiting & Token Bucket:**
  - Add in-memory or Redis-backed rate limiting middleware to protect the routing engine against malicious spikes.
