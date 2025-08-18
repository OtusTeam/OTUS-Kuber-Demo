# k3d Cluster Configurations

This folder contains example YAML configurations for creating local Kubernetes clusters using [k3d](https://github.com/k3d-io/k3d).

## Usage

1. Install k3d:
   ```sh
   curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
   # or use Homebrew
   brew install k3d
   ```

2. Create a single-node cluster:
   ```sh
   k3d cluster create --config single-node.yaml
   ```

3. List clusters:
   ```sh
   k3d cluster list
   ```

4. Delete cluster:
   ```sh
   k3d cluster delete otus-demo
   ```

## Files
- `single-node.yaml`: Minimal single-node cluster with custom node labels and reserved resources.

## References
- [k3d documentation](https://k3d.io/)
- [k3d config file reference](https://k3d.io/v5.6.0/usage/configfile/)

---
For more advanced scenarios (multi-node, ingress, volumes, etc.), see k3d documentation or request additional examples.
