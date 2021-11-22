package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodGet, "/v1/product/:id", app.getOneProduct)
	router.HandlerFunc(http.MethodGet, "/v1/products", app.getAllProducts)

	router.HandlerFunc(http.MethodPost, "/test", app.test)

	return app.enableCORS(router)
}
