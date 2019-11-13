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
			"ListResources",
			"GET",
			"/apis/resource",
			handler.ListResources,
		},
		{
			"ListResourcesWithKind",
			"GET",
			"/apis/resource/{kind}",
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
			"/apis/resource/{kind}/{resource}",
			handler.UpdateResource,
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})

	return c.Handler(router)
}
