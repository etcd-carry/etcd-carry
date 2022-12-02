package testetcd

import (
	"go.etcd.io/etcd/embed"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"net"
	"net/url"
	"orcastack.io/etcd-mirror/pkg/testing/util"
	"os"
	"strconv"
	"testing"
	"time"
)

func NewTestConfig(t *testing.T) *embed.Config {
	ports, err := util.GetAvailableTestPorts(2)
	if err != nil {
		t.Fatal(err)
	}
	clientURL := url.URL{Scheme: "http", Host: net.JoinHostPort("localhost", strconv.Itoa(ports[0]))}
	peerURL := url.URL{Scheme: "http", Host: net.JoinHostPort("localhost", strconv.Itoa(ports[1]))}

	cfg := embed.NewConfig()
	cfg.LPUrls = []url.URL{peerURL}
	cfg.APUrls = []url.URL{peerURL}
	cfg.LCUrls = []url.URL{clientURL}
	cfg.ACUrls = []url.URL{clientURL}
	cfg.InitialCluster = cfg.InitialClusterFromName(cfg.Name)
	cfg.ZapLoggerBuilder = embed.NewZapCoreLoggerBuilder(zaptest.NewLogger(t, zaptest.Level(zapcore.ErrorLevel)).Named("etcd-server"), nil, nil)
	cfg.Dir = t.TempDir()
	os.Chmod(cfg.Dir, 0700)

	return cfg
}

func RunEtcdServer(t *testing.T, cfg *embed.Config) *embed.Etcd {
	if cfg == nil {
		cfg = NewTestConfig(t)
	}

	e, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(e.Close)

	select {
	case <-e.Server.ReadyNotify():
	case <-time.After(30 * time.Second):
		e.Server.Stop()
		t.Fatal("etcd server start timeout")
	}

	return e
}
