# `kconf`: kubeconfigs easily

The `kconf` helps you with kubeconfigs.

One of the annoying operations with `KUBECONFIG` for me is merging a new cluster
into existing file. I don't want to modify my `KUBECONFIG` everytime, when I
need to add a new cluster. I was tired if manual merging of YAMLs.

I wanted something simple. I want add a new cluster into my existing
`KUBECONFIG` file quickly. Ideally from various sources. Ideally something, I
can use UNIX piping, or I can copy&paste new cluster context into.

I wrote this tool for myself, to relive pains with context manipulation.

## How to

- import a cluster from existing file to your `KUBECONFIG` file

```shell
cat ./new-cluster.yaml | kconf import
```

- import a new cluster context into your `KUBECONFIG` file

```yaml
kconf import << EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
  ...
users:
- name:
  ...
contexts:
- context:
  ...
EOF
```

- import a base64 decoded context to specific file

```shell
kconf import --base64 --kubeconfig=/path/to/kube.conf << EOF
LSBjb250ZXh0OgogICAgY2x1c3RlcjogIDxjbHV
zdGVyLW5hbWU+CiAgICB1c2VyOiAgPGNsdXN0ZX
ItbmFtZS11c2VyPgogIG5hbWU6ICA8Y2x1c3Rlc
i1uYW1lPg==
EOF
```

- export a full context (with user and cluster) from your `KUBECONFIG` file

```shell
kconf export your-k8s-cluster >> ./your-k8s-cluster.conf
```

- remove a full context (with user and cluster) from your
`KUBECONFIG` file

```shell
kconf rm your-k8s-cluster
```

## Installation

Currently, kconf is not available in any package system like Homebrew etc.
But it's simple binary file which can be easily installed by downloading.

### MacOS or Linux

You can install kconf with curl or get with those one-liners:

```bash
curl -s https://installme.sh/sn3d/kconf | sh
```

or

```bash
wget -q -O - https://installme.sh/sn3d/kconf | sh
```

### Windows

Download the correct binary for your platform from [project's GitHub](https://github.com/sn3d/kconf/releases/).
Uncompress the binary to you PATH.

# Feedback & Bugs

Feedback is more than welcome. Did you found a bug? Is something not behaving as expected? Feature or bug, feel free to create [issue](https://github.com/sn3d/kconf/issues).
