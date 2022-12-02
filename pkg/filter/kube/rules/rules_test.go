package rules

import (
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"path"
	"reflect"
	"testing"
)

const ValidMirrorRules = `
filters:
  sequential:
    - group: ""
      sequence: 1
      resources:
        - version: v1
          kind: Namespace
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: In
              values:
                - unit
  secondary:
    - group: ""
      resources:
        - version: v1
          kind: ConfigMap
        - version: v1
          Kind: Secret
      namespace: unit-test
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: Exists
      excludes:
        - resource:
            version: v1
            kind: Secret
          name: exclude-me-secret
          namespace: unit-test
`
const ValidOnlySequentialMirrorRules = `
filters:
  sequential:
    - group: ""
      sequence: 1
      resources:
        - version: v1
          kind: Namespace
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: In
              values:
                - unit
`
const ValidOnlySecondaryMirrorRules = `
filters:
  secondary:
    - group: ""
      resources:
        - version: v1
          kind: ConfigMap
        - version: v1
          Kind: Secret
      namespace: unit-test
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: Exists
      excludes:
        - resource:
            version: v1
            kind: Secret
          name: exclude-me-secret
          namespace: unit-test
`
const InvalidMirrorRules = `
filters:
  primary:
    - group: ""
      sequence: 1
      resources:
        - version: v1
          kind: Namespace
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: In
              values:
                - unit
  secondary:
    - group: ""
      resources:
        version: v1
        kind: ConfigMap
      namespace: unit-test
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: Exists
`
const InvalidSequentialMirrorRules = `
filters:
  sequential:
    - group: ""
      sequence: 1
      resources:
        version: v1
        kind: Namespace
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: In
              values:
                - unit
  secondary:
    - group: ""
      resources:
        - version: v1
          kind: ConfigMap
        - version: v1
          Kind: Secret
      namespace: unit-test
`
const InvalidSecondaryMirrorRules = `
filters:
  sequential:
    - group: ""
      sequence: 1
      resources:
        - version: v1
          kind: Namespace
  secondary:
    - group: ""
      resources:
        - version: v1
          kind: ConfigMap
        - version: v1
          Kind: Secret
      namespace: unit-test
      labelSelectors:
        - matchExpressions:
            - key: test.io/namespace-kind
              operator: Exists
      excludes:
          resource:
            version: v1
            kind: Secret
          name: exclude-me-secret
          namespace: unit-test
`

func TestLoadMirrorRules(t *testing.T) {
	testCase := []struct {
		name         string
		rulesContent string
		expected     *Rules
		expectErrors bool
	}{
		{
			name:         "valid mirror rules",
			rulesContent: ValidMirrorRules,
			expected: &Rules{Filters: struct {
				Sequential []Sequential `yaml:"sequential"`
				Secondary  []Secondary  `yaml:"secondary"`
			}(struct {
				Sequential []Sequential
				Secondary  []Secondary
			}{
				Sequential: []Sequential{
					{
						Group:    "",
						Sequence: 1,
						Resources: []metav1.GroupVersionKind{
							{
								Version: "v1",
								Kind:    "Namespace",
							},
						},
						LabelSelectors: []metav1.LabelSelector{
							{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "test.io/namespace-kind",
										Operator: "In",
										Values:   []string{"unit"},
									},
								},
							},
						},
					},
				},
				Secondary: []Secondary{
					{
						Group: "",
						Resources: []metav1.GroupVersionKind{
							{
								Version: "v1",
								Kind:    "ConfigMap",
							},
							{
								Version: "v1",
								Kind:    "Secret",
							},
						},
						Namespace: "unit-test",
						LabelSelectors: []metav1.LabelSelector{
							{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "test.io/namespace-kind",
										Operator: "Exists",
									},
								},
							},
						},
						Excludes: []Exclude{
							{
								Resource:  schema.GroupVersionKind{Version: "v1", Kind: "Secret"},
								Name:      "exclude-me-secret",
								Namespace: "unit-test",
							},
						},
					},
				},
			})},
			expectErrors: false,
		},
		{
			name:         "valid only sequential rules",
			rulesContent: ValidOnlySequentialMirrorRules,
			expected: &Rules{Filters: struct {
				Sequential []Sequential `yaml:"sequential"`
				Secondary  []Secondary  `yaml:"secondary"`
			}(struct {
				Sequential []Sequential
				Secondary  []Secondary
			}{
				Sequential: []Sequential{
					{
						Group:    "",
						Sequence: 1,
						Resources: []metav1.GroupVersionKind{
							{
								Version: "v1",
								Kind:    "Namespace",
							},
						},
						LabelSelectors: []metav1.LabelSelector{
							{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "test.io/namespace-kind",
										Operator: "In",
										Values:   []string{"unit"},
									},
								},
							},
						},
					},
				},
				Secondary: nil,
			})},
			expectErrors: false,
		},
		{
			name:         "valid only secondary rules",
			rulesContent: ValidOnlySecondaryMirrorRules,
			expected: &Rules{Filters: struct {
				Sequential []Sequential `yaml:"sequential"`
				Secondary  []Secondary  `yaml:"secondary"`
			}(struct {
				Sequential []Sequential
				Secondary  []Secondary
			}{
				Sequential: nil,
				Secondary: []Secondary{
					{
						Group: "",
						Resources: []metav1.GroupVersionKind{
							{
								Version: "v1",
								Kind:    "ConfigMap",
							},
							{
								Version: "v1",
								Kind:    "Secret",
							},
						},
						Namespace: "unit-test",
						LabelSelectors: []metav1.LabelSelector{
							{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									{
										Key:      "test.io/namespace-kind",
										Operator: "Exists",
									},
								},
							},
						},
						Excludes: []Exclude{
							{
								Resource: schema.GroupVersionKind{
									Version: "v1",
									Kind:    "Secret",
								},
								Name:      "exclude-me-secret",
								Namespace: "unit-test",
							},
						},
					},
				},
			})},
			expectErrors: false,
		},
		{
			name:         "invalid mirror rules",
			rulesContent: InvalidMirrorRules,
			expectErrors: true,
		},
		{
			name:         "invalid sequential rules",
			rulesContent: InvalidSequentialMirrorRules,
			expectErrors: true,
		},
		{
			name:         "invalid secondary rules",
			rulesContent: InvalidSecondaryMirrorRules,
			expectErrors: true,
		},
	}

	baseDir := os.TempDir()
	tempDir, err := ioutil.TempDir(baseDir, "etcd-rules")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rulesFile := path.Join(tempDir, "rules.yaml")
			if err = ioutil.WriteFile(rulesFile, []byte(tc.rulesContent), 0644); err != nil {
				t.Error(err)
			}

			rules, err := LoadMirrorRules(rulesFile)
			if err != nil && !tc.expectErrors {
				t.Errorf("expected no errors, but error found %+v", err)
			}
			if err == nil && tc.expectErrors {
				t.Error("expected errors, but no errors found")
			}
			if err == nil && !tc.expectErrors {
				if !reflect.DeepEqual(tc.expected, rules) {
					t.Errorf("Difference detected on:\\n%s", cmp.Diff(tc.expected, rules))
				}
			}
		})
	}
}
