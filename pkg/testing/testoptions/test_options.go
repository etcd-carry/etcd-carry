package testoptions

import (
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/mirror/options"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testetcd"
	"github.com/spf13/pflag"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"io/ioutil"
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

	args = []string{
		"--debug=true",
		fmt.Sprintf("--mirror-rule=%s", rulesFile),
		fmt.Sprintf("--encryption-provider-config=%s", encryptionFile),
		"--kube-prefix=/registry",
		"--max-txn-ops=100",
		"--rev=520",
		"--dial-timeout=1s",
		"--keepalive-time=1s",
		"--keepalive-timeout=1s",
		"--source-insecure-skip-tls-verify=false",
		"--source-insecure-transport=true",
		"--dest-insecure-skip-tls-verify=false",
		"--dest-insecure-transport=true",
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

	sourceCertsDir, err := ioutil.TempDir(testDir, "source_certs")
	if err != nil {
		t.Fatal(err)
	}
	sourceCAFile, sourceCertFile, sourceKeyFile := testetcd.GetTLSCerts(t, sourceCertsDir)
	sourceConfig := testetcd.NewTestConfig(t)
	sourceConfig.ClientTLSInfo = transport.TLSInfo{
		TrustedCAFile: sourceCAFile,
		CertFile:      sourceCertFile,
		KeyFile:       sourceKeyFile,
	}
	for i := range sourceConfig.ACUrls {
		sourceConfig.ACUrls[i].Scheme = "https"
	}
	for i := range sourceConfig.LCUrls {
		sourceConfig.LCUrls[i].Scheme = "https"
	}
	etcdSource := testetcd.RunEtcdServer(t, sourceConfig)

	destCertsDir, err := ioutil.TempDir(testDir, "dest_certs")
	if err != nil {
		t.Fatal(err)
	}
	destCAFile, destCertFile, destKeyFile := testetcd.GetTLSCerts(t, destCertsDir)
	destConfig := testetcd.NewTestConfig(t)
	destConfig.ClientTLSInfo = transport.TLSInfo{
		TrustedCAFile: destCAFile,
		CertFile:      destCertFile,
		KeyFile:       destKeyFile,
	}
	for i := range destConfig.ACUrls {
		destConfig.ACUrls[i].Scheme = "https"
	}
	for i := range destConfig.LCUrls {
		destConfig.LCUrls[i].Scheme = "https"
	}
	etcdDest := testetcd.RunEtcdServer(t, destConfig)

	args := GetArguments(t, testDir)
	args = append(args, fmt.Sprintf("--source-cacert=%s", sourceCAFile))
	args = append(args, fmt.Sprintf("--source-cert=%s", sourceCertFile))
	args = append(args, fmt.Sprintf("--source-key=%s", sourceKeyFile))
	args = append(args, fmt.Sprintf("--source-endpoints=%s", strings.Join(etcdSource.Server.Cluster().ClientURLs(), ",")))
	args = append(args, fmt.Sprintf("--dest-cacert=%s", destCAFile))
	args = append(args, fmt.Sprintf("--dest-cert=%s", destCertFile))
	args = append(args, fmt.Sprintf("--dest-key=%s", destKeyFile))
	args = append(args, fmt.Sprintf("--dest-endpoints=%s", strings.Join(etcdDest.Server.Cluster().ClientURLs(), ",")))

	fs := pflag.NewFlagSet("mirrorcontexttest", pflag.ContinueOnError)
	o := options.NewMirrorOptions()
	for _, f := range o.Flags().FlagSets {
		fs.AddFlagSet(f)
	}
	fs.Parse(args)

	return o
}
