package args

import "strings"

func KeyValue(arg string) (string, string) {
	keValue := strings.SplitN(arg, "=", 2)
	if len(keValue) < 2 {
		return keValue[0], ""
	}
	return keValue[0], keValue[1]
}
