package wait

import (
	"errors"
	mirrorcontext "orcastack.io/etcd-mirror/pkg/mirror/context"
	"testing"
)

func TestUntilSucceed(t *testing.T) {
	ctx := &mirrorcontext.Context{}

	ch := make(chan struct{})
	close(ch)
	UntilSucceed(ctx, func(ctx *mirrorcontext.Context) error {
		t.Fatal("should not have been invoked")
		return nil
	}, 0, ch)

	ch = make(chan struct{})
	called := make(chan struct{})
	go func() {
		UntilSucceed(ctx, func(ctx *mirrorcontext.Context) error {
			called <- struct{}{}
			return nil
		}, 0, ch)
	}()
	<-called

	go func() {
		retry := 0
		UntilSucceed(ctx, func(ctx *mirrorcontext.Context) error {
			retry++
			if retry == 3 {
				return nil
			}
			called <- struct{}{}
			return errors.New("function runs error")
		}, 0, ch)
	}()
	<-called
	<-called
	close(called)
	close(ch)
}
