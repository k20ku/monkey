package repl

import "strings"

func hasPrefix(target string, rs ...string) bool {
	res := false
	for _, r := range rs {
		res = res || strings.HasPrefix(strings.TrimSpace(target), r)
	}
	return res
}
