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
				"--source-cacert=/etc/etcd/mirror/source/ca.crt",
				"--source-cert=/etc/etcd/mirror/source/server.crt",
				"--source-key=/etc/etcd/mirror/source/server.key",
				"--source-endpoints=10.20.100.111:2379",
				"--source-insecure-skip-tls-verify=false",
				"--source-insecure-transport=true",
				"--dest-cacert=/etc/etcd/mirror/dest/ca.crt",
				"--dest-cert=/etc/etcd/mirror/dest/server.crt",
				"--dest-key=/etc/etcd/mirror/dest/server.key",
				"--dest-endpoints=10.20.100.112:2379",
				"--dest-insecure-skip-tls-verify=false",
				"--dest-insecure-transport=true",
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
			name:     "invalid Transport --source-endpoints option",
			args:     []string{""},
			errorMsg: "--source-endpoints must be specified",
			isValid:  false,
		},
		{
			name:     "invalid Transport --dest-endpoints option",
			args:     []string{""},
			errorMsg: "--dest-endpoints must be specified",
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
