package wait

import (
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/wait"
	mirrorcontext "orcastack.io/etcd-mirror/pkg/mirror/context"
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
