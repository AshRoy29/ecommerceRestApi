package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *application) getOneProduct(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	product, err := app.models.DB.Get(id)

	err = app.writeJSON(w, http.StatusOK, product, "product")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := app.models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, products, "products")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}

func (app *application) deleteProduct(w http.ResponseWriter, r *http.Request) {

}

func (app *application) insertProduct(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {

}

func (app *application) searchProducts(w http.ResponseWriter, r *http.Request) {

}

func (app *application) test(w http.ResponseWriter, r *http.Request) {
	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}

	err := app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
