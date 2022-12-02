package layer2

import (
	"bytes"
	"regexp"
	"strings"
)

var match = regexp.MustCompile

func convertKubeGroupRegex(src string) string {
	var buffer bytes.Buffer

	s := strings.ReplaceAll(src, ".", "\\.")
	buffer.WriteString(strings.ReplaceAll(s, "*", "[[:alnum:]]*"))
	buffer.WriteString("$")

	return buffer.String()
}

func IsResourceGroupRegexMatch(group, ruleGroup string) bool {
	keyExp := convertKubeGroupRegex(ruleGroup)
	p := match(keyExp)
	return p.FindString(group) == group
}
