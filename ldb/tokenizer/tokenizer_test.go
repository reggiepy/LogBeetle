package tokenizer

import (
	"fmt"
	"testing"
)

func Test_CutWords(t *testing.T) {
	ws := CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造，Java和Go都学得不错，Java和Go都不错的")
	fmt.Println(ws)
}

func Test_CutWords2(t *testing.T) {
	ws := CutForSearch("2024-08-28 13:59:04,797 -- INFO  -- 【76】 【903b234c90f4461294d49a7090858158】 【76--903b234c90f4461294d49a7090858158】退出计算 (inner_wrapper in new_runService.py):45]")
	fmt.Println(ws)
}
