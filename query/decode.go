package query

import (
	"strconv"
)

func parse(str string) map[string]string {
	if str[len(str)-1] != '/' {
		str += "/"
	}
	return deScan(str)
}

func Decode(str string) map[string]string {
	return parse(str)
}

func deScan(str string) map[string]string {
	var r, t string
	runes := []rune(str)
	var kmp = make(map[string]string)
	for o := 0; o < len(runes); o++ {
		var a = runes[o]
		if '/' == a {
			if t != "" {
				kmp[t] = r
			}
			t, r = "", ""
		} else if '@' == a {
			o += 1
			a = runes[o]
			if 'A' == a {
				r += "@"
			} else if 'S' == a {
				r += "/"
			} else if '=' == a {
				t, r = r, ""
			}
		} else {
			s := strconv.QuoteRuneToGraphic(a)
			r += s[1 : len(s)-1]
		}
	}
	return kmp
}
