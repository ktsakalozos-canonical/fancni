## System Boundaries and Module Responsibilities
- **CNI Plugin**: Implemented in `./cmd/fancni/main.go`, responsible for managing network interfaces and integrating with Kubernetes.
- **IP Address Management (IPAM)**: Located in `./internal/ipam`, currently using a file-based approach. Transition to a distributed solution is essential for scalability in multi-node environments.
- **Networking Utilities**: The `internal/netutil` package provides essential helper functions for network operations, such as IP address manipulation and network interface management.
- **iptables Management**: Managed in `internal/iptables`, with plans to migrate to nftables for improved functionality and performance.
- **Configuration Management**: Handled by `internal/config`, allowing customizable settings for the CNI plugin based on deployment requirements.
- **Testing**: E2E tests are located in `tests/e2e/test-e2e.sh` and need to be transitioned to Go-based tests for better maintainability.

## Architectural Decisions
- **Single Executable for Bridge Creation**: Consider integrating `fanctl` into the Go codebase to reduce external dependencies and enhance maintainability.
- **File-based IPAM**: Prioritize transitioning to a distributed IPAM solution for improved scalability in multi-node environments.
- **Helm for Deployment**: Utilize Helm charts located in `deploy/helm/fancni/` for managing configurations and deployments.
- **Testing Framework**: The current E2E testing framework relies on shell scripts, which should be replaced with Go-based tests to improve reliability.

## Dependency Graph
- **External Dependencies**:
  - `github.com/coreos/go-iptables` (for iptables management)
  - `github.com/vishvananda/netlink` (for network link management)
  - `github.com/vishvananda/netns` (for network namespace management)
  - `golang.org/x/sys` (for system calls)

- **Internal Dependencies**:
  - `internal/cni` should be decoupled from other internal packages to enhance maintainability and testability. Define interfaces for interactions between packages.
  - `internal/ipam` is tightly coupled with `internal/netutil` for network operations.

## Incomplete/In-progress Work
- **IPAM Development**: Transition from file-based IPAM to a distributed solution. Create a roadmap with milestones and resource allocation.
- **Iptables Migration**: Develop a migration plan to support nftables, including timelines and resource allocation. Assess the impact on existing functionality.
- **Observability Enhancements**: Implement metrics and tracing capabilities. Review `misc/nightly-dreams/observability.md` for ideas and prioritize implementation.
- **Testing Framework Expansion**: Transition E2E tests from shell scripts (`tests/e2e/test-e2e.sh`) to Go-based tests, focusing on critical paths first to improve reliability and maintainability.
- **Helm Chart Refinement**: Enhance the Helm chart in `deploy/helm/fancni/` with advanced templating and validation features to ensure proper deployment configurations.

## Areas of Technical Debt
- **IPAM Scalability**: Initiate a project to design a scalable IPAM solution. Engage with the community for best practices and potential collaboration.
- **Error Recovery Mechanisms**: Develop strategies for handling partial failures during CNI operations. Implement retry logic and fallback mechanisms.
- **Containerization of Dependencies**: Investigate containerizing the `fanctl` binary to eliminate PATH dependencies and simplify deployment.
- **Packet Filtering Abstraction**: Plan for an abstraction layer that supports both iptables and nftables to future-proof networking capabilities.
- **Logging Improvements**: Implement log rotation and configurable log levels. Consider adopting structured logging for better analysis and monitoring.
- **Dynamic Configuration Support**: Explore options for dynamic configuration reloads to minimize downtime during configuration changes.
- **Documentation Enhancements**: Improve user-facing documentation, focusing on `README.md` and `ARCHITECTURE.md`, to provide clear guidance for users and developers.

## General Recommendations
- **Code Reviews**: Establish a routine for code reviews to ensure adherence to architectural decisions and coding standards.
- **Community Engagement**: Engage with the Kubernetes community for feedback and contributions, leveraging community knowledge and best practices.
- **Continuous Integration**: Regularly review and update CI workflows in `.github/workflows/` to ensure robustness and alignment with project goals.
