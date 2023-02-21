package context

import (
	"context"
	mirrorclientv3 "github.com/etcd-carry/etcd-carry/pkg/client/v3"
	"github.com/etcd-carry/etcd-carry/pkg/filter/kube/layer2"
	"github.com/etcd-carry/etcd-carry/pkg/mirror/options"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/options/encryptionconfig"
	"k8s.io/apiserver/pkg/storage/value"
	"os"
	"path"
	"strings"
)

type Context struct {
	context.Context
	options.MirrorOptions
	DestClient   *clientv3.Client
	SourceClient *clientv3.Client
	MirrorFilter *layer2.Filter
	Transformers map[string]value.Transformer
}

func NewMirrorContext(o *options.MirrorOptions) (*Context, error) {
	var err error
	mirrorCtx := &Context{
		Context:       context.Background(),
		MirrorOptions: *o,
		Transformers:  make(map[string]value.Transformer),
	}

	if o.Generic.Debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 4))
	} else {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, os.Stderr))
	}

	mirrorCtx.DestClient, err = mirrorclientv3.New(mirrorclientv3.ConfigSpec{
		Endpoints:        mirrorCtx.MirrorOptions.Transport.DestTransport.ServerList,
		DialTimeout:      mirrorCtx.MirrorOptions.Transport.DialTimeout,
		KeepAliveTime:    mirrorCtx.MirrorOptions.Transport.KeepAliveTime,
		KeepAliveTimeout: mirrorCtx.MirrorOptions.Transport.KeepAliveTimeout,
		Secure: &mirrorclientv3.SecureCfg{
			Cert:               mirrorCtx.MirrorOptions.Transport.DestTransport.CertFile,
			Key:                mirrorCtx.MirrorOptions.Transport.DestTransport.KeyFile,
			Cacert:             mirrorCtx.MirrorOptions.Transport.DestTransport.CACertFile,
			InsecureTransport:  mirrorCtx.MirrorOptions.Transport.DestTransport.Insecure,
			InsecureSkipVerify: mirrorCtx.MirrorOptions.Transport.DestTransport.InsecureSkipVerify,
		}})
	if err != nil {
		return nil, err
	}
	mirrorCtx.SourceClient, err = mirrorclientv3.New(mirrorclientv3.ConfigSpec{
		Endpoints:        mirrorCtx.MirrorOptions.Transport.SourceTransport.ServerList,
		DialTimeout:      mirrorCtx.MirrorOptions.Transport.DialTimeout,
		DialOptions:      []grpc.DialOption{grpc.WithBlock()},
		KeepAliveTime:    mirrorCtx.MirrorOptions.Transport.KeepAliveTime,
		KeepAliveTimeout: mirrorCtx.MirrorOptions.Transport.KeepAliveTimeout,
		Secure: &mirrorclientv3.SecureCfg{
			Cert:               mirrorCtx.MirrorOptions.Transport.SourceTransport.CertFile,
			Key:                mirrorCtx.MirrorOptions.Transport.SourceTransport.KeyFile,
			Cacert:             mirrorCtx.MirrorOptions.Transport.SourceTransport.CACertFile,
			InsecureTransport:  mirrorCtx.MirrorOptions.Transport.SourceTransport.Insecure,
			InsecureSkipVerify: mirrorCtx.MirrorOptions.Transport.SourceTransport.InsecureSkipVerify,
		}})
	if err != nil {
		return nil, err
	}

	mirrorCtx.MirrorFilter, err = layer2.NewFilter(mirrorCtx.MirrorOptions.Generic.MirrorRulesConfigFilepath)
	if err != nil {
		return nil, err
	}

	var transformers map[schema.GroupResource]value.Transformer
	transformers, err = encryptionconfig.GetTransformerOverrides(mirrorCtx.MirrorOptions.Etcd.EncryptionProviderConfigFilepath)
	if err != nil {
		return nil, err
	}
	for gr, transformer := range transformers {
		mirrorCtx.SetTransformer(gr, transformer)
	}

	return mirrorCtx, nil
}

func (ctx *Context) SetTransformer(gr schema.GroupResource, transformer value.Transformer) {
	key := path.Join(ctx.Etcd.KubePrefix, gr.Group, gr.Resource)
	ctx.Transformers[key] = transformer
}

func (ctx *Context) GetTransformer(key string) value.Transformer {
	for prefix, transformer := range ctx.Transformers {
		if strings.HasPrefix(key, prefix) {
			return transformer
		}
	}
	return value.IdentityTransformer
}
