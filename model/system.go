package model

type SystemInfoResponse struct {
	GoroutineNumber int          `json:"goroutine_number" form:"goroutineNumber"`
	ConsumerInfo    ConsumerInfo `json:"consumer_info" form:"consumerInfo"`
}

type ConsumerInfo struct {
	ConsumerCount int `json:"consumer_count" form:"consumerCount"`
}
