package scan

import "github.com/maliceio/engine/api/server/router"

// scanRouter is a router to talk with the plugin controller
type scanRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router {
	r := &scanRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the plugin controller
func (r *scanRouter) Routes() []router.Route {
	return r.routes
}

func (r *scanRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRoute("/scans", r.doScan),
	}
}
