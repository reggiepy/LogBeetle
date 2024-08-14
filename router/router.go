package router

import (
	"github.com/reggiepy/LogBeetle/router/private"
	"github.com/reggiepy/LogBeetle/router/public"
)

var LbRouter = new(struct {
	RouterPublic  public.RouterPublic
	RouterPrivate private.RouterPrivate
})
