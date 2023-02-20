package options

import (
	"github.com/spf13/pflag"
)

type GenericOptions struct {
	Debug                     bool
	MirrorRulesConfigFilepath string
}

func NewGenericOptions() *GenericOptions {
	return &GenericOptions{
		Debug:                     false,
		MirrorRulesConfigFilepath: "/etc/mirror/rules.yaml",
	}
}

func (s *GenericOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Debug, "debug", s.Debug, "enable client-side debug logging")
	fs.StringVar(&s.MirrorRulesConfigFilepath, "mirror-rule", s.MirrorRulesConfigFilepath, "Specify the rules to start mirroring")
}

func (s *GenericOptions) Validation() []error {
	var errors []error
	return errors
}
