package public

import "github.com/reggiepy/LogBeetle/api"

type RouterPublic struct {
	RouterIndex      RouterIndex
	RouterMessage    RouterMessage
	RouterNsq        RouterNsq
	RouterSystem     RouterSystem
	RouterLog        RouterLog
	RouterStorageMnt RouterStorageMnt
}

var apiPublic = api.LbApi.ApiPublic
