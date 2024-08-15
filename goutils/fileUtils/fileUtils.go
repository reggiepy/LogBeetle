package fileUtils

import "os"

// FileExists 判断文件是否存在
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
