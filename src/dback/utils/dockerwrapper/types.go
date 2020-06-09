package dockerwrapper

import (
	"regexp"
	"strconv"
)

type ExcludePattern string

func NewExcludePatterns(strArr []string) []ExcludePattern {
	res := []ExcludePattern{}

	for _, curPattern := range strArr {
		_, err := regexp.Compile(curPattern)
		check(err, `exlude pattern must be valid regular expression, but the pattern `+
			curPattern+` is invalid`)

		res = append(res, ExcludePattern(curPattern))
	}

	return res
}

type EmulateFlag bool

func NewEmulateFlag(str string) EmulateFlag {
	res, err := strconv.ParseBool(str)
	check(err, `cannot convert "emulate" flag value "`+str+`" to boolean`)

	return EmulateFlag(res)
}
