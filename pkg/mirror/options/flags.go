package options

import (
	"bytes"
	"fmt"
	"github.com/spf13/pflag"
	"io"
)

type SectionFlagSet struct {
	Sequence []string
	FlagSets map[string]*pflag.FlagSet
}

func (s *SectionFlagSet) FlagSet(section string) *pflag.FlagSet {
	if s.FlagSets == nil {
		s.FlagSets = map[string]*pflag.FlagSet{}
	}
	if _, ok := s.FlagSets[section]; !ok {
		s.FlagSets[section] = pflag.NewFlagSet(section, pflag.ExitOnError)
		s.Sequence = append(s.Sequence, section)
	}
	return s.FlagSets[section]
}

func PrintSections(w io.Writer, sfs SectionFlagSet) {
	for _, section := range sfs.Sequence {
		fs := sfs.FlagSets[section]
		if !fs.HasFlags() {
			continue
		}

		var buf bytes.Buffer
		fmt.Fprintf(&buf, "\n%s:\n\n%s", section, fs.FlagUsages())
		fmt.Fprint(w, buf.String())
	}
}
