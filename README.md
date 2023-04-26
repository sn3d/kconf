# `kconf`: kubeconfigs easily

[![Go Report Card](https://goreportcard.com/badge/github.com/sn3d/kconf)](https://goreportcard.com/report/github.com/sn3d/kconf)
[![codebeat badge](https://codebeat.co/badges/4d9db5e8-918f-4561-a9de-5c27c1f509ad)](https://codebeat.co/projects/github-com-sn3d-kconf-main)

The `kconf` helps you with kubeconfigs.

One of the annoying operations with `KUBECONFIG` for me is merging a new cluster
into existing file. I don't want to modify my `KUBECONFIG` everytime, when I
need to add a new cluster. I was tired if manual merging of YAMLs.

I wanted something simple. I want add a new cluster into my existing
`KUBECONFIG` file quickly. Ideally from various sources. Ideally something, I
can use UNIX piping, or I can copy&paste new cluster context into.

I wrote this tool for myself, to relive pains with context manipulation. All 
you need is:

```shell
brew install sn3d/tap/kconf
```

## Import and Export

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

## Basic manipulation

I'm using `kconf` more and more for manipulation with contexts in kubeconfig.
I've implemented few more simple but useful functionalities like `mv`, `rm` 
and `ls`. These subcommands are well known in Unix world and they're manipulating
with contexts same was as Unix commands with files.


The subcommand `kconf ls` will print you all contexts in your `KUBECONFIG` 
file. Also it supports `-l` flag which print all contexts in long listed 
format.

```shell
$ kconf ls -l                                                                                                                                     
CONTEXT          CLUSTER          USER             NAMESPACE
kind-cluster1    kind-cluster1    kind-cluster1    mytest
rancher-desktop  rancher-desktop  rancher-desktop
blue             blue-cluster     john             default
```

With `mv` command, you can rename any context with his user and cluster. 
That means not only context will change the name, but also user and cluster
will be named by context

```shell
$ kconf mv blue cyan
$ kconf ls -l
CONTEXT          CLUSTER          USER             NAMESPACE
kind-cluster1    kind-cluster1    kind-cluster1    mytest
rancher-desktop  rancher-desktop  rancher-desktop
cyan             cyan             cyan             default
```

And last but not least is `rm`. This commant will remove context with 
associated user and cluster.

```shell
$ kconf rm rancher-desktop
$ kconf ls -l                                                                                                                                     
CONTEXT          CLUSTER          USER             NAMESPACE
kind-cluster1    kind-cluster1    kind-cluster1    mytest
cyan             cyan             cyan             default
```

## Splitting kubeconfig

You can split given Kubeconfig into smaller pieces with `kconf split` subdommand. 
Each context will be saved into separated file. The file will follow given prefix, 
and context name as suffix.

```shell
$ kconf split cluster-
cluster-red
cluster-green
cluster-blue
```

You can also set additional postfix and get even better names of files:

```shell
$ kconf split --additional-postfix=.yaml cluster-
cluster-red.yaml
cluster-green.yaml
cluster-blue.yaml
```

You can change the postfix from context name to 2 digit number by using `-d` flag:

```shell
$ kconf split -d --additional-postfix=.yaml cluster-
cluster-01.yaml
cluster-02.yaml
cluster-03.yaml
```

## Installation

Currently, kconf is not available in any package system like Homebrew etc.
But it's simple binary file which can be easily installed by downloading.

### Homebrew (MacOS or Linux)

The preferred method for is to use the Homebrew.

```bash
brew install sn3d/tap/kconf
```

### Curl or wget (MacOS or Linux)

If you don't have brew on your system, you can install `kconf` 
with `curl` or `wget` with one of those one-liners:

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
