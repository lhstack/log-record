package utils

import (
	md52 "crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	var md5 = md52.New()
	md5.Write([]byte(str))
	sum := md5.Sum(nil)
	return hex.EncodeToString(sum)
}
