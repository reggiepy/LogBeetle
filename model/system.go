package model

type SystemInfoResponse struct {
	GoroutineNumber int          `json:"goroutine_number" form:"goroutineNumber"`
	ConsumerInfo    ConsumerInfo `json:"consumer_info" form:"consumerInfo"`
	StartTime       string       `json:"start_time" form:"start_time"`
}

type ConsumerInfo struct {
	ConsumerCount int `json:"consumer_count" form:"consumerCount"`
}

type SystemInfoRequest struct {
}
