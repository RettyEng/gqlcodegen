package generator

import (
	"strings"
)

func capitalizeFirst(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}
