package str

import (
	"strings"
)

func StringContains(slices []string, comparizon string) bool {
	for _, a := range slices {
		if a == comparizon {
			return true
		}
	}

	return false
}

func StringContainsPrefix(prefixSlices []string, s string) bool {

	for _, a := range prefixSlices {
		if strings.HasPrefix(s, a) {
			return true
		}
	}

	return false
}