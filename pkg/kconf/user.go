package kconf

import apiv1 "k8s.io/client-go/tools/clientcmd/api/v1"

func (c *KubeConfig) getUser(name string) *apiv1.NamedAuthInfo {
	for i := range c.AuthInfos {
		if c.AuthInfos[i].Name == name {
			return &c.AuthInfos[i]
		}
	}
	return nil
}

func (c *KubeConfig) addToUsers(users ...apiv1.NamedAuthInfo) {
	c.AuthInfos = append(c.AuthInfos, users...)
}

func (c *KubeConfig) removeFromUsers(name string) {
	for idx, k := range c.AuthInfos {
		if k.Name == name {
			c.AuthInfos[idx] = c.AuthInfos[len(c.AuthInfos)-1] // copy last element to index
			c.AuthInfos = c.AuthInfos[:len(c.AuthInfos)-1]     // truncate slice
			return
		}
	}
}

func (c *KubeConfig) renameUser(src, dest string) {
	user := c.getUser(src)
	if user != nil {
		user.Name = dest
	}
}
