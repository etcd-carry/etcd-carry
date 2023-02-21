package mirror

import (
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testcodec"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testoptions"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testsync"
	"strings"
	"testing"
)

func TestNewSyncer(t *testing.T) {
	o := testoptions.GetMirrorOptions(t)
	o.Etcd.StartReversion = 0
	mirrorCtx, err := mirrorcontext.NewMirrorContext(o)
	if err != nil {
		t.Fatal(err)
	}

	syncer := NewSyncer(mirrorCtx)
	testsync.PrepareTestData(t, mirrorCtx.SourceClient)

	testCase := []struct {
		name string
		keys map[string]bool
	}{
		{
			name: "SyncSequential_matched",
			keys: map[string]bool{
				string(testcodec.SampleCrdKey):               true,
				string(testcodec.SampleNamespaceMatchedKey1): true,
				string(testcodec.SampleNamespaceMatchedKey2): true,
			},
		},
		{
			name: "SyncSequential_mismatched",
			keys: map[string]bool{
				string(testcodec.SampleNamespaceMismatchedKey1): false,
			},
		},
		{
			name: "SyncSecondary_matched",
			keys: map[string]bool{
				string(testcodec.SampleConfigmapMatchedKey1): true,
				string(testcodec.SampleSecretMatchedKey1):    true,
				string(testcodec.SampleSecretMatchedKey2):    true,
				string(testcodec.SampleServiceMatchedKey1):   true,
			},
		},
		{
			name: "SyncSecondary_mismatched",
			keys: map[string]bool{
				string(testcodec.SampleConfigmapMismatchedKey1): false,
				string(testcodec.SampleSecretMismatchedKey1):    false,
				string(testcodec.SampleSecretMismatchedKey2):    false,
			},
		},
		{
			name: "SyncUpdates",
			keys: map[string]bool{
				string(testcodec.SampleCrdResourceKey): false,
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			if strings.HasPrefix(tc.name, "SyncSequential") {
				seqRC, errC := syncer.SyncSequential()
				if err := <-errC; err != nil {
					t.Fatal(err)
				}
				switch tc.name {
				case "SyncSequential_matched":
					for r := range seqRC {
						for _, kv := range r.Kvs {
							if !tc.keys[string(kv.Key)] {
								t.Fatal("should not happened")
							}
						}
					}
				case "SyncSequential_mismatched":
					for r := range seqRC {
						for _, kv := range r.Kvs {
							if tc.keys[string(kv.Key)] {
								t.Fatal("should not happened")
							}
						}
					}
				}
			} else if strings.HasPrefix(tc.name, "SyncSecondary") {
				secRC, errC := syncer.SyncSecondary()
				if err := <-errC; err != nil {
					t.Fatal(err)
				}
				switch tc.name {
				case "SyncSecondary_matched":
					for r := range secRC {
						for _, kv := range r.Kvs {
							if !tc.keys[string(kv.Key)] {
								t.Fatal("should not happened")
							}
						}
					}
				case "SyncSecondary_mismatched":
					for r := range secRC {
						for _, kv := range r.Kvs {
							if tc.keys[string(kv.Key)] {
								t.Fatal("should not happened")
							}
						}
					}
				}
			} else if strings.HasPrefix(tc.name, "SyncUpdates") {
				if _, err := mirrorCtx.SourceClient.Put(mirrorCtx.Context, string(testcodec.SampleCrdResourceKey), string(testcodec.SampleCrdResourceValue)); err != nil {
					t.Fatal(err)
				}
				wc := syncer.SyncUpdates()
				wr := <-wc
				if wr.Err() != nil {
					err := wr.Err()
					t.Fatal(err)
				}
				for _, e := range wr.Events {
					if tc.keys[string(e.Kv.Key)] {
						t.Fatal("should not happened")
					}
				}
			}
		})
	}
}
