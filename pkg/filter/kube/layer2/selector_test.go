package layer2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestIsSelectorLooseOverlap(t *testing.T) {
	testCase := []struct {
		name      string
		label     map[string]string
		selector  *metav1.LabelSelector
		isMatched bool
	}{
		{
			name:  "key is in and value matches",
			label: map[string]string{"test.io/namespace-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{"A"},
					},
				},
			},
			isMatched: true,
		},
		{
			name:  "key is in and matches one of the values",
			label: map[string]string{"test.io/namespace-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{"A", "B", "C"},
					},
				},
			},
			isMatched: true,
		},
		{
			name:  "key is not in and mismatches value",
			label: map[string]string{"test.io/namespace-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpIn,
						Values:   []string{"B", "C"},
					},
				},
			},
			isMatched: false,
		},
		{
			name:  "key exists and matches",
			label: map[string]string{"test.io/namespace-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			},
			isMatched: true,
		},
		{
			name:  "key does not exist and mismatches",
			label: map[string]string{"test.io/service-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			},
			isMatched: false,
		},
		{
			name:  "key is not in and matches",
			label: map[string]string{"test.io/namespace-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpNotIn,
						Values:   []string{"B", "C"},
					},
				},
			},
			isMatched: true,
		},
		{
			name:  "key is in and mismatches",
			label: map[string]string{"test.io/namespace-kind": "A"},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "test.io/namespace-kind",
						Operator: metav1.LabelSelectorOpNotIn,
						Values:   []string{"A", "B", "C"},
					},
				},
			},
			isMatched: false,
		},
		//{
		//	name:  "key dose not exist and matches",
		//	label: map[string]string{"test.io/service-kind": "A"},
		//	selector: &metav1.LabelSelector{
		//		MatchExpressions: []metav1.LabelSelectorRequirement{
		//			{
		//				Key:      "test.io/namespace-kind",
		//				Operator: metav1.LabelSelectorOpDoesNotExist,
		//			},
		//		},
		//	},
		//	isMatched: true,
		//},
		//{
		//	name:  "key dose not exist and mismatches",
		//	label: map[string]string{"test.io/namespace-kind": "A"},
		//	selector: &metav1.LabelSelector{
		//		MatchExpressions: []metav1.LabelSelectorRequirement{
		//			{
		//				Key:      "test.io/namespace-kind",
		//				Operator: metav1.LabelSelectorOpDoesNotExist,
		//			},
		//		},
		//	},
		//	isMatched: false,
		//},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			if IsSelectorLooseOverlap(tc.label, tc.selector) != tc.isMatched {
				t.Errorf("should not happened, lable %+v, selector %+v", tc.label, tc.selector.MatchExpressions[0])
			}
		})
	}
}
