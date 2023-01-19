package v3

import (
	"context"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testetcd"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	baseDir := os.TempDir()
	certsDir, err := ioutil.TempDir(baseDir, "etcd_certs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(certsDir)
	caFile, certFile, keyFile := testetcd.GetTLSCerts(t, certsDir)

	etcdConfig := testetcd.NewTestConfig(t)
	etcdConfig.ClientTLSInfo = transport.TLSInfo{
		TrustedCAFile: caFile,
		CertFile:      certFile,
		KeyFile:       keyFile,
	}
	for i := range etcdConfig.ACUrls {
		etcdConfig.ACUrls[i].Scheme = "https"
	}
	for i := range etcdConfig.LCUrls {
		etcdConfig.LCUrls[i].Scheme = "https"
	}

	etcdServer := testetcd.RunEtcdServer(t, etcdConfig)

	client, err := New(ConfigSpec{
		Endpoints:        etcdServer.Server.Cluster().ClientURLs(),
		DialTimeout:      2 * time.Second,
		DialOptions:      []grpc.DialOption{grpc.WithBlock()},
		KeepAliveTime:    2 * time.Second,
		KeepAliveTimeout: 2 * time.Second,
		Secure: &SecureCfg{
			Cacert:             caFile,
			Cert:               certFile,
			Key:                keyFile,
			InsecureTransport:  true,
			InsecureSkipVerify: false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err = client.Get(context.TODO(), "foo"); err != nil {
		t.Fatal(err)
	}
}
