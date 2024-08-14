package service

import (
	"github.com/reggiepy/LogBeetle/service/private"
	"github.com/reggiepy/LogBeetle/service/public"
)

var LbService = new(struct {
	ServicePublic  public.ServicePublic
	ServicePrivate private.ServicePrivate
})
