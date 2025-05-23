package boot

import (
	"github.com/reggiepy/LogBeetle/ldb"
	"github.com/reggiepy/LogBeetle/ldb/tokenizer"
)

func Ldb() {
	tokenizer.InitSego()
	// 默认引擎空转一下，触发未建索引继续建
	go ldb.NewDefaultEngine().AddTextLog("", "", "")
}
