package easycache

import "strings"

func defaultKeyGenerator(slug string, params ...string) string {
	return slug + ":" + strings.Join(params, ":")
}
