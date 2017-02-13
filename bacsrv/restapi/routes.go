package restapi

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Status",
		"GET",
		"/status",
		StatusIndex,
	},
	Route{
		"Clients List",
		"GET",
		"/clients",
		ShowClients,
	},
	Route{
		"Client Status",
		"GET",
		"/clients/{clientName}",
		ShowClientStatus,
	},
	Route{
		"Client Integration",
		"POST",
		"/clients/integrate",
		IntegrateNewClient,
	},
}
