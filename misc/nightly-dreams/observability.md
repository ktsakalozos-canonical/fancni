# Morpheus — observability

_Generated: 2026-05-13T10:17:49Z_

---

## 1. Enhanced Logging

**Description:**  
Implement structured and contextual logging throughout the internal modules, especially in `fan`, `ipam`, and `cni` packages. Use log levels and include trace IDs for easier debugging in distributed environments.

**Feasibility:**  
High; can leverage existing Go logging libraries and incrementally refactor code for better log coverage.

---

## 2. Metrics Exporter

**Description:**  
Add Prometheus-compatible metrics exporter for runtime statistics such as IP allocations, plugin invocations, and error rates. Integrate with Kubernetes for monitoring via Helm charts.

**Feasibility:**  
High; standard libraries exist and integration is straightforward given current architecture.

---

## 3. Tracing Support

**Description:**  
Integrate distributed tracing (e.g., OpenTelemetry) in key flows like IPAM allocation and network setup. This allows performance and latency analysis across the CNI plugin lifecycle.

**Feasibility:**  
Medium; some refactoring needed but Go has strong support for OpenTelemetry.

---

## 4. Configurable Debug Mode

**Description:**  
Provide a runtime option (via config or environment variable) to enable verbose debug mode, increasing log verbosity and possibly activating additional self-checks.

**Feasibility:**  
High; can be implemented with minimal changes and is user-friendly.

---

## 5. Health Endpoints

**Description:**  
Expose HTTP endpoints for liveness and readiness checks, useful for Kubernetes to monitor the CNI plugin and its sidecar/init containers.

**Feasibility:**  
Medium; requires lightweight HTTP server integration but is common in Kubernetes environments.

---

## 6. Error Reporting Workflow

**Description:**  
Establish a standardized workflow for error reporting, including automatic log collection and optionally telemetry opt-in. Document this in the README and Helm chart.

**Feasibility:**  
Medium; some process and documentation work needed, but technical implementation is manageable.

---

## 7. Customizable Log Sink

**Description:**  
Allow logs to be sent to different sinks (stdout, file, syslog, etc.) configurable via Helm values or environment variables.

**Feasibility:**  
High; Go logging libraries support multiple backends and this is mostly configuration work.

---

## 8. Event Hooks

**Description:**  
Add support for emitting custom events (e.g., IP allocation, subnet exhaustion) that can be consumed by external monitoring or alerting systems.

**Feasibility:**  
Medium; requires an event bus abstraction but fits well with observability goals.

---

## 9. Automated Connectivity Tests

**Description:**  
Extend the existing `deploy/test/connectivity-test.yaml` to automate periodic connectivity checks and report failures to the metrics exporter or logs.

**Feasibility:**  
High; builds on existing test infrastructure, can be triggered on schedule.

---

## 10. Documentation of Observability Features

**Description:**  
Update the README and Helm chart documentation to describe all observability features, including examples of integration with monitoring stacks.

**Feasibility:**  
High; straightforward and increases user adoption and operational clarity.
