package event

import (
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/constant"
	kubeschema "github.com/etcd-carry/etcd-carry/pkg/filter/kube/schema"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"go.etcd.io/etcd/client/v3"
	corev1 "k8s.io/api/core/v1"
)

type MirrorEvent interface {
	Event() *clientv3.Event
	ProcessEvent(ctx *mirrorcontext.Context) bool
	DeleteEventMatched(ctx *mirrorcontext.Context) bool
}

type event struct {
	ev        *clientv3.Event
	isDeleted bool
	isCreated bool
}

func ParseEvent(e *clientv3.Event) (MirrorEvent, error) {
	if !e.IsCreate() && e.PrevKv == nil {
		// If the previous value is nil, error. One example of how this is possible is if the previous value has been compacted already.
		return nil, fmt.Errorf("etcd event received with PrevKv=nil (key=%q, modRevision=%d, type=%s)", string(e.Kv.Key), e.Kv.ModRevision, e.Type.String())
	}

	ret := &event{
		ev:        e,
		isDeleted: e.Type == clientv3.EventTypeDelete,
		isCreated: e.IsCreate(),
	}

	return ret, nil
}

func (e *event) Event() *clientv3.Event {
	return e.ev
}

func (e *event) ProcessEvent(ctx *mirrorcontext.Context) bool {
	var val []byte

	if !e.isDeleted {
		val = e.ev.Kv.Value
	}

	// We need to decode prevValue, only if this is deletion event
	if e.ev.PrevKv != nil && len(e.ev.PrevKv.Value) > 0 && e.isDeleted {
		val = e.ev.PrevKv.Value
	}

	resource, err := kubeschema.Decode(ctx, e.ev.Kv.Key, val)
	if err != nil {
		return false
	}

	seqMatched := resource.FilterSequentialByRules(ctx)
	if seqMatched {
		if corev1.SchemeGroupVersion.WithKind(constant.KubeKindNamespace) == resource.GVK() {
			if !e.isDeleted {
				ctx.MirrorFilter.Namespace[resource.Object().GetName()] = resource.Object().GetLabels()
			}

			if e.ev.PrevKv != nil && len(e.ev.PrevKv.Value) > 0 && e.isDeleted {
				delete(ctx.MirrorFilter.Namespace, resource.Object().GetName())
			}
		}
	}

	secondaryMatched := resource.FilterSecondaryByRules(ctx)
	if secondaryMatched {
		kv, err := resource.MutateKubeResource(e.ev.Kv)
		if err != nil {
			return false
		}
		e.ev.Kv = kv
	}

	return seqMatched || secondaryMatched
}

func (e *event) DeleteEventMatched(ctx *mirrorcontext.Context) bool {
	var val []byte

	if !e.isDeleted {
		return false
	}

	if e.ev.PrevKv != nil && len(e.ev.PrevKv.Value) > 0 && e.isDeleted {
		val = e.ev.PrevKv.Value
	}

	resource, err := kubeschema.Decode(ctx, e.ev.Kv.Key, val)
	if err != nil {
		return false
	}

	return resource.FilterSequentialByRules(ctx) || resource.FilterSecondaryByRules(ctx)
}
