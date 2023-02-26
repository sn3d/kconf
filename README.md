# `kconf`: kubeconfigs easily

The `kconf` helps you with kubeconfigs. 

One of the annoying operations with KUBECONFIG for me is merging a new cluster. 
I don't want to modify my KUBECONFIG everytime, when I need to add a new cluster. 
I was tired if manual merging of YAMLs.

I wanted something simple. I want add a new cluster into my existing 
KUBECONFIG file quickly. Ideally from various sources. Ideally something, I 
can use UNIX piping, or I can copy&paste new cluster context into.

I wrote this tool for myself, to releave pains with context manipulation.

## How to

- how to import a new context into your existing `KUBECONFIG` 
file

```
kconf import << EOF
apiVersion: v1
kind: Config
...
EOF
```

- how to import base64 decoded context into your existing 
kubeconfig

```
kconf import --base64 --kubeconfig=/path/to/kube.conf << EOF
LSBjb250ZXh0OgogICAgY2x1c3RlcjogIDxjbHV
zdGVyLW5hbWU+CiAgICB1c2VyOiAgPGNsdXN0ZX
ItbmFtZS11c2VyPgogIG5hbWU6ICA8Y2x1c3Rlc
i1uYW1lPg==
EOF
```

- how to export full context from your `KUBECONFIG` file

```
kconf export your-k8s-cluster >> ./your-k8s-cluster.conf
```

- how to remove full context (with user and cluster) from your
`KUBECONFIG` file

```
kconf rm your-k8s-cluster
```
