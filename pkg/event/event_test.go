package event

import (
	"github.com/etcd-carry/etcd-carry/pkg/mirror"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testcodec"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testoptions"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testsync"
	"testing"
)

func TestEvent_ProcessEvent(t *testing.T) {
	o := testoptions.GetMirrorOptions(t)
	o.Etcd.StartReversion = 0
	mirrorCtx, err := mirrorcontext.NewMirrorContext(o)
	if err != nil {
		t.Fatal(err)
	}

	syncer := mirror.NewSyncer(mirrorCtx)
	testsync.PrepareSequentialTestData(t, mirrorCtx.MasterClient)

	testCase := []struct {
		name string
		keys map[string]bool
	}{
		{
			name: "SyncUpdates",
			keys: map[string]bool{
				string(testcodec.SampleCrdResourceKey): false,
			},
		},
	}

	_, errC := syncer.SyncSequential()
	if err := <-errC; err != nil {
		t.Error(err)
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := mirrorCtx.MasterClient.Put(mirrorCtx.Context, string(testcodec.SampleSecretMatchedKey2), string(testcodec.SampleSecretMatchedValue2)); err != nil {
				t.Fatal(err)
			}
			wc := syncer.SyncUpdates()
			wr := <-wc
			if wr.Err() != nil {
				err := wr.Err()
				t.Fatal(err)
			}
			for _, ev := range wr.Events {
				parsedEvent, err := ParseEvent(ev)
				if err != nil {
					t.Fatal(err)
				}
				if !(parsedEvent.ProcessEvent(mirrorCtx) || parsedEvent.DeleteEventMatched(mirrorCtx)) {
					t.Errorf("should not happened")
				}
			}
		})
	}
}

//func TestEvent_ProcessEvent(t *testing.T) {
//	o := testoptions.GetMirrorOptions(t)
//	o.Etcd.StartReversion = 0
//	mirrorCtx, err := mirrorcontext.NewMirrorContext(o)
//	if err != nil {
//		t.Fatal(err)
//	}
//	ctx, cancel := context.WithCancel(context.TODO())
//	mirrorCtx.Context = ctx
//	t.Cleanup(cancel)
//
//	syncer := mirror.NewSyncer(mirrorCtx)
//	testsync.PrepareSequentialTestData(t, mirrorCtx.MasterClient)
//
//	_, errC := syncer.SyncSequential()
//	if err := <-errC; err != nil {
//		t.Error(err)
//	}
//
//	testCase := map[string]struct {
//		name         string
//		action       string
//		key          []byte
//		value        []byte
//		isMatched    bool
//		expectErrors bool
//	}{
//		string(testcodec.SampleSecretMatchedKey1): {
//			name:         "delete matched secret",
//			action:       "DELETE",
//			key:          testcodec.SampleSecretMatchedKey1,
//			value:        testcodec.SampleSecretMatchedValue1,
//			isMatched:    true,
//			expectErrors: false,
//		},
//		//string(testcodec.SampleSecretMismatchedKey2): {
//		//	name:   "delete mismatched secret",
//		//	action: "DELETE",
//		//	//key:          testcodec.SampleSecretMismatchedKey2,
//		//	value:        testcodec.SampleSecretMismatchedValue2,
//		//	isMatched:    false,
//		//	expectErrors: false,
//		//},
//		//string(testcodec.SampleSecretMatchedKey2): {
//		//	name:   "put matched secret",
//		//	action: "PUT",
//		//	//key:          testcodec.SampleSecretMatchedKey2,
//		//	value:        testcodec.SampleSecretMatchedValue2,
//		//	isMatched:    true,
//		//	expectErrors: false,
//		//},
//		//string(testcodec.SampleSecretMismatchedKey1): {
//		//	name:   "put mismatched secret",
//		//	action: "PUT",
//		//	//key:          testcodec.SampleSecretMismatchedKey1,
//		//	value:        testcodec.SampleSecretMismatchedValue1,
//		//	isMatched:    false,
//		//	expectErrors: false,
//		//},
//	}
//
//	//var testCaseCount = 0
//	//for k, v := range testCase {
//	//	switch v.action {
//	//	case "DELETE":
//	//		if _, err = mirrorCtx.MasterClient.Delete(mirrorCtx, k); err != nil {
//	//			t.Error(err)
//	//		}
//	//	case "PUT":
//	//		if _, err = mirrorCtx.MasterClient.Put(mirrorCtx, k, string(v.value)); err != nil {
//	//			t.Error(err)
//	//		}
//	//	}
//	//	testCaseCount++
//	//}
//
//	//wc := syncer.SyncUpdates()
//	//for wr := range wc {
//	//	if wr.Err() != nil {
//	//		err := wr.Err()
//	//		t.Error(err)
//	//	}
//	//
//	//	if testCaseCount == 0 {
//	//		return
//	//	}
//	//	for _, ev := range wr.Events {
//	//		testCaseCount--
//	//		parsedEvent, err := ParseEvent(ev)
//	//		v, ok := testCase[string(parsedEvent.Event().Kv.Key)]
//	//		if !ok {
//	//			t.Error("should not happened")
//	//		}
//	//
//	//		fmt.Println("xxxxxxxxx:", string(parsedEvent.Event().Kv.Key), err)
//	//		if err != nil && !v.expectErrors {
//	//			t.Errorf("expect no error, but errors found %+v", err)
//	//		}
//	//		if err == nil && v.expectErrors {
//	//			t.Errorf("expect errors, but no error found")
//	//		}
//	//		if err == nil && !v.expectErrors {
//	//			if (parsedEvent.ProcessEvent(mirrorCtx) || parsedEvent.DeleteEventMatched(mirrorCtx)) != v.isMatched {
//	//				t.Errorf("should not happened")
//	//			}
//	//		}
//	//		fmt.Println("pppppppppppp1:", testCaseCount)
//	//	}
//	//	fmt.Println("pppppppppppp2:", testCaseCount)
//	//}
//	for _, tc := range testCase {
//		t.Run(t.Name(), func(t *testing.T) {
//			if _, err = mirrorCtx.MasterClient.Delete(mirrorCtx, string(tc.key)); err != nil {
//				t.Error(err)
//			}
//			wc := syncer.SyncUpdates()
//			wr := <-wc
//			if wr.Err() != nil {
//				err := wr.Err()
//				t.Fatal(err)
//			}
//
//			for _, ev := range wr.Events {
//				println("xxxxxxxxxxxxxxxxxxxxxxx")
//				parsedEvent, err := ParseEvent(ev)
//				if err != nil && !tc.expectErrors {
//					t.Errorf("expect no error, but errors found %+v", err)
//				}
//				if err == nil && tc.expectErrors {
//					t.Errorf("expect errors, but no error found")
//				}
//				if err == nil && !tc.expectErrors {
//					if (parsedEvent.ProcessEvent(mirrorCtx) || parsedEvent.DeleteEventMatched(mirrorCtx)) != tc.isMatched {
//						t.Errorf("should not happened")
//					}
//				}
//			}
//		})
//	}
//}
