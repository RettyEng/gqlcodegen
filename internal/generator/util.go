package generator

import (
	"strings"
)

const commentOnTop = "// DO NOT EDIT. this sourcecode is generated by gqlcodegen.\n"

func capitalizeFirst(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}
