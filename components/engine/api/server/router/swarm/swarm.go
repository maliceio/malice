package swarm

import "github.com/maliceio/engine/api/server/router"

// swarmRouter is a router to talk with the plugin controller
type swarmRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router {
	r := &swarmRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the plugin controller
func (r *swarmRouter) Routes() []router.Route {
	return r.routes
}

func (r *swarmRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRoute("/swarm/start", r.startWeb),
		router.NewGetRoute("/swarm/stop", r.stopWeb),
		router.NewGetRoute("/swarm/backup", r.backUpWeb),
	}
}
