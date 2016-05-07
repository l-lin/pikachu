package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Returns the routers for services and instances
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Availables routes
var routes = Routes{
	Route{
		"Services",
		"GET",
		"/service",
		GetServices,
	},
	Route{
		"Services",
		"GET",
		"/service/{serviceId}",
		GetService,
	},
	Route{
		"SaveService",
		"POST",
		"/service",
		SaveService,
	},
	Route{
		"UpdateService",
		"PUT",
		"/service/{serviceId}",
		UpdateService,
	},
	Route{
		"DeleteService",
		"DELETE",
		"/service/{serviceId}",
		DeleteService,
	},
	Route{
		"Instances",
		"GET",
		"/service/{serviceId}/instance",
		GetInstances,
	},
	Route{
		"SaveInstance",
		"POST",
		"/service/{serviceId}/instance",
		SaveInstance,
	},
	Route{
		"Instance",
		"GET",
		"/instance/{instanceId}",
		GetInstance,
	},
	Route{
		"UpdateInstance",
		"PUT",
		"/instance/{instanceId}",
		UpdateInstance,
	},
	Route{
		"DeleteInstance",
		"DELETE",
		"/instance/{instanceId}",
		DeleteInstance,
	},
}
