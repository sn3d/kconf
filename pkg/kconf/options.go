package kconf

type CtxmodOptions struct {
	Cluster string
	User    string
}

type ImportOptions struct {
	As string
}

type ExportOptions struct {
	As string
}

type ClustermodOptions struct {
	ServerURL string
}

type UsermodOptions struct {
	Token string
}
