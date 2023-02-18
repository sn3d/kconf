# `kconf`: managing kubeconfigs easily

The `kconf` helps you deal with multiple kubeconfigs easily. It provide 
operations like import and export full contexts (with users and clusters) 
etc.

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
