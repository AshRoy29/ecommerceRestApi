package main

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.signin)

	router.HandlerFunc(http.MethodGet, "/v1/product/:id", app.getOneProduct)
	router.HandlerFunc(http.MethodGet, "/v1/products", app.getAllProducts)

	router.POST("/test", app.wrap(secure.ThenFunc(app.test)))

	//router.HandlerFunc(http.MethodPost, "/test", app.test)

	router.HandlerFunc(http.MethodPost, "/image", app.insertProduct)

	return app.enableCORS(router)
}
