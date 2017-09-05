package search

import "github.com/maliceio/engine/api/server/router"

// searchRouter is a router to talk with the plugin controller
type searchRouter struct {
	backend Backend
	routes  []router.Route
}

// NewRouter initializes a new plugin router
func NewRouter(b Backend) router.Router {
	r := &searchRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

// Routes returns the available routers to the plugin controller
func (r *searchRouter) Routes() []router.Route {
	return r.routes
}

func (r *searchRouter) initRoutes() {
	r.routes = []router.Route{
		router.NewGetRoute("/search", r.doSearch),
	}
}
