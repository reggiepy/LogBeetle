package model

type DeleteStoreRequest struct {
	StoreName string `json:"store_name" form:"store_name"`
}
