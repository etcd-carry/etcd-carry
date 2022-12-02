package layer2

import (
	"github.com/etcd-carry/etcd-carry/pkg/filter/kube/rules"
	"reflect"
)

type Filter struct {
	Namespace map[string]map[string]string
	Rules     *rules.Rules
}

func NewFilter(ruleConfigPath string) (*Filter, error) {
	rules, err := rules.LoadMirrorRules(ruleConfigPath)
	if err != nil {
		return nil, err
	}
	return &Filter{
		Namespace: make(map[string]map[string]string),
		Rules:     rules,
	}, nil
}

func FilterStringToMap(src map[string]interface{}) map[string]string {
	if src == nil {
		return nil
	}

	dest := make(map[string]string)
	for k, v := range src {
		kt := reflect.TypeOf(k)
		vt := reflect.TypeOf(v)
		if kt.Kind() == reflect.String && vt.Kind() == reflect.String {
			dest[k] = v.(string)
		}
	}
	return dest
}
