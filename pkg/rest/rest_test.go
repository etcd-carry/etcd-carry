package rest

import (
	"bytes"
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/testing/util"
	"github.com/etcd-carry/etcd-carry/pkg/util/signal"
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func readAll(r io.Reader) (b []byte, err error) {
	var bufa [64]byte
	buf := bufa[:]
	for {
		n, err := r.Read(buf)
		if n == 0 && err == nil {
			return nil, fmt.Errorf("read: n=0 with err=nil")
		}
		b = append(b, buf[:n]...)
		if err == io.EOF {
			n, err = r.Read(buf)
			if n != 0 || err != io.EOF {
				return nil, fmt.Errorf("read: n=%d err=%#v after EOF", n, err)
			}
			return b, nil
		}
		if err != nil {
			return b, err
		}
	}
}

func TestDaemon_Run(t *testing.T) {
	ports, err := util.GetAvailableTestPorts(1)
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewRestfulServing(net.JoinHostPort("0.0.0.0", strconv.Itoa(ports[0])))
	if err != nil {
		t.Fatal(err)
	}

	if err = d.Run(signal.SetupSignalHandler()); err != nil {
		t.Fatal(err)
	}

	client := http.DefaultClient
	client.Timeout = 3 * time.Second

	testCase := []struct {
		name  string
		ready bool
	}{
		{
			name:  "test ready",
			ready: true,
		},
		{
			name:  "test not ready",
			ready: false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			d.SetReady(tc.ready)
			resp, err := client.Get(fmt.Sprintf("http://localhost:%d/ready", ports[0]))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if err == nil {
				buf, err := readAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(buf, []byte(fmt.Sprintf("{\"ready\":%v}", tc.ready))) {
					t.Fatalf("should not happened, expect %+v, but %+v", tc.ready, string(buf))
				}
			}
		})
	}
}
