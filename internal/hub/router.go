package hub

import (
	"github.com/irlite/matrixd/internal/matrix/auth"
	"github.com/irlite/matrixd/internal/matrix/send"
	"github.com/irlite/matrixd/internal/wire"
)

type route struct {
	Type   string
	Action string
}

var routes = map[route]handler{
	{"ipc", "ping"}:    handle(ipcPing),
	{"auth", "login"}:  handle(auth.Login),
	{"auth", "logout"}: handle(auth.Logout),
	{"send", "text"}:   handle(send.Text),
}

func routeMessage(req request) wire.Outgoing {
	h, found := routes[route{req.Type, req.Action}]
	if !found {
		return fail(req, "unknown command "+req.Type+"."+req.Action)
	}
	return h(req)
}

func ipcPing(_ string, _ struct{}) (any, error) {
	return "pong", nil
}
