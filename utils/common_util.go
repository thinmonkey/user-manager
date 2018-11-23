package utils

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"time"
	"github.com/thinmonkey/user-manager/utils/log"
)

/**
检查error并抛出panic异常
*/
func CheckErr(err error) {
	if err != nil {
		log.Info(errors.WithStack(err))
	}
}

/**
将实例转换为string类型的字符串
*/
func Marshal(v interface{}) string {
	dataJson, err := json.Marshal(v)
	CheckErr(err)
	return fmt.Sprintf("%s\n", dataJson)
}

/**
将时间字符串转换成Long型时间
*/
func TimeConvertFromStringToLong(string string) int64 {
	tm2, err := time.Parse(string, string)
	CheckErr(err)
	return tm2.Unix()
}

//判断文件是否存在
func Exists(path string) (bool, error) {

	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	// 检测是否为路径不存在的错误
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

// 反转字符串
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
