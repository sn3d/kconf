package usr

import apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

type UserItem struct {
	User *apiv1.NamedAuthInfo
}

func (i UserItem) Title() string       { return i.User.Name }
func (i UserItem) FilterValue() string { return i.User.Name }
