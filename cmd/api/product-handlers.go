package main

import (
	"ecom-api/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ProductPayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	Description string `json:"description"`
	//Image string `json:"image"`

}

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
	r.ParseMultipartForm(10 * 1024 * 1024)
	file, handler, err := r.FormFile("image")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fmt.Println("file info")
	fmt.Println("file name:", handler.Filename)
	fmt.Println("file size:", handler.Size)
	fmt.Println("file type:", handler.Header.Get("Content-Type"))

	tempFile, err := ioutil.TempFile("img", "img-*.jpg")
	if err != nil {
		log.Println(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	tempFile.Write(fileBytes)

	fmt.Println("SUCCESS")
}

func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {

}

func (app *application) searchProducts(w http.ResponseWriter, r *http.Request) {

}

func (app *application) test(w http.ResponseWriter, r *http.Request) {

	var payload ProductPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	log.Println(payload.Title)

	var product models.Product

	product.ID, _ = strconv.Atoi(payload.ID)
	product.Title = payload.Title
	product.Price, _ = strconv.Atoi(payload.Price)
	product.Description = payload.Description

	//app.insertProduct(w, r)

	log.Println(product.Price)

	type jsonResp struct {
		OK bool `json:"ok"`
	}

	ok := jsonResp{
		OK: true,
	}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

}
