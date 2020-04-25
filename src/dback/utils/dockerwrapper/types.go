package dockerwrapper

import (
	"strconv"
)

type ExcludePattern string

func NewExcludePatterns(strArr []string) []ExcludePattern {
	res := []ExcludePattern{}
	//toodo: check the each string is valid regExp

	return res
}

type EmulateFlag bool

func NewEmulateFlag(str string) EmulateFlag {
	res, err := strconv.ParseBool(str)
	check(err, `cannot convert "emulate" flag value "`+str+`" to boolean`)

	return EmulateFlag(res)
}
