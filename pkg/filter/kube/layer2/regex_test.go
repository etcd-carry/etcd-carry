package layer2

import "testing"

func TestIsResourceGroupRegexMatch(t *testing.T) {
	testCase := []struct {
		name      string
		group     string
		ruleGroup string
		isMatched bool
	}{
		{
			name:      "match specific group 2",
			group:     "test.io",
			ruleGroup: "test.io",
			isMatched: true,
		},
		{
			name:      "match specific group 3",
			group:     "*",
			ruleGroup: "test.io",
			isMatched: false,
		},
		{
			name:      "match specific group 4",
			group:     "repo.test.io",
			ruleGroup: "test.io",
			isMatched: false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			if IsResourceGroupRegexMatch(tc.group, tc.ruleGroup) != tc.isMatched {
				t.Errorf("should not happened, match group %+v, rule group %+v", tc.group, tc.ruleGroup)
			}
		})
	}
}
