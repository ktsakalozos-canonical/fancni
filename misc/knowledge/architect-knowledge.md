## System Boundaries and Module Responsibilities
- **CNI Plugin**: The core functionality of the project is to implement a Container Network Interface (CNI) plugin. The main entry point is located at `./cmd/fancni/main.go`.
- **IP Address Management (IPAM)**: The `internal/ipam` package manages IP address allocation. Currently, it uses a file-based approach, which should transition to a distributed solution for multi-node environments.
- **Networking Utilities**: The `internal/netutil` package provides helper functions for network operations, such as IP address manipulation and network interface management.
- **iptables Management**: The `internal/iptables` package manages interactions with iptables, with a potential future migration to nftables.
- **Configuration Management**: The `internal/config` package handles configuration settings for the CNI plugin, ensuring that the plugin can be customized based on deployment needs.

## Architectural Decisions
- **Single Executable for Bridge Creation**: The reliance on `fanctl` for bridge creation should be evaluated for potential integration into the Go codebase to reduce external dependencies and improve maintainability.
- **File-based IPAM**: The current implementation of IPAM is file-based. A distributed IPAM solution is essential for supporting multi-node environments and should be prioritized.

## Dependency Graph
- **External Dependencies**: The project relies on several indirect dependencies:
  - `github.com/coreos/go-iptables` (for iptables management)
  - `github.com/vishvananda/netlink` (for network link management)
  - `github.com/vishvananda/netns` (for network namespace management)
  - `golang.org/x/sys` (for system calls)
- **Internal Dependencies**: The `internal/cni` package should be decoupled from other internal packages to improve maintainability and testability. Consider creating interfaces for interactions between packages.

## Incomplete/In-progress Work
- **IPAM Development**: Transition from file-based IPAM to a distributed solution. Create a roadmap for this transition, including milestones and resource allocation.
- **Iptables Migration**: Develop a migration plan to support nftables, including timelines and resource allocation. Evaluate the impact on existing functionality.
- **Observability Enhancements**: Implement metrics and tracing capabilities. Review `misc/nightly-dreams/observability.md` for ideas and prioritize implementation.
- **Testing Framework Expansion**: Move E2E tests from shell scripts (`tests/e2e/test-e2e.sh`) to Go-based tests, focusing on critical paths first to enhance reliability and maintainability.
- **Helm Chart Refinement**: Improve the Helm chart located in `deploy/helm/fancni/` with advanced templating and validation features to ensure proper deployment configurations.

## Areas of Technical Debt
- **IPAM Scalability**: Initiate a project to design a scalable IPAM solution. Engage with the community for best practices and potential collaboration.
- **Error Recovery Mechanisms**: Develop strategies for handling partial failures during CNI operations. Implement retry logic and fallback mechanisms.
- **Containerization of Dependencies**: Investigate containerizing the `fanctl` binary to eliminate PATH dependencies and simplify deployment.
- **Packet Filtering Abstraction**: Plan for an abstraction layer that supports both iptables and nftables to future-proof the networking capabilities.
- **Logging Improvements**: Implement log rotation and configurable log levels. Consider adopting structured logging for better analysis and monitoring.
- **Dynamic Configuration Support**: Explore options for dynamic configuration reloads to minimize downtime during configuration changes.
- **Documentation Enhancements**: Improve user-facing documentation, focusing on `README.md` and `ARCHITECTURE.md`, to provide clear guidance for users and developers.

## General Recommendations
- **Code Reviews**: Establish a routine for code reviews to ensure adherence to architectural decisions and coding standards.
- **Community Engagement**: Engage with the Kubernetes community for feedback and contributions, leveraging community knowledge and best practices.
- **Continuous Integration**: Regularly review and update CI workflows in `.github/workflows/` to ensure robustness and alignment with project goals.
