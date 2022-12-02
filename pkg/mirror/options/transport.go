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
	MasterTransport  TransportConfig
	SlaveTransport   TransportConfig
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
		MasterTransport: TransportConfig{
			Insecure:           true,
			InsecureSkipVerify: false,
			CACertFile:         "/etc/kubernetes/master/etcd/ca.crt",
			CertFile:           "/etc/kubernetes/master/etcd/server.crt",
			KeyFile:            "/etc/kubernetes/master/etcd/server.key",
		},
		SlaveTransport: TransportConfig{
			Insecure:           true,
			InsecureSkipVerify: false,
			CACertFile:         "/etc/kubernetes/slave/etcd/ca.crt",
			CertFile:           "/etc/kubernetes/slave/etcd/server.crt",
			KeyFile:            "/etc/kubernetes/slave/etcd/server.key",
		},
		DialTimeout:      defaultDialTimeout,
		KeepAliveTime:    defaultKeepAliveTime,
		KeepAliveTimeout: defaultKeepAliveTimeOut,
	}
}

func (s *TransportOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.MasterTransport.Insecure, "master-insecure-transport", s.MasterTransport.Insecure, "disable transport security for client connections")
	fs.BoolVar(&s.MasterTransport.InsecureSkipVerify, "master-insecure-skip-tls-verify", s.MasterTransport.InsecureSkipVerify, "skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)")
	fs.StringSliceVar(&s.MasterTransport.ServerList, "master-endpoints", s.MasterTransport.ServerList, "List of etcd servers to connect with (scheme://ip:port), comma separated")
	fs.StringVar(&s.MasterTransport.CACertFile, "master-cacert", s.MasterTransport.CACertFile, "verify certificates of TLS-enabled secure servers using this CA bundle")
	fs.StringVar(&s.MasterTransport.CertFile, "master-cert", s.MasterTransport.CertFile, "identify secure client using this TLS certificate file")
	fs.StringVar(&s.MasterTransport.KeyFile, "master-key", s.MasterTransport.KeyFile, "identify secure client using this TLS key file")

	fs.BoolVar(&s.SlaveTransport.Insecure, "slave-insecure-transport", s.SlaveTransport.Insecure, "Disable transport security for client connections for the destination cluster")
	fs.BoolVar(&s.SlaveTransport.InsecureSkipVerify, "slave-insecure-skip-tls-verify", s.SlaveTransport.InsecureSkipVerify, "skip server certificate verification (CAUTION: this option should be enabled only for testing purposes)")
	fs.StringSliceVar(&s.SlaveTransport.ServerList, "slave-endpoints", s.SlaveTransport.ServerList, "List of etcd servers to connect with (scheme://ip:port) for the destination cluster, comma separated")
	fs.StringVar(&s.SlaveTransport.CACertFile, "slave-cacert", s.SlaveTransport.CACertFile, "Verify certificates of TLS enabled secure servers using this CA bundle for the destination cluster")
	fs.StringVar(&s.SlaveTransport.CertFile, "slave-cert", s.SlaveTransport.CertFile, "Identify secure client using this TLS certificate file for the destination cluster")
	fs.StringVar(&s.SlaveTransport.KeyFile, "slave-key", s.SlaveTransport.KeyFile, "Identify secure client using this TLS key file for the destination cluster")

	fs.DurationVar(&s.DialTimeout, "dial-timeout", s.DialTimeout, "dial timeout for client connections")
	fs.DurationVar(&s.KeepAliveTime, "keepalive-time", s.KeepAliveTime, "keepalive time for client connections")
	fs.DurationVar(&s.KeepAliveTimeout, "keepalive-timeout", s.KeepAliveTimeout, "keepalive timeout for client connections")
}

func (s *TransportOptions) Validate() []error {
	var errors []error
	if len(s.MasterTransport.ServerList) == 0 {
		errors = append(errors, fmt.Errorf("--master-endpoints must be specified"))
	}
	if len(s.SlaveTransport.ServerList) == 0 {
		errors = append(errors, fmt.Errorf("--slave-endpoints must be specified"))
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
