package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

type RestfulServingOptions struct {
	BindAddress net.IP
	BindPort    int
}

func NewDaemonOptions() *RestfulServingOptions {
	return &RestfulServingOptions{
		BindAddress: net.ParseIP("0.0.0.0"),
		BindPort:    10520,
	}
}

func (s *RestfulServingOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IPVar(&s.BindAddress, "bind-address", s.BindAddress, "the address the metric endpoint and ready/healthz binds to")
	fs.IntVar(&s.BindPort, "bind-port", s.BindPort, "the port on which to serve restful")
}

func (s *RestfulServingOptions) Validate() []error {
	var errors []error
	if s.BindPort < 0 || s.BindPort > 65535 {
		errors = append(errors, fmt.Errorf("--bind-port %v must be between 0 and 65535", s.BindPort))
	}
	return errors
}
