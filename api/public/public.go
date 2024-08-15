package public

import "github.com/reggiepy/LogBeetle/service"

type ApiPublic struct {
	ApiMessage ApiMessage
	ApiIndex   ApiIndex
	ApiNsq     ApiNsq
}

var serverPublic = service.LbService.ServicePublic
