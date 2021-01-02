package utils

import (
	"strings"
)

func HasOpenIDScope(scope string) bool {
	for _, s := range strings.Split(scope, " ") {
		if s == "openid" {
			return true
		}
	}
	return false
}
