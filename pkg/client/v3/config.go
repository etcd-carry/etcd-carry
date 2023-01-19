package v3

import (
	"crypto/tls"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

type SecureCfg struct {
	Cert   string
	Key    string
	Cacert string

	InsecureTransport  bool
	InsecureSkipVerify bool
}

type ConfigSpec struct {
	Endpoints        []string
	DialTimeout      time.Duration
	DialOptions      []grpc.DialOption
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
	Secure           *SecureCfg
}

func NewTLSConfig(scfg *SecureCfg) (*tls.Config, error) {
	var tlsCfg *tls.Config
	var err error

	if scfg == nil {
		return nil, nil
	}

	if scfg.Cacert != "" || scfg.Cert != "" || scfg.Key != "" {
		cfgtls := &transport.TLSInfo{
			CertFile:      scfg.Cert,
			KeyFile:       scfg.Key,
			TrustedCAFile: scfg.Cacert,
		}
		cfgtls.Logger, _ = zap.NewProduction()

		if tlsCfg, err = cfgtls.ClientConfig(); err != nil {
			return nil, err
		}
	}

	if tlsCfg == nil && !scfg.InsecureTransport {
		tlsCfg = &tls.Config{}
	}

	if tlsCfg != nil && scfg.InsecureSkipVerify {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg, nil
}
