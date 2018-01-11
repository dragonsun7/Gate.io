package main

import (
	"regexp"
	"strconv"
)

/* 字符串转换为十进制整数（支持逗号分隔） */
func StrToInt(s string) int64  {
	reg := regexp.MustCompile(",")
	new := reg.ReplaceAllString(s, "")
	ret, err := strconv.ParseInt(new, 10, 64)
	if err != nil {
		return 0
	} else {
		return ret
	}
}

/* 字符串转换为浮点数（支持逗号分隔） */
func StrToFloat(s string) float64 {
	reg := regexp.MustCompile(",")
	new := reg.ReplaceAllString(s, "")
	ret, err := strconv.ParseFloat(new, 64)
	if err != nil {
		return 0.0
	} else {
		return ret
	}
}
