package mirror

import (
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/constant"
	kubeschema "github.com/etcd-carry/etcd-carry/pkg/filter/kube/schema"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

const (
	batchLimit = 1000
)

type Syncer interface {
	SyncSequential() (<-chan clientv3.GetResponse, chan error)
	SyncSecondary() (<-chan clientv3.GetResponse, chan error)
	SyncUpdates() clientv3.WatchChan
}

func NewSyncer(ctx *mirrorcontext.Context) Syncer {
	return &syncer{Context: ctx}
}

type syncer struct {
	*mirrorcontext.Context
}

func (s *syncer) SyncSequential() (<-chan clientv3.GetResponse, chan error) {
	var prefixNum = 0
	for _, rule := range s.MirrorFilter.Rules.Filters.Sequential {
		prefixNum = prefixNum + len(rule.Resources)
	}

	respChan := make(chan clientv3.GetResponse, prefixNum)
	errChan := make(chan error, 1)

	if s.Etcd.StartReversion == 0 {
		resp, err := s.MasterClient.Get(s.Context, constant.KubeRootPrefix)
		if err != nil {
			errChan <- err
			defer close(respChan)
			defer close(errChan)
			return respChan, errChan
		}
		s.Etcd.StartReversion = resp.Header.Revision
	}

	go func() {
		defer close(respChan)
		defer close(errChan)

		for _, rule := range s.MirrorFilter.Rules.Filters.Sequential {
			var keys []string

			for _, rs := range rule.Resources {
				var key = constant.KubeRootPrefix
				if rs.Group != "" {
					key = key + strings.ToLower(rs.Group) + "/"
				}
				key = key + strings.ToLower(rs.Kind+"s") + "/"
				keys = append(keys, key)
			}

			for _, key := range keys {
				var kvs []*mvccpb.KeyValue
				opts := []clientv3.OpOption{clientv3.WithRev(s.Etcd.StartReversion), clientv3.WithRange(clientv3.GetPrefixRangeEnd(key))}

				resp, err := s.MasterClient.Get(s.Context, key, opts...)
				if err != nil {
					errChan <- err
					return
				}

				for _, kv := range resp.Kvs {
					resource, err := kubeschema.Decode(s.Context, kv.Key, kv.Value)
					if err != nil {
						errChan <- err
						return
					}

					if resource.FilterSequentialByRules(s.Context) {
						kvs = append(kvs, kv)

						// TODO:
						if corev1.SchemeGroupVersion.WithKind(constant.KubeKindNamespace) == resource.GVK() {
							s.MirrorFilter.Namespace[resource.Object().GetName()] = resource.Object().GetLabels()
						}
					}
				}

				fmt.Printf("Key: %v, Matched: %v, Filtered: %v\n", key, resp.Count, len(kvs))
				resp.Kvs = kvs
				resp.Count = int64(len(kvs))
				respChan <- *resp
			}
		}
	}()

	return respChan, errChan
}

func (s *syncer) SyncSecondary() (<-chan clientv3.GetResponse, chan error) {
	if s.Etcd.StartReversion == 0 {
		panic("unexpect revision = 0. Calling SyncBase before SyncNamespace finishes?")
	}

	var respChan = make(chan clientv3.GetResponse, 1024)
	var errChan = make(chan error, 1)

	go func() {
		defer close(respChan)
		defer close(errChan)

		var key = constant.KubeRootPrefix

		opts := []clientv3.OpOption{clientv3.WithLimit(batchLimit), clientv3.WithRev(s.Etcd.StartReversion), clientv3.WithRange(clientv3.GetPrefixRangeEnd(key))}

		for {
			var kvs []*mvccpb.KeyValue

			resp, err := s.MasterClient.Get(s.Context, key, opts...)
			if err != nil {
				errChan <- err
				return
			}

			for _, kv := range resp.Kvs {
				resource, err := kubeschema.Decode(s.Context, kv.Key, kv.Value)
				if err != nil {
					errChan <- err
					return
				}

				if resource.FilterSecondaryByRules(s.Context) {
					out, err := resource.MutateKubeResource(kv)
					if err != nil {
						errChan <- err
						return
					}

					kvs = append(kvs, out)
				}
			}

			fmt.Printf("Key: %v, Matched: %v, Filtered: %v\n", key, resp.Count, len(kvs))

			if len(kvs) != 0 {
				out := resp
				out.Count = int64(len(kvs))
				out.Kvs = kvs
				respChan <- *resp
			}

			if !resp.More {
				fmt.Println("No more keys matched: ", key, resp.Count)
				return
			}
			key = string(append(resp.Kvs[len(resp.Kvs)-1].Key, 0))
		}
	}()

	return respChan, errChan
}

func (s *syncer) SyncUpdates() clientv3.WatchChan {
	if s.Etcd.StartReversion == 0 {
		panic("unexpect revision = 0. Calling SyncUpdates before SyncBase finishes?")
	}

	fmt.Println("Start to watch updates from revision", s.Etcd.StartReversion)

	return s.MasterClient.Watch(s.Context, constant.KubeRootPrefix, clientv3.WithPrefix(), clientv3.WithRev(s.Etcd.StartReversion+1), clientv3.WithPrevKV())
}
