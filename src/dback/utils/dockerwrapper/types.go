package dockerwrapper

import (
	"regexp"
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

type InterruptPattern string

func NewInterruptPatterns(strArr []string) []InterruptPattern {
	res := []InterruptPattern{}

	for _, curPattern := range strArr {
		_, err := regexp.Compile(curPattern)
		check(err, `interrupt pattern must be valid regular expression, but the pattern `+
			curPattern+` is invalid`)

		res = append(res, InterruptPattern(curPattern))
	}

	return res
}
