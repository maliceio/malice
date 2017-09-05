package watch

import "github.com/maliceio/engine/api/server/router"

// watchRouter is a router to talk with the plugin controller
type watchRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router {
	r := &watchRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the plugin controller
func (r *watchRouter) Routes() []router.Route {
	return r.routes
}

func (r *watchRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRoute("/watch/start", r.startWeb),
		router.NewGetRoute("/watch/stop", r.stopWeb),
		router.NewGetRoute("/watch/backup", r.backUpWeb),
	}
}
