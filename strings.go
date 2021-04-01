package goutil

func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

type Strings []string

func (s Strings) Contains(v string) bool {
	if len(s) == 0 {
		return false
	}
	for _, _v := range s {
		if _v == v {
			return true
		}
	}
	return false
}
