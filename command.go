package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	WhiteSpacePattern string = " {2,}"
)

type command interface {
	Execute(f *file) error
}

type searchCommand struct {
	search     string
	replace    string
	matchCase  bool
	regexMode  bool
	whiteSpace bool
}

func (s searchCommand) Execute(f *file) error {
	var pattern string
	if s.regexMode {
		pattern = s.search
	} else {
		pattern = regexp.QuoteMeta(s.search)
	}
	if !s.matchCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("Invalid regex: '%v'", s.search)
	}

	fullName := f.getFullName()
	fullName = re.ReplaceAllString(fullName, s.replace)

	// untuk kasus pada penggunaan grup regex pada bagian replace
	// grup berdasarkan indeks %0, %1, dst
	// grup berdasarkan nama %nama yang didapat dari contoh (?P<nama>[a-zA-Z]+)
	if s.regexMode {
		matches := re.FindStringSubmatch(f.getFullName())
		// group %index
		for i, match := range matches {
			if match != "" {
				fullName = PercentReplace(fullName, strconv.Itoa(i), match)
			}
		}

		// group %name
		for i, name := range re.SubexpNames() {
			if name != "" {
				fullName = PercentReplace(fullName, name, matches[i])
			}
		}
	}

	if s.whiteSpace {
		f.setFullName(fullName)
	} else {
		whiteSpaceRegex := regexp.MustCompile(WhiteSpacePattern)
		fullName = whiteSpaceRegex.ReplaceAllString(fullName, " ")
		fullName = strings.Trim(fullName, " ")
		f.setFullName(fullName)
		f.name = strings.Trim(f.name, " ")
	}

	return nil
}

type deleteCommand struct {
	value      string
	matchCase  bool
	regexMode  bool
	whiteSpace bool
}

func (d deleteCommand) Execute(f *file) error {
	var pattern string
	if d.regexMode {
		pattern = d.value
	} else {
		pattern = regexp.QuoteMeta(d.value)
	}
	if !d.matchCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("Invalid regex: '%v'", d.value)
	}

	fullName := f.getFullName()
	fullName = re.ReplaceAllString(fullName, "")

	if d.whiteSpace {
		f.setFullName(fullName)
	} else {
		whiteSpaceRegex := regexp.MustCompile(WhiteSpacePattern)
		fullName = whiteSpaceRegex.ReplaceAllString(fullName, " ")
		fullName = strings.Trim(fullName, " ")
		f.setFullName(fullName)
		f.name = strings.Trim(f.name, " ")
	}

	return nil
}

type templateCommand struct {
	value string
}

func (t templateCommand) Execute(f *file) error {
	fullName := PercentReplace(t.value, "f", f.getFullName())
	fullName = PercentReplace(fullName, "n", f.name)
	fullName = PercentReplace(fullName, "x", f.getExt())
	f.setFullName(fullName)
	return nil
}

type moveCommand struct {
	pattern        string
	destinationDir string
	matchCase      bool
	regexMode      bool
}

func (m moveCommand) Execute(f *file) error {
	var pattern string
	if m.regexMode {
		pattern = m.pattern
	} else {
		pattern = regexp.QuoteMeta(m.pattern)
	}
	if !m.matchCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("Invalid regex: '%v'", m.pattern)
	}

	if re.MatchString(f.getFullName()) {
		f.baseDir = m.destinationDir
	}

	return nil
}

func PercentReplace(s string, name string, new string) string {
	re := regexp.MustCompile("%" + name)
	return re.ReplaceAllStringFunc(s, func(match string) string {
		idx := strings.Index(s, match)
		if idx > 0 && s[idx-1] == '%' {
			return strings.TrimPrefix(match, "%")
		}
		return new
	})
}
