/*
 * @Author: yujiajie
 * @Date: 2025-01-15 15:22:39
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-15 15:33:10
 * @FilePath: /Go-Base/pkg/discov/etcd/config.go
 * @Description:
 */
package etcd

import "errors"

var (
	// errEmptyEtcdHosts indicates that etcd hosts are empty.
	errEmptyEtcdHosts = errors.New("empty etcd hosts")
	// errEmptyEtcdKey indicates that etcd key is empty.
	errEmptyEtcdKey = errors.New("empty etcd key")
)

type EtcdConf struct {
	Hosts              []string
	Key                string
	ID                 int64  `json:",optional"`
	User               string `json:",optional"`
	Pass               string `json:",optional"`
	CertFile           string `json:",optional"`
	CertKeyFile        string `json:",optional=CertFile"`
	CACertFile         string `json:",optional=CertFile"`
	InsecureSkipVerify bool   `json:",optional"`
}

// HasAccount returns if account provided.
func (c EtcdConf) HasAccount() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}

// HasID returns if ID provided.
func (c EtcdConf) HasID() bool {
	return c.ID > 0
}

// HasTLS returns if TLS CertFile/CertKeyFile/CACertFile are provided.
func (c EtcdConf) HasTLS() bool {
	return len(c.CertFile) > 0 && len(c.CertKeyFile) > 0 && len(c.CACertFile) > 0
}

// Validate validates c.
func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errEmptyEtcdHosts
	} else if len(c.Key) == 0 {
		return errEmptyEtcdKey
	} else {
		return nil
	}
}
