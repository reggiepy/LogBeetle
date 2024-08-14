package api

import (
	"github.com/reggiepy/LogBeetle/api/private"
	"github.com/reggiepy/LogBeetle/api/public"
)

var LbApi = new(struct {
	ApiPublic  public.ApiPublic
	ApiPrivate private.ApiPrivate
})
