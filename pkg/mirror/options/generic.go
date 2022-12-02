package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
)

const (
	StandAloneMode    = "standalone"
	ActiveStandbyMode = "active-standby"
)

type GenericOptions struct {
	Debug                     bool
	Mode                      string
	MirrorRulesConfigFilepath string
}

func NewGenericOptions() *GenericOptions {
	return &GenericOptions{
		Debug:                     false,
		Mode:                      StandAloneMode,
		MirrorRulesConfigFilepath: "/etc/mirror/rules.yaml",
	}
}

func (s *GenericOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.Debug, "debug", s.Debug, "enable client-side debug logging")
	fs.StringVar(&s.Mode, "mode", s.Mode, "running mode, standalone or active-standby")
	fs.StringVar(&s.MirrorRulesConfigFilepath, "mirror-rule", s.MirrorRulesConfigFilepath, "Specify the rules to start mirroring")
}

func (s *GenericOptions) Validation() []error {
	var errors []error
	supportedMode := sets.NewString(StandAloneMode, ActiveStandbyMode)
	if !supportedMode.Has(s.Mode) {
		errors = append(errors, fmt.Errorf("--mode only supports %v", strings.Join(supportedMode.List(), ",")))
	}

	return errors
}
