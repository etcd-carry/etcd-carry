package schema

import (
	"fmt"
	"github.com/etcd-carry/etcd-carry/pkg/constant"
	"github.com/etcd-carry/etcd-carry/pkg/filter/kube/layer2"
	"github.com/etcd-carry/etcd-carry/pkg/filter/kube/rules"
	mirrorcontext "github.com/etcd-carry/etcd-carry/pkg/mirror/context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type KubeResource struct {
	object *unstructured.Unstructured
}

func (k *KubeResource) GVK() schema.GroupVersionKind {
	return k.object.GroupVersionKind()
}

func (k *KubeResource) Object() *unstructured.Unstructured {
	return k.object
}

func (k *KubeResource) FilterSequentialByRules(ctx *mirrorcontext.Context) bool {
	for _, rule := range ctx.MirrorFilter.Rules.Filters.Sequential {
		// 1. match group
		// 2. match gvk
		// 3. match selector
		if k.groupMatches(rule.Group) &&
			k.gvkMatches(rule.Resources) &&
			k.selectorMatches(rule.LabelSelectors) {
			fmt.Println("[Matched Sequential]:", "GVK:", k.GVK().String(), "Name:", k.Object().GetName(),
				"Namespace:", k.Object().GetNamespace(), "Selector:", k.Object().GetLabels())
			return true
		}
	}

	return false
}

func (k *KubeResource) FilterSecondaryByRules(ctx *mirrorcontext.Context) bool {
	for _, rule := range ctx.MirrorFilter.Rules.Filters.Secondary {
		// 1. match group
		// 2. match gvk
		// 3. match namespace
		// 4. match selector
		// 5. match excluded
		// 6. match field
		if k.groupMatches(rule.Group) &&
			k.gvkMatches(rule.Resources) &&
			k.namespaceMatches(rule.Namespace, rule.NamespaceSelectors, ctx.MirrorFilter) &&
			k.selectorMatches(rule.LabelSelectors) &&
			!k.exclusionMatches(rule.Excludes) &&
			k.fieldMatches(rule.FieldSelectors) {
			fmt.Println("[Matched Object]:", "GVK:", k.GVK().String(), "Name:", k.Object().GetName(),
				"Namespace:", k.Object().GetNamespace(), "Selector:", k.Object().GetLabels())
			fmt.Println("[Matched Rule]:", "Group:", rule.Group, "Resource:", rule.Resources, "Namespace:", rule.Namespace,
				"Selector:", rule.LabelSelectors, "Exclude:", rule.Excludes, "Field:", rule.FieldSelectors)
			return true
		}
	}
	return false
}

func (k *KubeResource) groupMatches(ruleGroup string) bool {
	switch {
	// allow all groups
	case ruleGroup == "*" || k.GVK().Group == ruleGroup:
		return true
	case (k.GVK().Group == "" || ruleGroup == "") && k.GVK().Group != ruleGroup:
		return false
	default:
		return layer2.IsResourceGroupRegexMatch(k.GVK().Group, ruleGroup)
	}
}

func (k *KubeResource) gvkMatches(ruleResource []metav1.GroupVersionKind) bool {
	if len(ruleResource) != 0 {
		for _, r := range ruleResource {
			if r.String() == k.GVK().String() {
				return true
			}
		}
		return false
	}
	return true
}

func (k *KubeResource) selectorMatches(ruleSelector []metav1.LabelSelector) bool {
	if len(ruleSelector) != 0 {
		for _, s := range ruleSelector {
			if layer2.IsSelectorLooseOverlap(k.Object().GetLabels(), &s) {
				return true
			}
		}
		return false
	}
	return true
}

func (k *KubeResource) namespaceMatches(ruleNamespace string, ruleNamespaceSelector []metav1.LabelSelector, filter *layer2.Filter) bool {
	switch {
	case ruleNamespace != "":
		return k.Object().GetNamespace() == ruleNamespace
	case len(ruleNamespaceSelector) != 0:
		if labels, ok := filter.Namespace[k.Object().GetNamespace()]; ok {
			for _, s := range ruleNamespaceSelector {
				if layer2.IsSelectorLooseOverlap(labels, &s) {
					return true
				}
			}
			return false
		}
		return false
	default:
		return true
	}
}

func (k *KubeResource) exclusionMatches(ruleExclude []rules.Exclude) bool {
	if len(ruleExclude) != 0 {
		for _, e := range ruleExclude {
			if e.Resource.String() == k.GVK().String() {
				switch {
				// cluster scope object
				case e.Name != "" && e.Namespace == "" && k.Object().GetName() == e.Name:
					return true
				// namespace scope object
				case e.Name != "" && e.Namespace != "" && k.Object().GetName() == e.Name && k.Object().GetNamespace() == e.Namespace:
					return true
				// match expressions object
				case e.Name == "" && len(e.LabelSelectors) != 0:
					// cluster scope
					// namespace scope
					if e.Namespace == "" {
						for _, s := range e.LabelSelectors {
							if layer2.IsSelectorLooseOverlap(k.Object().GetLabels(), &s) {
								return true
							}
						}
					} else {
						for _, s := range e.LabelSelectors {
							if layer2.IsSelectorLooseOverlap(k.Object().GetLabels(), &s) && k.Object().GetNamespace() == e.Namespace {
								return true
							}
						}
					}
				// namespace scope resource
				case e.Name == "" && len(e.LabelSelectors) == 0 && e.Namespace != "":
					return k.Object().GetNamespace() == e.Namespace
				// cluster scope resource
				case e.Name == "" && len(e.LabelSelectors) == 0 && e.Namespace == "":
					return true
				}
			}
		}
		return false
	}
	return false
}

func (k *KubeResource) fieldMatches(fieldSelector []metav1.LabelSelector) bool {
	if len(fieldSelector) != 0 {
		fields := layer2.FilterStringToMap(k.Object().UnstructuredContent())
		for _, s := range fieldSelector {
			if layer2.IsSelectorLooseOverlap(fields, &s) {
				return true
			}
		}
		return false
	}
	return true
}

func (k *KubeResource) MutateKubeResource(kvs *mvccpb.KeyValue) (*mvccpb.KeyValue, error) {
	var out = kvs

	switch k.GVK() {
	case corev1.SchemeGroupVersion.WithKind(constant.KubeKindService):
		var svc = &corev1.Service{}

		delete(k.Object().Object["spec"].(map[string]interface{}), "clusterIP")
		if err := runtime.DefaultUnstructuredConverter.FromUnstructured(k.Object().UnstructuredContent(), svc); err != nil {
			fmt.Println("[Error] Convert to Service error,", "Object:", k.Object().UnstructuredContent(), "error:", err)
			return nil, err
		}

		data, err := runtime.Encode(Codecs.LegacyCodec(corev1.SchemeGroupVersion), runtime.Object(svc))
		if err != nil {
			fmt.Println("[Error] Encode Service data error,", "GVK:", k.GVK().String(), "name:", svc.Name, "namespace:", svc.Namespace, "error:", err)
			return nil, err
		}

		out.Value = data
		return out, nil
	default:
		return out, nil
	}
}
