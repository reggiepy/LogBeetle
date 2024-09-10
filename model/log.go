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
