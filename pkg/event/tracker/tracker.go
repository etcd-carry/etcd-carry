package tracker

import (
	"context"
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/constant"
	"github.com/etcd-carry/etcd-carry/pkg/event"
	"github.com/etcd-carry/etcd-carry/pkg/kv"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"go.etcd.io/etcd/clientv3"
	"strconv"
)

type EventTracker struct {
	kv.KeyValueDB
	ctx    context.Context
	cancel context.CancelFunc
	errc   chan error
}

func NewEventTracker(inCtx *mirrorcontext.Context) (*EventTracker, error) {
	instance, err := kv.NewRocksDBStore(inCtx.KeyValue.KVPath)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(inCtx)
	return &EventTracker{
		KeyValueDB: instance,
		ctx:        ctx,
		cancel:     cancel,
		errc:       make(chan error, 1),
	}, nil
}

func (t *EventTracker) Track(inCtx *mirrorcontext.Context) {
	resp, err := inCtx.MasterClient.Get(inCtx, "foo")
	if err != nil {
		t.errc <- err
	}

	go func() {
		defer fmt.Println("watch is stopping...")
		watchCh := inCtx.MasterClient.Watch(t.ctx, constant.KubeRootPrefix, clientv3.WithPrefix(), clientv3.WithRev(resp.Header.Revision), clientv3.WithPrevKV())
		for wr := range watchCh {
			if wr.Err() != nil {
				t.errc <- wr.Err()
				return
			}

			for _, ev := range wr.Events {
				parsedEvent, err := event.ParseEvent(ev)
				if err != nil {
					t.errc <- err
					return
				}

				if !parsedEvent.DeleteEventMatched(inCtx) {
					continue
				}

				fmt.Println("DELETE ", string(parsedEvent.Event().Kv.Key))
				if err := t.Put(parsedEvent.Event().Kv.Key, []byte(strconv.FormatInt(wr.Header.GetRevision(), 10))); err != nil {
					fmt.Println("track delete event error:", err)
					t.errc <- err
					return
				}
			}
		}
	}()
}

func (t *EventTracker) Replay(ctx *mirrorcontext.Context, stopCh <-chan struct{}) {
	defer close(t.errc)

	t.cancel()
	select {
	case <-t.ctx.Done():
		fmt.Println("context is stopping...")
	case <-stopCh:
		fmt.Println("do nothing, exiting...")
		return
	}

	fmt.Println("start to replay...")
	var ops []clientv3.Op
	it := t.Iterator()
	for it.SeekToFirst(); it.Valid(); it.Next() {
		if len(ops) == int(ctx.Etcd.MaxTxnOps) {
			_, err := ctx.SlaveClient.Txn(ctx).Then(ops...).Commit()
			if err != nil {
				t.errc <- err
				return
			}
			ops = []clientv3.Op{}
		}
		ops = append(ops, clientv3.OpDelete(string(it.Key().Data())))
	}

	if len(ops) != 0 {
		_, err := ctx.SlaveClient.Txn(ctx).Then(ops...).Commit()
		if err != nil {
			t.errc <- err
			return
		}
	}
}

func (t *EventTracker) Err() chan error {
	return t.errc
}
