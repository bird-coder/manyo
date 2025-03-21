/*
 * @Author: yujiajie
 * @Date: 2025-01-23 17:24:20
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 17:30:26
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/accountmanager.go
 * @Description:
 */
package internal

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"
)

var (
	accounts   = make(map[string]Account)
	tlsConfigs = make(map[string]*tls.Config)
	lock       sync.RWMutex
)

type Account struct {
	User string
	Pass string
}

func AddAccount(endpoints []string, user, pass string) {
	lock.Lock()
	defer lock.Unlock()

	accounts[getClusterKey(endpoints)] = Account{
		User: user,
		Pass: pass,
	}
}

func AddTLS(endpoints []string, certFile, certKeyFile, caFile string, insecureSkipVerify bool) error {
	cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
	if err != nil {
		return err
	}

	caData, err := os.ReadFile(caFile)
	if err != nil {
		return err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	lock.Lock()
	defer lock.Unlock()
	tlsConfigs[getClusterKey(endpoints)] = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
		InsecureSkipVerify: insecureSkipVerify,
	}

	return nil
}

func GetAccount(endpoints []string) (Account, bool) {
	lock.RLock()
	defer lock.RUnlock()

	account, ok := accounts[getClusterKey(endpoints)]
	return account, ok
}

func GetTLS(endpoints []string) (*tls.Config, bool) {
	lock.RLock()
	defer lock.RUnlock()

	cfg, ok := tlsConfigs[getClusterKey(endpoints)]
	return cfg, ok
}
