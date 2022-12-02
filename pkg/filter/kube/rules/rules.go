package rules

import (
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Rules struct {
	Filters struct {
		Sequential []Sequential `yaml:"sequential"`
		Secondary  []Secondary  `yaml:"secondary"`
	} `yaml:"filters"`
}

type Sequential struct {
	Group          string                    `yaml:"group"`
	Sequence       int                       `yaml:"sequence"`
	Resources      []metav1.GroupVersionKind `yaml:"resources"`
	LabelSelectors []metav1.LabelSelector    `yaml:"labelSelectors"`
}

type Secondary struct {
	Group              string                    `yaml:"group"`
	Resources          []metav1.GroupVersionKind `yaml:"resources"`
	Excludes           []Exclude                 `yaml:"excludes"`
	Namespace          string                    `yaml:"namespace"`
	NamespaceSelectors []metav1.LabelSelector    `yaml:"namespaceSelectors"`
	LabelSelectors     []metav1.LabelSelector    `yaml:"labelSelectors"`
	FieldSelectors     []metav1.LabelSelector    `yaml:"fieldSelectors"`
}

type Exclude struct {
	Resource       schema.GroupVersionKind `yaml:"resource"`
	Name           string                  `yaml:"name"`
	Namespace      string                  `yaml:"namespace"`
	LabelSelectors []metav1.LabelSelector  `yaml:"labelSelectors"`
}

func LoadMirrorRules(path string) (*Rules, error) {
	var rules Rules
	if path == "" {
		return nil, fmt.Errorf("mirror rules not defined")
	}

	yamlData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(yamlData, &rules); err != nil {
		fmt.Println("unmarshal error: ", err)
		return nil, err
	}
	return &rules, nil
}
