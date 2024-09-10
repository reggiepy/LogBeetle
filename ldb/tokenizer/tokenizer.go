/**
 * TODO 待整理
 */
package tokenizer

import (
	"github.com/reggiepy/LogBeetle/com"
	"path/filepath"
	"runtime"
)

var sego *com.TokenizerSego

// 初始化装载字典
func InitSego() {
	_, filename, _, _ := runtime.Caller(0) // 当前go文件所在路径
	dictfile := filepath.Join(filepath.Dir(filename), "dict.txt")
	sego = com.NewTokenizerSego(dictfile)
}

// 按搜索引擎模式进行分词后返回分词数组
func CutForSearch(text string) []string {
	return sego.CutForSearch(text)
}

// 按搜索引擎模式进行分词后返回分词数组，可自定义添加或删除分词
func CutForSearchEx(text string, addWords []string, delWords []string) []string {
	return sego.CutForSearchEx(text, addWords, delWords)
}
