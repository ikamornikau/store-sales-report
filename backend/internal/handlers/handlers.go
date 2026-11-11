package handlers

import "fmt"

func MakeRoute(prefix, method, pattern string) string {
	return fmt.Sprintf("%s %s%s", method, prefix, pattern)
}
