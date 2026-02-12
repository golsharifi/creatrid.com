package platform

import "strconv"

func parseIntOrZero(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
