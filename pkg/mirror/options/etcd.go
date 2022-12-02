package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

type EtcdOptions struct {
	KubePrefix                       string
	StartReversion                   int64
	MaxTxnOps                        uint
	EncryptionProviderConfigFilepath string
}

const defaultMaxTxnOps = uint(128)

func NewEtcdOptions() *EtcdOptions {
	return &EtcdOptions{
		KubePrefix:                       "/registry",
		StartReversion:                   0,
		MaxTxnOps:                        defaultMaxTxnOps,
		EncryptionProviderConfigFilepath: "/etc/mirror/secrets-encryption.yaml",
	}
}

func (s *EtcdOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.KubePrefix, "kube-prefix", s.KubePrefix, "the prefix to all kubernetes resources passed to etcd")
	fs.Int64Var(&s.StartReversion, "rev", s.StartReversion, "Specify the kv revision to start to mirror")
	fs.UintVar(&s.MaxTxnOps, "max-txn-ops", s.MaxTxnOps, "Maximum number of operations permitted in a transaction during syncing updates")
	fs.StringVar(&s.EncryptionProviderConfigFilepath, "encryption-provider-config", s.EncryptionProviderConfigFilepath, "The file containing configuration for encryption providers to be used for storing secrets in etcd")
}

func (s *EtcdOptions) Validate() []error {
	var errors []error
	if s.StartReversion < 0 {
		errors = append(errors, fmt.Errorf("--rev can not be negative value"))
	}

	return errors
}
