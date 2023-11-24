package kconf

import (
	"fmt"

	apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

func (kc *KubeConfig) GetUser(name string) *apiv1.NamedAuthInfo {
	for i := range kc.AuthInfos {
		if kc.AuthInfos[i].Name == name {
			return &kc.AuthInfos[i]
		}
	}
	return nil
}

// change user for given context. If context is empty string, then
// the current context will be used
func (c *KubeConfig) ChangeUser(context, user string) error {
	ctx := c.GetContext(context)
	if ctx == nil {
		return fmt.Errorf("no context %s in kubeconfing", context)
	}

	usr := c.GetUser(user)
	if usr == nil {
		return fmt.Errorf("no user %s in kubeconfig", user)
	}

	ctx.Context.AuthInfo = user
	return nil
}

func (kc *KubeConfig) addToUsers(users ...apiv1.NamedAuthInfo) {
	kc.AuthInfos = append(kc.AuthInfos, users...)
}

func (kc *KubeConfig) removeFromUsers(name string) {
	for idx, k := range kc.AuthInfos {
		if k.Name == name {
			kc.AuthInfos[idx] = kc.AuthInfos[len(kc.AuthInfos)-1] // copy last element to index
			kc.AuthInfos = kc.AuthInfos[:len(kc.AuthInfos)-1]     // truncate slice
			return
		}
	}
}

func (kc *KubeConfig) RenameUser(src, dest string) {
	user := kc.GetUser(src)
	if user != nil {
		user.Name = dest
	}
}
