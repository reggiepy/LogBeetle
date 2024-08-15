package model

import (
	"github.com/gin-gonic/gin"
	"github.com/reggiepy/LogBeetle/global"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Result(code int, data interface{}, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		message,
	})
}

func ResponseSuccessMessage(message string, c *gin.Context) {
	Result(200, nil, message, c)
}

func ResponseSuccessMessageWithLog(message string, c *gin.Context) {
	global.LbLogger.Info(message)
	ResponseSuccessMessage(message, c)
}

func ResponseSuccessData(data interface{}, c *gin.Context) {
	Result(200, data, "Success", c)
}

func ResponseSuccessMessageData(data interface{}, message string, c *gin.Context) {
	Result(200, data, message, c)
}

func ResponseSuccessMessageDataWithLog(data interface{}, message string, c *gin.Context) {
	global.LbLogger.Info(message)
	ResponseSuccessMessageData(data, message, c)
}

func ResponseErrorMessage(message string, c *gin.Context) {
	Result(500, nil, message, c)
}

func ResponseErrorMessageWithLog(message string, c *gin.Context) {
	global.LbLogger.Error(message)
	ResponseErrorMessage(message, c)
}

func ResponseErrorData(data interface{}, c *gin.Context) {
	Result(500, data, "Failed", c)
}

func ResponseErrorMessageData(data interface{}, message string, c *gin.Context) {
	Result(500, data, message, c)
}

type ResponsePage struct {
	Records  interface{} `json:"records"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type ResponsePageWithParentId struct {
	Records    interface{} `json:"records"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	ParentCode string      `json:"parentCode"`
}
