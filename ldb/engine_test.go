package ldb

import (
	"github.com/reggiepy/LogBeetle/com"
	"github.com/reggiepy/LogBeetle/ldb/search"
	"github.com/reggiepy/LogBeetle/ldb/storage/logdata"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	now := time.Now().Format("2006-01-02 15:04:05.000")
	traceId := com.RandomHashString()
	data := &logdata.LogDataModel{
		Text:       "测试用的日志，字段名为Text，" + "字段Date的格式为YYYY-MM-DD HH:MM:SS.SSS，必须格式一致才能正常使用时间范围检索条件。" + "随机3位字符串：" + com.RandomString(3) + "，第" + com.IntToString(1) + "条",
		Date:       now,
		System:     "demo1",
		ServerName: "default",
		ServerIp:   "127.0.0.1",
		ClientIp:   "127.0.0.1",
		TraceId:    traceId,
		LogLevel:   "INFO",
		User:       "tuser-" + com.RandomString(1),
	}
	engine := NewDefaultEngine()
	engine.AddLogDataModel(data)
}

func TestEngineTotalCount(t *testing.T) {
	engine := NewDefaultEngine()
	cond := &search.SearchCondition{
		StoreName:        "",       // 日志仓条件
		SearchKey:        "测试用的日志", // 输入的查询关键词
		CurrentStoreName: "",       // 滚动查询时定位用日志仓
		CurrentId:        0,        // 滚动查询时定位用ID
		Forward:          true,     // 是否向下滚动查询
		OldNearId:        0,        // 相邻检索旧ID
		NewNearId:        0,        // 相邻检索新ID
		NearStoreName:    "",       // 相邻检索时新ID对应的日志仓
		DatetimeFrom:     "",       // 日期范围（From）
		DatetimeTo:       "",       // 日期范围（To）
		OrgSystem:        "",  // 系统
		User:             "",  // 用户
		Loglevel:         "DEBUG",  // 日志级别
	}
	engine.Search(cond)
}
