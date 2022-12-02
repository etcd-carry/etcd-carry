package testoptions

import (
	"fmt"
	"github.com/spf13/pflag"
	"go.etcd.io/etcd/pkg/transport"
	"io/ioutil"
	"orcastack.io/etcd-mirror/pkg/mirror/options"
	"orcastack.io/etcd-mirror/pkg/testing/testetcd"
	"orcastack.io/etcd-mirror/pkg/testing/util"
	"os"
	"path"
	"strings"
	"testing"
)

const MirrorRulesConfig = `
filters:
  sequential:
    - group: apiextensions.k8s.io
      sequence: 1
      resources:
        - group: apiextensions.k8s.io
          version: v1beta1
          kind: CustomResourceDefinition
    - group: ""
      sequence: 2
      resources:
        - version: v1
          kind: Namespace
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: In
              values:
                - unit
                - unique
  secondary:
    - group: "*"
      namespaceSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: In
              values:
                - unique
      excludes:
        - resource:
            version: v1
            kind: Secret
          labelSelectors:
            - matchExpressions:
                - key: test.io/namespace-kind
                  operator: Exists
    - group: ""
      resources:
        - version: v1
          Kind: Secret
      namespace: unit-test
      fieldSelectors:
        - matchExpressions:
            - key: type
              operator: NotIn
              values:
                - kubernetes.io/service-account-token
      excludes:
        - resource:
            version: v1
            kind: Secret
          name: exclude-me-secret
          namespace: unit-test
    - group: ""
      resources:
        - version: v1
          kind: ConfigMap
      namespace: unit-test
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: Exists
    - group: ""
      resources:
        - version: v1
          kind: Service
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: Exists
`

const EncryptionProviderConfig = `
kind: EncryptionConfig
apiVersion: v1
resources:
  - resources:
    - secrets
    providers:
    - aescbc:
        keys:
        - name: key
          secret: GPG4RC0Vyk7+Mz/niQPttxLIeL4HF96oRCcBRyKNpfM=
    - identity: {}
`

func GetArguments(t *testing.T, testDir string) (args []string) {
	rulesFile := path.Join(testDir, "rules.yaml")
	if err := ioutil.WriteFile(rulesFile, []byte(MirrorRulesConfig), 0644); err != nil {
		t.Fatal(err)
	}
	encryptionFile := path.Join(testDir, "encryption.yaml")
	if err := ioutil.WriteFile(encryptionFile, []byte(EncryptionProviderConfig), 0644); err != nil {
		t.Fatal(err)
	}
	bindPorts, err := util.GetAvailableTestPorts(1)
	if err != nil {
		t.Fatal(err)
	}

	args = []string{
		"--debug=true",
		fmt.Sprintf("--mirror-rule=%s", rulesFile),
		"--mode=active-standby",
		fmt.Sprintf("--encryption-provider-config=%s", encryptionFile),
		"--kube-prefix=/registry",
		"--max-txn-ops=100",
		"--rev=520",
		"--dial-timeout=1s",
		"--keepalive-time=1s",
		"--keepalive-timeout=1s",
		"--master-insecure-skip-tls-verify=false",
		"--master-insecure-transport=true",
		"--slave-insecure-skip-tls-verify=false",
		"--slave-insecure-transport=true",
		fmt.Sprintf("--db-path=%s", path.Join(testDir, "db")),
		"--bind-address=0.0.0.0",
		fmt.Sprintf("--bind-port=%v", bindPorts[0]),
	}
	return args
}

func GetMirrorOptions(t *testing.T) (ctx *options.MirrorOptions) {
	baseDir := os.TempDir()
	testDir, err := ioutil.TempDir(baseDir, "etcd_mirror")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(testDir) })

	masterCertsDir, err := ioutil.TempDir(testDir, "master_certs")
	if err != nil {
		t.Fatal(err)
	}
	masterCAFile, masterCertFile, masterKeyFile := testetcd.GetTLSCerts(t, masterCertsDir)
	masterConfig := testetcd.NewTestConfig(t)
	masterConfig.ClientTLSInfo = transport.TLSInfo{
		TrustedCAFile: masterCAFile,
		CertFile:      masterCertFile,
		KeyFile:       masterKeyFile,
	}
	for i := range masterConfig.ACUrls {
		masterConfig.ACUrls[i].Scheme = "https"
	}
	for i := range masterConfig.LCUrls {
		masterConfig.LCUrls[i].Scheme = "https"
	}
	etcdMaster := testetcd.RunEtcdServer(t, masterConfig)

	slaveCertsDir, err := ioutil.TempDir(testDir, "slave_certs")
	if err != nil {
		t.Fatal(err)
	}
	slaveCAFile, slaveCertFile, slaveKeyFile := testetcd.GetTLSCerts(t, slaveCertsDir)
	slaveConfig := testetcd.NewTestConfig(t)
	slaveConfig.ClientTLSInfo = transport.TLSInfo{
		TrustedCAFile: slaveCAFile,
		CertFile:      slaveCertFile,
		KeyFile:       slaveKeyFile,
	}
	for i := range slaveConfig.ACUrls {
		slaveConfig.ACUrls[i].Scheme = "https"
	}
	for i := range slaveConfig.LCUrls {
		slaveConfig.LCUrls[i].Scheme = "https"
	}
	etcdSlave := testetcd.RunEtcdServer(t, slaveConfig)

	args := GetArguments(t, testDir)
	args = append(args, fmt.Sprintf("--master-cacert=%s", masterCAFile))
	args = append(args, fmt.Sprintf("--master-cert=%s", masterCertFile))
	args = append(args, fmt.Sprintf("--master-key=%s", masterKeyFile))
	args = append(args, fmt.Sprintf("--master-endpoints=%s", strings.Join(etcdMaster.Server.Cluster().ClientURLs(), ",")))
	args = append(args, fmt.Sprintf("--slave-cacert=%s", slaveCAFile))
	args = append(args, fmt.Sprintf("--slave-cert=%s", slaveCertFile))
	args = append(args, fmt.Sprintf("--slave-key=%s", slaveKeyFile))
	args = append(args, fmt.Sprintf("--slave-endpoints=%s", strings.Join(etcdSlave.Server.Cluster().ClientURLs(), ",")))

	fs := pflag.NewFlagSet("mirrorcontexttest", pflag.ContinueOnError)
	o := options.NewMirrorOptions()
	for _, f := range o.Flags().FlagSets {
		fs.AddFlagSet(f)
	}
	fs.Parse(args)

	return o
}
