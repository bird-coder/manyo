/*
 * @Author: yujiajie
 * @Date: 2025-01-15 15:50:32
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 16:59:47
 * @FilePath: /manyo/pkg/discov/etcd/publisher.go
 * @Description:
 */
package etcd

import (
	"time"

	"github.com/bird-coder/manyo/pkg/discov/etcd/internal"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/syncx"
	"github.com/bird-coder/manyo/pkg/threading"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	PubOption func(client *Publisher)

	Publisher struct {
		endpoints  []string
		key        string
		fullKey    string
		id         int64
		value      string
		lease      clientv3.LeaseID
		quit       *syncx.DoneChan
		pauseChan  chan struct{}
		resumeChan chan struct{}
	}
)

func NewPublisher(endpoints []string, key, value string, opts ...PubOption) *Publisher {
	publisher := &Publisher{
		endpoints:  endpoints,
		key:        key,
		value:      value,
		quit:       syncx.NewDoneChan(),
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

func (p *Publisher) KeepAlive() error {
	cli, err := p.doRegister()
	if err != nil {
		return err
	}
	return p.keepAliveAsync(cli)
}

func (p *Publisher) Parse() {
	p.pauseChan <- struct{}{}
}

func (p *Publisher) Resume() {
	p.resumeChan <- struct{}{}
}

func (p *Publisher) Stop() {
	p.quit.Close()
}

func (p *Publisher) doKeepAlive() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-p.quit.Done():
			return nil
		default:
			cli, err := p.doRegister()
			if err != nil {
				logger.Error("etcd publisher doRegister: %s", err.Error())
				break
			}
			if err := p.keepAliveAsync(cli); err != nil {
				logger.Error("etcd publisher keepAliveAsync: %s", err.Error())
				break
			}
			return nil
		}
	}
	return nil
}

func (p *Publisher) doRegister() (internal.EtcdClient, error) {
	cli, err := internal.GetRegistry().GetConn(p.endpoints)
	if err != nil {
		return nil, err
	}
	p.lease, err = p.register(cli)
	return cli, err
}

func (p *Publisher) keepAliveAsync(cli internal.EtcdClient) error {
	ch, err := cli.KeepAlive(cli.Ctx(), p.lease)
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					p.revoke(cli)
					if err := p.doKeepAlive(); err != nil {
						logger.Error("etcd publisher KeepAlive: %s", err.Error())
					}
					return
				}
			case <-p.pauseChan:
				logger.Info("paused etcd renew, key: %s, value: %s", p.key, p.value)
				p.revoke(cli)
				select {
				case <-p.resumeChan:
					if err := p.doKeepAlive(); err != nil {
						logger.Error("etcd publisher KeepAlive: %s", err.Error())
					}
					return
				case <-p.quit.Done():
					return
				}
			case <-p.quit.Done():
				p.revoke(cli)
				return
			}
		}
	})

	return nil
}

func (p *Publisher) register(client internal.EtcdClient) (clientv3.LeaseID, error) {
	resp, err := client.Grant(client.Ctx(), TimeToLive)
	if err != nil {
		return clientv3.NoLease, err
	}

	lease := resp.ID
	if p.id > 0 {
		p.fullKey = makeEtcdKey(p.key, p.id)
	} else {
		p.fullKey = makeEtcdKey(p.key, int64(lease))
	}
	_, err = client.Put(client.Ctx(), p.fullKey, p.value, clientv3.WithLease(lease))

	return lease, err
}

func (p *Publisher) revoke(cli internal.EtcdClient) {
	if _, err := cli.Revoke(cli.Ctx(), p.lease); err != nil {
		logger.Error("etcd publisher revoke: %s", err.Error())
	}
}

func WithId(id int64) PubOption {
	return func(client *Publisher) {
		client.id = id
	}
}

func WithPubEtcdAccount(user, pass string) PubOption {
	return func(client *Publisher) {
		RegisterAccount(client.endpoints, user, pass)
	}
}

func WithPubEtcdTLS(certFile, certKeyFile, caFile string, insecureSkipVerify bool) PubOption {
	return func(client *Publisher) {
		err := RegisterTLS(client.endpoints, certFile, certKeyFile, caFile, insecureSkipVerify)
		logger.Error("%v", err)
	}
}
