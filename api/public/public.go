package public

import "github.com/reggiepy/LogBeetle/service"

type ApiPublic struct {
	ApiMessage ApiMessage
	ApiIndex   ApiIndex
}

var serverPublic = service.LbService.ServicePublic
