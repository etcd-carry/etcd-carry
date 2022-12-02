package layer2

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func convertSelectorToMatchExpressions(selector *metav1.LabelSelector) map[string]metav1.LabelSelectorRequirement {
	matchExps := map[string]metav1.LabelSelectorRequirement{}
	for _, exp := range selector.MatchExpressions {
		matchExps[exp.Key] = exp
	}

	for k, v := range selector.MatchLabels {
		matchExps[k] = metav1.LabelSelectorRequirement{
			Operator: metav1.LabelSelectorOpIn,
			Values:   []string{v},
		}
	}

	return matchExps
}

func isMatchExpOverlap(value string, matchExp metav1.LabelSelectorRequirement) bool {
	switch matchExp.Operator {
	case metav1.LabelSelectorOpIn:
		if sliceOverlaps(value, matchExp.Values) {
			return true
		}
	case metav1.LabelSelectorOpExists:
		return true
	case metav1.LabelSelectorOpNotIn:
		if !sliceOverlaps(value, matchExp.Values) {
			return true
		}
	}
	return false
}

// IsSelectorLooseOverlap
// labels: orcastack.io/namespace-kind: tenant
// labels: orcastack.io/namespace-kind: ss
// labels: orcastack.io/namespace-kind: rss
// matchExpressions:
//   - key: orcastack.io/namespace-kind
//     operator: In
//     values:
//   - tenant
//   - rss
func IsSelectorLooseOverlap(labels map[string]string, selector *metav1.LabelSelector) bool {
	if selector == nil {
		return true
	}

	if labels == nil {
		return false
	}

	matchExps := convertSelectorToMatchExpressions(selector)
	for k, exp := range matchExps {
		value, ok := labels[k]
		if !ok {
			return false
		}
		if !isMatchExpOverlap(value, exp) {
			return false
		}
	}

	return true
}

func sliceOverlaps(a string, b []string) bool {
	keyExist := make(map[string]bool, len(b))
	for _, key := range b {
		keyExist[key] = true
	}

	return keyExist[a]
}
