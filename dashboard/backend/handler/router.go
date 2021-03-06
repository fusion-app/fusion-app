package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(handler *APIHandler) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	var publicRoutes = Routes{
		{
			"ListResourcesWithKind",
			"GET",
			"/apis/resource",
			handler.ListResourcesWithKind,
		},
		{
			"CreateResource",
			"POST",
			"/apis/resource",
			handler.CreateResource,
		},
		{
			"UpdateResource",
			"PUT",
			"/apis/resource",
			handler.UpdateResource,
		},
		{
			"BindResource",
			"POST",
			"/apis/resource/bind",
			handler.BindResource,
		},
		{
			"UnBindResource",
			"DELETE",
			"/apis/resource/bind",
			handler.UnBindResource,
		},
		{
			"QueryApp",
			"GET",
			"/apis/app",
			handler.ListApp,
		},
		{
			"UpdateApp",
			"PUT",
			"/apis/app",
			handler.UpdateApp,
		},
		{
			"CreateApp",
			"POST",
			"/apis/app",
			handler.CreateApp,
		},
		{
			"QueryAppInstance",
			"GET",
			"/apis/app_instance",
			handler.QueryAppInstance,
		},
		{
			"CreateAppInstance",
			"POST",
			"/apis/app_instance",
			handler.CreateAppInstance,
		},
		{
			"DeleteAppInstance",
			"DELETE",
			"/apis/app_instance",
			handler.DeleteAppInstance,
		},
		{
			"ListAppInstance",
			"GET",
			"/apis/app_instances",
			handler.ListAppInstance,
		},
		{
			"UpdateAppInstance",
			"PUT",
			"/apis/app_instance",
			handler.UpdateAppInstance,
		},
	}

	// The public route is always accessible
	for _, route := range publicRoutes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	// Handle static files.
	router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(handler.frontDir))))

	c := cors.Default()

	return c.Handler(router)
}
