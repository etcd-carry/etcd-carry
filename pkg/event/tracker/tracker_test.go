package tracker

import (
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/client/probe/grpc"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testcodec"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testoptions"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testsync"
	"github.com/etcd-carry/etcd-carry/pkg/util/signal"
	"github.com/etcd-carry/etcd-carry/pkg/util/wait"
	"testing"
	"time"
)

func TestNewEventTracker(t *testing.T) {
	stopCh := signal.SetupSignalHandler()
	o := testoptions.GetMirrorOptions(t)
	o.Etcd.StartReversion = 0
	mirrorCtx, err := mirrorcontext.NewMirrorContext(o)
	if err != nil {
		t.Fatal(err)
	}
	testsync.PrepareTestData(t, mirrorCtx.MasterClient)
	testsync.PrepareTestData(t, mirrorCtx.SlaveClient)

	deleteMatchedCase := map[string]string{
		string(testcodec.SampleCrdKey):               string(testcodec.SampleCrdValue),
		string(testcodec.SampleConfigmapMatchedKey1): string(testcodec.SampleConfigmapMatchedValue1),
	}
	deleteMismatchedCase := map[string]string{
		string(testcodec.SampleSecretMismatchedKey1):    string(testcodec.SampleSecretMismatchedValue1),
		string(testcodec.SampleNamespaceMismatchedKey1): string(testcodec.SampleNamespaceMismatchedValue1),
	}
	putCase := map[string]string{
		"/registry/test1": "{TestString: \"foo\"}",
		"/registry/test2": "{\"TestString\": \"foo\"}",
	}

	tracker, err := NewEventTracker(mirrorCtx)
	if err != nil {
		t.Fatal(err)
	}
	tracker.Track(mirrorCtx)

	for k := range deleteMatchedCase {
		if _, err = mirrorCtx.MasterClient.Delete(mirrorCtx.Context, k); err != nil {
			t.Error(err)
		}
	}

	for k := range deleteMismatchedCase {
		if _, err = mirrorCtx.MasterClient.Delete(mirrorCtx.Context, k); err != nil {
			t.Error(err)
		}
	}

	for k, v := range putCase {
		if _, err = mirrorCtx.MasterClient.Put(mirrorCtx.Context, k, v); err != nil {
			t.Error(err)
		}
	}

	probeFunc := func(ctx *mirrorcontext.Context) error {
		result, msg, err := ctx.Probe.Probe(ctx.SlaveClient.ActiveConnection())
		if err == nil && result == grpc.Success {
			t.Log("probe ok!")
			return nil
		}
		t.Logf("probe failed: %v", msg)
		return fmt.Errorf(msg)
	}
	// wait slave etcd server become available
	wait.UntilSucceed(mirrorCtx, probeFunc, 3*time.Second, stopCh)

	time.Sleep(5 * time.Second)
	tracker.Replay(mirrorCtx, stopCh)
	if err = <-tracker.Err(); err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)
	for k := range deleteMatchedCase {
		resp, err := mirrorCtx.SlaveClient.Get(mirrorCtx.Context, k)
		if err != nil {
			t.Error(err)
		}
		if resp.Count != 0 {
			t.Error("should not happened")
		}
	}
	for k := range deleteMismatchedCase {
		resp, err := mirrorCtx.SlaveClient.Get(mirrorCtx.Context, k)
		if err != nil {
			t.Error(err)
		}
		if resp.Count == 0 {
			t.Error("should not happened")
		}
	}
	for k := range putCase {
		resp, err := mirrorCtx.SlaveClient.Get(mirrorCtx.Context, k)
		if err != nil {
			t.Error(err)
		}
		if resp.Count != 0 {
			t.Error("should not happened")
		}
	}
}
