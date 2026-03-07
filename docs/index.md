# kubetray

Local K8s, served on a tray.

`kubetray` is a CLI tool to create local Kubernetes development environments with a single command.

## Quick Start

```bash
kubetray start
export KUBECONFIG=~/.kubetray/kubeconfig
kubectl get nodes
```

## Installation

```bash
git clone https://github.com/ugurozkn/kubetray.git
cd kubetray
make install
```

## Commands

```bash
kubetray start
kubetray stop
kubetray clean
```

## Project Links

- Repository: https://github.com/ugurozkn/kubetray
- README: https://github.com/ugurozkn/kubetray/blob/main/README.md
- License: https://github.com/ugurozkn/kubetray/blob/main/LICENSE
