package app

import (
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/event"
	"github.com/etcd-carry/etcd-carry/pkg/mirror"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"github.com/etcd-carry/etcd-carry/pkg/mirror/options"
	"github.com/etcd-carry/etcd-carry/pkg/util/signal"
	"github.com/spf13/cobra"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"k8s.io/apimachinery/pkg/util/errors"
)

func NewEtcdMirrorCommand() *cobra.Command {
	o := options.NewMirrorOptions()
	c := &cobra.Command{
		Use:   "etcd-carry",
		Short: "A simple command line for etcd mirroring",
		RunE: func(cmd *cobra.Command, args []string) error {
			if errs := o.Validation(); len(errs) != 0 {
				return errors.NewAggregate(errs)
			}
			return Run(o, signal.SetupSignalHandler())
		},
	}

	fs := c.Flags()
	sfs := o.Flags()
	for _, f := range sfs.FlagSets {
		fs.AddFlagSet(f)
	}

	c.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Fprintf(c.OutOrStderr(), "Usage:\n  %s\n", c.UseLine())
		options.PrintSections(c.OutOrStderr(), sfs)
		return nil
	})
	c.SetHelpFunc(func(c *cobra.Command, args []string) {
		fmt.Fprintf(c.OutOrStdout(), "%s\n\nUsage:\n  %s\n", c.Short, c.UseLine())
		options.PrintSections(c.OutOrStdout(), sfs)
	})

	return c
}

func Run(o *options.MirrorOptions, stopCh <-chan struct{}) error {
	mirrorCtx, err := mirrorcontext.NewMirrorContext(o)
	if err != nil {
		return err
	}

	return makeMirror(mirrorCtx, stopCh)
}

func makeMirror(ctx *mirrorcontext.Context, stopCh <-chan struct{}) error {
	select {
	case <-stopCh:
		return nil
	default:
	}

	startRev := ctx.MirrorOptions.Etcd.StartReversion
	if startRev < 0 {
		startRev = 0
	}

	s := mirror.NewSyncer(ctx)

	if startRev == 0 {
		seqRC, seqErrC := s.SyncSequential()
		for r := range seqRC {
			for _, kv := range r.Kvs {
				fmt.Println("PUT", string(kv.Key))
				if _, err := ctx.DestClient.Put(ctx, string(kv.Key), string(kv.Value)); err != nil {
					return err
				}
			}
		}
		if err := <-seqErrC; err != nil {
			return err
		}

		secRC, secErrC := s.SyncSecondary()

		for r := range secRC {
			for _, kv := range r.Kvs {
				fmt.Println("PUT", string(kv.Key))
				if _, err := ctx.DestClient.Put(ctx, string(kv.Key), string(kv.Value)); err != nil {
					return err
				}
			}
		}
		if err := <-secErrC; err != nil {
			return err
		}
	}

	wc := s.SyncUpdates()
	for wr := range wc {
		if wr.Err() != nil {
			err := wr.Err()
			// If there is an error on server (e.g. compaction), the channel will return it before closed.
			return err
		}

		var lastRev int64
		var ops []clientv3.Op

		for _, ev := range wr.Events {
			nextRev := ev.Kv.ModRevision
			if lastRev != 0 && nextRev > lastRev {
				_, err := ctx.DestClient.Txn(ctx).Then(ops...).Commit()
				if err != nil {
					return err
				}
				ops = []clientv3.Op{}
			}
			lastRev = nextRev

			if len(ops) == int(ctx.Etcd.MaxTxnOps) {
				_, err := ctx.DestClient.Txn(ctx).Then(ops...).Commit()
				if err != nil {
					return err
				}
				ops = []clientv3.Op{}
			}

			parsedEvent, err := event.ParseEvent(ev)
			if err != nil {
				return err
			}

			if !parsedEvent.ProcessEvent(ctx) {
				continue
			}

			fmt.Println(parsedEvent.Event().Type.String(), string(parsedEvent.Event().Kv.Key))

			switch parsedEvent.Event().Type {
			case mvccpb.PUT:
				ops = append(ops, clientv3.OpPut(string(parsedEvent.Event().Kv.Key), string(parsedEvent.Event().Kv.Value)))
			case mvccpb.DELETE:
				ops = append(ops, clientv3.OpDelete(string(parsedEvent.Event().Kv.Key)))
			default:
				panic("unexpected event type")
			}
		}

		if len(ops) != 0 {
			_, err := ctx.DestClient.Txn(ctx).Then(ops...).Commit()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
