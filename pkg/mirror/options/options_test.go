package options

import (
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/errors"
	"os"
	"testing"
)

func TestMirrorOptions_Flags(t *testing.T) {
	fs := pflag.NewFlagSet("addflagstest", pflag.ContinueOnError)
	o := NewMirrorOptions()
	for _, f := range o.Flags().FlagSets {
		fs.AddFlagSet(f)
	}

	testCase := []struct {
		name     string
		args     []string
		errorMsg string
		isValid  bool
	}{
		{
			name: "valid Mirror options",
			args: []string{
				"--debug=true",
				"--mirror-rule=/etc/etcd/mirror/rules.yaml",
				"--encryption-provider-config=/etc/etcd/mirror/encryption.yaml",
				"--kube-prefix=/registry",
				"--max-txn-ops=100",
				"--rev=520",
				"--dial-timeout=1s",
				"--keepalive-time=1s",
				"--keepalive-timeout=1s",
				"--master-cacert=/etc/etcd/mirror/master/ca.crt",
				"--master-cert=/etc/etcd/mirror/master/server.crt",
				"--master-key=/etc/etcd/mirror/master/server.key",
				"--master-endpoints=10.20.100.111:2379",
				"--master-insecure-skip-tls-verify=false",
				"--master-insecure-transport=true",
				"--slave-cacert=/etc/etcd/mirror/slave/ca.crt",
				"--slave-cert=/etc/etcd/mirror/slave/server.crt",
				"--slave-key=/etc/etcd/mirror/slave/server.key",
				"--slave-endpoints=10.20.100.112:2379",
				"--slave-insecure-skip-tls-verify=false",
				"--slave-insecure-transport=true",
			},
			errorMsg: "",
			isValid:  true,
		},
		{
			name:     "invalid Etcd --rev option, negative value",
			args:     []string{"--rev=-1"},
			errorMsg: "--rev can not be negative value",
			isValid:  false,
		},
		{
			name:     "invalid Etcd --max-txn-ops option, negative value",
			args:     []string{"--max-txn-ops=-1"},
			errorMsg: "--max-txn-ops can not be negative value",
			isValid:  false,
		},
		{
			name:     "invalid Transport --dial-timeout option, negative value",
			args:     []string{"--dial-timeout=-1"},
			errorMsg: "--dial-timeout can not be negative value",
			isValid:  false,
		},
		{
			name:     "invalid Transport --keepalive-time option, negative value",
			args:     []string{"--keepalive-time=-1"},
			errorMsg: "--keepalive-time can not be negative value",
			isValid:  false,
		},
		{
			name:     "invalid Transport --keepalive-timeout option, negative value",
			args:     []string{"--keepalive-timeout=-1"},
			errorMsg: "--keepalive-timeout can not be negative value",
			isValid:  false,
		},
		{
			name:     "invalid Transport --master-endpoints option",
			args:     []string{""},
			errorMsg: "--master-endpoints must be specified",
			isValid:  false,
		},
		{
			name:     "invalid Transport --slave-endpoints option",
			args:     []string{""},
			errorMsg: "--slave-endpoints must be specified",
			isValid:  false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			fs.Parse(tc.args)
			errList := o.Validation()
			if len(errList) != 0 && tc.isValid {
				t.Errorf("expected no errors, but error found %+v", errors.NewAggregate(errList))
			}
			if len(errList) == 0 && !tc.isValid {
				t.Errorf("expected errors %+v, but no errors found", tc.errorMsg)
			}
			PrintSections(os.Stderr, o.Flags())
		})
	}
}
