package context

import (
	"github.com/etcd-carry/etcd-carry/pkg/testing/testoptions"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/storage/value"
	"path"
	"testing"
)

func TestNewMirrorContext(t *testing.T) {
	o := testoptions.GetMirrorOptions(t)
	ctx, err := NewMirrorContext(o)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := ctx.SourceClient.Get(ctx.Context, "foo")
	if err != nil || resp.Header.GetRevision() != 1 {
		t.Fatal(err)
	}
	resp, err = ctx.DestClient.Get(ctx.Context, "foo")
	if err != nil || resp.Header.GetRevision() != 1 {
		t.Fatal(err)
	}

	gr := schema.ParseGroupResource("secrets")
	key := path.Join(ctx.Etcd.KubePrefix, gr.Group, gr.Resource)
	if ctx.GetTransformer(key) == value.IdentityTransformer {
		t.Fatal("should not be IdentityTransformer")
	}

	gr = schema.ParseGroupResource("comfigmaps")
	key = path.Join(ctx.Etcd.KubePrefix, gr.Group, gr.Resource)
	if ctx.GetTransformer(key) != value.IdentityTransformer {
		t.Fatal("should be IdentityTransformer")
	}
}
