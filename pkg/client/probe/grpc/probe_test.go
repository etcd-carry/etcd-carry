package grpc

import (
	"github.com/etcd-carry/etcd-carry/pkg/testing/testetcd"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestGrpcProbe_Probe(t *testing.T) {
	etcdServer := testetcd.RunEtcdServer(t, nil)
	tlsConfig, err := etcdServer.Config().ClientTLSInfo.ClientConfig()
	if err != nil {
		t.Fatal(err)
	}
	client1, err := clientv3.New(clientv3.Config{
		TLS:         tlsConfig,
		Endpoints:   etcdServer.Server.Cluster().ClientURLs(),
		DialTimeout: 2 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, msg, err := New().Probe(client1.ActiveConnection())
	if err == nil && result == Success {
		t.Log("probe ok!")
	} else {
		t.Fatalf("probe failed: %v", msg)
	}

	client2, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"empty"},
	})
	if err != nil {
		t.Fatal(err)
	}
	result, _, err = New().Probe(client2.ActiveConnection())
	if err == nil || result != Failure {
		t.Fatalf("should not have happened")
	}
}
