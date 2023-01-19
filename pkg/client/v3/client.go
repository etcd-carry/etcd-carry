package v3

import (
	"go.etcd.io/etcd/client/v3"
)

func New(cs ConfigSpec) (*clientv3.Client, error) {
	tlscfg, err := NewTLSConfig(cs.Secure)
	if err != nil {
		return nil, err
	}

	cfg := clientv3.Config{
		Endpoints:            cs.Endpoints,
		DialTimeout:          cs.DialTimeout,
		DialOptions:          cs.DialOptions,
		DialKeepAliveTime:    cs.KeepAliveTime,
		DialKeepAliveTimeout: cs.KeepAliveTimeout,
		TLS:                  tlscfg,
	}

	return clientv3.New(cfg)
}
