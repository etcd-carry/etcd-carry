package context

import (
	"context"
	mirrorgrpcprobe "github.com/etcd-carry/etcd-carry/pkg/client/probe/grpc"
	mirrorclientv3 "github.com/etcd-carry/etcd-carry/pkg/client/v3"
	"github.com/etcd-carry/etcd-carry/pkg/filter/kube/layer2"
	"github.com/etcd-carry/etcd-carry/pkg/mirror/options"
	"github.com/etcd-carry/etcd-carry/pkg/rest"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/options/encryptionconfig"
	"k8s.io/apiserver/pkg/storage/value"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

type Context struct {
	context.Context
	options.MirrorOptions
	Probe          mirrorgrpcprobe.Probe
	SlaveClient    *clientv3.Client
	MasterClient   *clientv3.Client
	MirrorFilter   *layer2.Filter
	Transformers   map[string]value.Transformer
	RestfulServing rest.RestfulServing
}

func NewMirrorContext(o *options.MirrorOptions) (*Context, error) {
	var err error
	mirrorCtx := &Context{
		Context:       context.Background(),
		MirrorOptions: *o,
		Probe:         mirrorgrpcprobe.New(),
		Transformers:  make(map[string]value.Transformer),
	}

	if o.Generic.Debug {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 4))
	} else {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, os.Stderr))
	}

	mirrorCtx.SlaveClient, err = mirrorclientv3.New(mirrorclientv3.ConfigSpec{
		Endpoints:        mirrorCtx.MirrorOptions.Transport.SlaveTransport.ServerList,
		DialTimeout:      mirrorCtx.MirrorOptions.Transport.DialTimeout,
		KeepAliveTime:    mirrorCtx.MirrorOptions.Transport.KeepAliveTime,
		KeepAliveTimeout: mirrorCtx.MirrorOptions.Transport.KeepAliveTimeout,
		Secure: &mirrorclientv3.SecureCfg{
			Cert:               mirrorCtx.MirrorOptions.Transport.SlaveTransport.CertFile,
			Key:                mirrorCtx.MirrorOptions.Transport.SlaveTransport.KeyFile,
			Cacert:             mirrorCtx.MirrorOptions.Transport.SlaveTransport.CACertFile,
			InsecureTransport:  mirrorCtx.MirrorOptions.Transport.SlaveTransport.Insecure,
			InsecureSkipVerify: mirrorCtx.MirrorOptions.Transport.SlaveTransport.InsecureSkipVerify,
		}})
	if err != nil {
		return nil, err
	}
	mirrorCtx.MasterClient, err = mirrorclientv3.New(mirrorclientv3.ConfigSpec{
		Endpoints:        mirrorCtx.MirrorOptions.Transport.MasterTransport.ServerList,
		DialTimeout:      mirrorCtx.MirrorOptions.Transport.DialTimeout,
		DialOptions:      []grpc.DialOption{grpc.WithBlock()},
		KeepAliveTime:    mirrorCtx.MirrorOptions.Transport.KeepAliveTime,
		KeepAliveTimeout: mirrorCtx.MirrorOptions.Transport.KeepAliveTimeout,
		Secure: &mirrorclientv3.SecureCfg{
			Cert:               mirrorCtx.MirrorOptions.Transport.MasterTransport.CertFile,
			Key:                mirrorCtx.MirrorOptions.Transport.MasterTransport.KeyFile,
			Cacert:             mirrorCtx.MirrorOptions.Transport.MasterTransport.CACertFile,
			InsecureTransport:  mirrorCtx.MirrorOptions.Transport.MasterTransport.Insecure,
			InsecureSkipVerify: mirrorCtx.MirrorOptions.Transport.MasterTransport.InsecureSkipVerify,
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

	mirrorCtx.RestfulServing, err = rest.NewRestfulServing(net.JoinHostPort(mirrorCtx.MirrorOptions.RestfulServing.BindAddress.String(), strconv.Itoa(mirrorCtx.MirrorOptions.RestfulServing.BindPort)))
	if err != nil {
		return nil, err
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
