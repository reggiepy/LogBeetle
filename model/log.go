package model

type SearchRequest struct {
	StoreName        string `json:"store_name" form:"storeName"`                // 门店名称编号或ID
	SearchKey        string `json:"search_key" form:"search_key"`               // 消费者信息
	CurrentStoreName string `json:"current_store_name" form:"currentStoreName"` // 当前门店名称
	CurrentId        string `json:"current_id" form:"currentId"`                // 当前ID
	Forward          string `json:"forward" form:"forward"`                     // 前向（可能是方向或者其他标识）
	OldNearId        string `json:"old_near_id" form:"oldNearId"`               // 旧邻近ID
	NewNearId        string `json:"new_near_id" form:"newNearId"`               // 新邻近ID
	NearStoreName    string `json:"near_store_name" form:"nearStoreName"`       // 邻近门店名称
	DatetimeFrom     string `json:"datetime_from" form:"datetimeFrom"`          // 开始时间
	DatetimeTo       string `json:"datetime_to" form:"datetimeTo"`              // 结束时间
	System           string `json:"system" form:"system"`                       // 系统信息
	User             string `json:"user" form:"user"`                           // 用户信息
	LogLevel         string `json:"log_level" form:"loglevel"`                  // 日志级别
}
type AddLogRequest struct {
	Id         string `json:"id,omitempty"`         // 从1开始递增
	Text       string `json:"text,omitempty"`       // 【必须】日志内容，多行时仅为首行，直接显示用，是全文检索对象
	Date       string `json:"date,omitempty"`       // 日期（格式YYYY-MM-DD HH:MM:SS.SSS）
	System     string `json:"system,omitempty"`     // 系统名
	ServerName string `json:"servername,omitempty"` // 服务器名
	ServerIp   string `json:"serverip,omitempty"`   // 服务器IP
	ClientIp   string `json:"clientip,omitempty"`   // 客户端IP
	TraceId    string `json:"traceid,omitempty"`    // 跟踪ID
	LogLevel   string `json:"loglevel,omitempty"`   // 日志级别（debug、info、warn、error）
	User       string `json:"user,omitempty"`       // 用户
}
