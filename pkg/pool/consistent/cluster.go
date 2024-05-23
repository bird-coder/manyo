/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2022-07-11 13:19:09
 * @LastEditTime: 2024-05-23 17:39:22
 * @LastEditors: yujiajie
 */
package consistent

import "github.com/bird-coder/manyo/pkg/hash"

type PeerPicker interface {
	AddNode(node *hash.Node)
	PickNode(key string) (string, error)
}

type Cluster struct {
	peerPicker PeerPicker
}

func NewCluster() *Cluster {
	cluster := &Cluster{
		peerPicker: hash.NewConsistent(),
	}
	return cluster
}
