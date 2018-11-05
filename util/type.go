package util

import "strconv"

func InterfaceToString(data interface{}) string {
	switch data.(type) {
	case string:
		return data.(string)
	case int:
		var i = data.(int)
		return strconv.Itoa(i)
	case int64:
		var i = data.(int64)
		return strconv.FormatInt(i, 10)
	}
	return ""
}
