package web

import "github.com/maliceio/engine/api/server/router"

// webRouter is a router to talk with the plugin controller
type webRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router {
	r := &webRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the plugin controller
func (r *webRouter) Routes() []router.Route {
	return r.routes
}

func (r *webRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRoute("/web/start", r.startWeb),
		router.NewGetRoute("/web/stop", r.stopWeb),
		router.NewGetRoute("/web/backup", r.backUpWeb),
	}
}
