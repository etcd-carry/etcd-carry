package schema

import (
	"fmt"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage/value"
)

type authenticatedDataString string

// AuthenticatedData implements the value.Context interface.
func (d authenticatedDataString) AuthenticatedData() []byte {
	return []byte(d)
}

var _ value.Context = authenticatedDataString("")

func Decode(ctx *mirrorcontext.Context, key, value []byte) (out *KubeResource, err error) {
	var unknown runtime.Unknown
	var decoder = Codecs.UniversalDeserializer()

	decryptedValue, _, err := ctx.GetTransformer(string(key)).TransformFromStorage(value, authenticatedDataString(key))
	if err != nil {
		fmt.Println("transform from storage error: ", string(key), err)
		return nil, err
	}

	if _, _, err := decoder.Decode(decryptedValue, nil, &unknown); err != nil {
		fmt.Println("decode into unknown object error:", string(key), err)
		return nil, err
	}

	unst := unstructured.Unstructured{}
	gvk := unknown.GroupVersionKind()
	obj, err := Scheme.New(gvk)
	switch {
	case runtime.IsNotRegisteredError(err):
		if _, _, err := decoder.Decode(decryptedValue, &gvk, &unst); err != nil {
			fmt.Println("decode into object error:", string(key), err)
			return nil, err
		}
		return &KubeResource{object: &unst}, nil
	case err != nil:
		fmt.Println("Couldn't create external object:", gvk, err)
		return nil, err
	default:
		if _, _, err := decoder.Decode(decryptedValue, &gvk, obj); err != nil {
			fmt.Println("decode into object error:", string(key), err)
			return nil, err
		}
	}

	unstrBody, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		fmt.Println("ToUnstructured failed:", err)
		return nil, err
	}
	unst.Object = unstrBody

	return &KubeResource{object: &unst}, nil
}
