package public

import "github.com/reggiepy/LogBeetle/service"

type ApiPublic struct {
	ApiMessage    ApiMessage
	ApiIndex      ApiIndex
	ApiNsq        ApiNsq
	ApiSystem     ApiSystem
	ApiLog        ApiLog
	ApiStorageMnt ApiStorageMnt
}

var serverPublic = service.LbService.ServicePublic
