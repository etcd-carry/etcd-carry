package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"time"
)

type TransportConfig struct {
	Insecure           bool
	InsecureSkipVerify bool
	ServerList         []string
	CACertFile         string
	CertFile           string
	KeyFile            string
}

type TransportOptions struct {
	SourceTransport  TransportConfig
	DestTransport    TransportConfig
	DialTimeout      time.Duration
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration
}

const (
	defaultDialTimeout      = 2 * time.Second
	defaultKeepAliveTime    = 2 * time.Second
	defaultKeepAliveTimeOut = 6 * time.Second
)

func NewTransportOptions() *TransportOptions {
	return &TransportOptions{
		SourceTransport: TransportConfig{
			Insecure:           true,
			InsecureSkipVerify: false,
			CACertFile:         "/etc/kubernetes/source/etcd/ca.crt",
			CertFile:           "/etc/kubernetes/source/etcd/server.crt",
			KeyFile:            "/etc/kubernetes/source/etcd/server.key",
		},
		DestTransport: TransportConfig{
			Insecure:           true,
			InsecureSkipVerify: false,
			CACertFile:         "/etc/kubernetes/dest/etcd/ca.crt",
			CertFile:           "/etc/kubernetes/dest/etcd/server.crt",
			KeyFile:            "/etc/kubernetes/dest/etcd/server.key",
		},
		DialTimeout:      defaultDialTimeout,
		KeepAliveTime:    defaultKeepAliveTime,
		KeepAliveTimeout: defaultKeepAliveTimeOut,
	}
}

func (s *TransportOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.SourceTransport.Insecure, "source-insecure-transport", s.SourceTransport.Insecure, "disable transport security for client connections")
	fs.BoolVar(&s.SourceTransport.InsecureSkipVerify, "source-insecure-skip-tls-verify", s.SourceTransport.InsecureSkipVerify, "skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)")
	fs.StringSliceVar(&s.SourceTransport.ServerList, "source-endpoints", s.SourceTransport.ServerList, "List of etcd servers to connect with (scheme://ip:port), comma separated")
	fs.StringVar(&s.SourceTransport.CACertFile, "source-cacert", s.SourceTransport.CACertFile, "verify certificates of TLS-enabled secure servers using this CA bundle")
	fs.StringVar(&s.SourceTransport.CertFile, "source-cert", s.SourceTransport.CertFile, "identify secure client using this TLS certificate file")
	fs.StringVar(&s.SourceTransport.KeyFile, "source-key", s.SourceTransport.KeyFile, "identify secure client using this TLS key file")

	fs.BoolVar(&s.DestTransport.Insecure, "dest-insecure-transport", s.DestTransport.Insecure, "Disable transport security for client connections for the destination cluster")
	fs.BoolVar(&s.DestTransport.InsecureSkipVerify, "dest-insecure-skip-tls-verify", s.DestTransport.InsecureSkipVerify, "skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)")
	fs.StringSliceVar(&s.DestTransport.ServerList, "dest-endpoints", s.DestTransport.ServerList, "List of etcd servers to connect with (scheme://ip:port) for the destination cluster, comma separated")
	fs.StringVar(&s.DestTransport.CACertFile, "dest-cacert", s.DestTransport.CACertFile, "Verify certificates of TLS enabled secure servers using this CA bundle for the destination cluster")
	fs.StringVar(&s.DestTransport.CertFile, "dest-cert", s.DestTransport.CertFile, "Identify secure client using this TLS certificate file for the destination cluster")
	fs.StringVar(&s.DestTransport.KeyFile, "dest-key", s.DestTransport.KeyFile, "Identify secure client using this TLS key file for the destination cluster")

	fs.DurationVar(&s.DialTimeout, "dial-timeout", s.DialTimeout, "dial timeout for client connections")
	fs.DurationVar(&s.KeepAliveTime, "keepalive-time", s.KeepAliveTime, "keepalive time for client connections")
	fs.DurationVar(&s.KeepAliveTimeout, "keepalive-timeout", s.KeepAliveTimeout, "keepalive timeout for client connections")
}

func (s *TransportOptions) Validate() []error {
	var errors []error
	if len(s.SourceTransport.ServerList) == 0 {
		errors = append(errors, fmt.Errorf("--source-endpoints must be specified"))
	}
	if len(s.DestTransport.ServerList) == 0 {
		errors = append(errors, fmt.Errorf("--dest-endpoints must be specified"))
	}
	if s.DialTimeout.Nanoseconds() < 0 {
		errors = append(errors, fmt.Errorf("--dial-timeout can not be negative value"))
	}
	if s.KeepAliveTime.Nanoseconds() < 0 {
		errors = append(errors, fmt.Errorf("--keepalive-time can not be negative value"))
	}
	if s.KeepAliveTimeout.Nanoseconds() < 0 {
		errors = append(errors, fmt.Errorf("--keepalive-timeout can not be negative value"))
	}
	return errors
}
