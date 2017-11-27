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
		"ClientDefinition Status",
		"GET",
		"/clients/{clientName}",
		ShowClientStatus,
	},
	Route{
		"ClientDefinition Integration",
		"POST",
		"/clients/integrate",
		IntegrateNewClient,
	},
	Route{
		"ClientDefinition Backup Trigger",
		"POST",
		"/client/{clientName}/fs",
		TriggerClientBackup,
	},
}
