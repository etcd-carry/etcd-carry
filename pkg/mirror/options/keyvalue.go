package options

import (
	"fmt"
	"github.com/spf13/pflag"
)

type KeyValueOptions struct {
	KVPath string
}

func NewKeyValueOptions() *KeyValueOptions {
	return &KeyValueOptions{
		KVPath: "/var/lib/mirror/db",
	}
}

func (s *KeyValueOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.KVPath, "db-path", s.KVPath, "the path where kv-db stores data")
}

func (s *KeyValueOptions) Validate() []error {
	var errors []error
	if len(s.KVPath) == 0 {
		errors = append(errors, fmt.Errorf("--db-path must be specified"))
	}
	return errors
}
