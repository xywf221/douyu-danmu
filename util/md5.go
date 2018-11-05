package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	return hex.EncodeToString(w.Sum(nil))
}
