package redis

import "crypto/tls"

func cloneTLSConfig(cfg *tls.Config) *tls.Config {
	return cfg.Clone()
}
