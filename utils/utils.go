package utils

import "fmt"

func Pluralize(num int, thing string) string {
	if num != 1 && num != 0 {
		thing += "s"
	}
	return fmt.Sprintf("%d %s", num, thing)
}
