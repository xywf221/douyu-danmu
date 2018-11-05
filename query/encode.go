package query

import (
	"douyuDm/util"
	"strconv"
)

func Encode(data map[string]interface{}) string {
	var str string
	for k, v := range data {
		str += addItem(k, util.InterfaceToString(v))
	}
	return str
}

func addItem(k, v string) string {
	var kp, vp string
	if k == "" {
		kp = ""
	} else {
		kp = scan(k) + "@="
	}
	vp = scan(v) + "/"
	return kp + vp;
}

func scan(str string) string {
	var t string
	for r := 0; r < len(str); r++ {
		var o = str[r]
		if o == '/' {
			t += "@S"
		} else if o == '@' {
			t += "@A"
		} else {
			t += strconv.QuoteRune(rune(o))[1:2]
		}
	}
	return t
}
