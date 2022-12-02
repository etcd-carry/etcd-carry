package grpc

import (
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"orcastack.io/etcd-mirror/pkg/testing/testetcd"
	"testing"
	"time"
)

func TestGrpcProbe_Probe(t *testing.T) {
	etcdServer := testetcd.RunEtcdServer(t, nil)
	tlsConfig, err := etcdServer.Config().ClientTLSInfo.ClientConfig()
	if err != nil {
		t.Fatal(err)
	}
	client, err := clientv3.New(clientv3.Config{
		TLS:         tlsConfig,
		Endpoints:   etcdServer.Server.Cluster().ClientURLs(),
		DialTimeout: 2 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}

	result, msg, err := New().Probe(client.ActiveConnection())
	if err == nil && result == Success {
		t.Log("probe ok!")
	} else {
		t.Fatalf("probe failed: %v", msg)
	}

	etcdServer.Close()
	result, _, err = New().Probe(client.ActiveConnection())
	if err == nil || result != Failure {
		t.Fatalf("should not have happened")
	}
}
