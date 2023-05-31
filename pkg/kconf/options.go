package kconf

type CtxmodOptions struct {
	Cluster string
	User    string
}

type ClustermodOptions struct {
	ServerURL string
}

type UsermodOptions struct {
	Token string
}
