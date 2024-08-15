package public

import "github.com/reggiepy/LogBeetle/api"

type RouterPublic struct {
	RouterIndex   RouterIndex
	RouterMessage RouterMessage
	RouterNsq     RouterNsq
}

var apiPublic = api.LbApi.ApiPublic
