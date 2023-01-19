package wait

import (
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/clock"
	"time"
)

func UntilSucceed(ctx *mirrorcontext.Context, f func(ctx *mirrorcontext.Context) error, period time.Duration, stopCh <-chan struct{}) {
	backoffUntil(ctx, f, wait.NewJitteredBackoffManager(period, 0.0, &clock.RealClock{}), stopCh)
}

func backoffUntil(ctx *mirrorcontext.Context, f func(ctx *mirrorcontext.Context) error, backoff wait.BackoffManager, stopCh <-chan struct{}) {
	var t clock.Timer
	for {
		select {
		case <-stopCh:
			return
		default:
		}

		if err := f(ctx); err == nil {
			return
		}

		t = backoff.Backoff()

		select {
		case <-stopCh:
			return
		case <-t.C():
		}
	}
}
