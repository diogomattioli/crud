package data

import (
	"regexp"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func VarsInt(ms map[string]string) (map[string]int, error) {

	mi := map[string]int{}

	for k, vs := range ms {

		vi, err := strconv.Atoi(vs)
		if err != nil {
			return nil, err
		}

		mi[k] = vi
	}

	return mi, nil
}
