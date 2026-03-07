# `kubetray`

**Local K8s, served on a tray.**

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)]()

> **Early development** — Core cluster management works. More features coming incrementally.

---

## What is `kubetray`?

`kubetray` is a CLI tool that creates local Kubernetes development environments with a single command.
It uses [k3s](https://k3s.io/) via [k3d](https://k3d.io/) to spin up lightweight clusters in Docker.

```bash
kubetray start
```

```
KubeTray
────────

✓ Platform: macOS (Apple Silicon) 26.2 (arm64)
✓ Dependencies: helm, kubectl, docker
✓ Container runtime: colima
✓ Cluster 'kubetray' created

Cluster is ready!

Cluster Details
───────────────
  Cluster name      kubetray
  Kubernetes        k3s (via k3d)
  Container runtime colima
  Resources         2 CPUs, 2G RAM
  Kubeconfig        ~/.kubetray/kubeconfig
  Setup time        45s

Quick Start
───────────
  export KUBECONFIG=~/.kubetray/kubeconfig
  kubectl get nodes
  kubectl get pods -A
```

---

## Usage

```bash
kubetray start                       # Start cluster (2 CPUs, 2G RAM)
kubetray start --cpus 4 --memory 8G  # Start with more resources
kubetray stop                        # Stop cluster (preserves data)
kubetray clean                       # Delete cluster completely
kubetray clean --force               # Delete without confirmation
```

---

## Installation

### Prerequisites

You need a container runtime and a couple of CLI tools:

- **Container runtime** (one of):
  [Colima](https://github.com/abiosoft/colima),
  [Docker Desktop](https://www.docker.com/products/docker-desktop/),
  or [OrbStack](https://orbstack.dev/)
- **helm** — `brew install helm`
- **kubectl** — `brew install kubectl`

### From source

```bash
git clone https://github.com/ugurozkn/kubetray.git
cd kubetray
make install
```

---

## Roadmap

`kubetray` is being built incrementally. Each feature ships as a separate PR.

- [x] **Cluster lifecycle** — `start`, `stop`, `clean`
- [ ] **Version command** — version, build info
- [ ] **Status command** — cluster health, node info, pod status
- [ ] **Component system** — install Prometheus, Grafana, Loki, ArgoCD via Helm
- [ ] **Log streaming** — `kubetray logs <component>`
- [ ] **Profiles** — pre-configured stacks (e.g. `kubetray start -p fullstack`)
- [ ] **Ingress** — `*.dev.local` DNS routing with Traefik
- [ ] **Project config** — per-project `.kubetray/` directory with custom Helm values
- [ ] **Docker Compose deploy** — `kubetray deploy -f docker-compose.yml`
- [ ] **App management** — scale, autoscale, resource limits, restart
- [ ] **CI/CD & releases** — GoReleaser, Homebrew tap, GitHub Actions

---

## Supported Platforms

| Platform | Architecture | Status |
|----------|-------------|--------|
| macOS | Apple Silicon (M1/M2/M3/M4) | **Supported** |
| macOS | Intel | **Supported** |
| Linux | amd64 | Planned |
| Linux | arm64 | Planned |

---

## Contributing

Contributions are welcome. The project is in early stages — there's plenty of room to help shape it.

```bash
git clone https://github.com/ugurozkn/kubetray.git
cd kubetray
make build    # Build binary
make test     # Run tests
```

---

## License

[MIT](LICENSE)
