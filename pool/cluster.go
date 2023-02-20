package pool

import "github.com/bird-coder/common/hash"

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
