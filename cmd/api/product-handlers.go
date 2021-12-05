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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ProductPayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Stock       string `json:"stock"`

	CategoryID string `json:"category"`
}

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

var dir string

var product models.Product

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

func (app *application) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.models.DB.GetAllCategory()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, categories, "categories")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllProductsByCategory(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	categoryID, err := strconv.Atoi(params.ByName("category_id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	products, err := app.models.DB.All(categoryID)
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
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	productImage, err := app.models.DB.DeleteProduct(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println("IMAGE", productImage)

	err = os.Remove(productImage)

	if err != nil {
		fmt.Println(err)
		return
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

//this was insertProduct
func (app *application) uploadImage(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(tempFile.Name())

	//files, _ := os.ReadDir(dir)
	path, _ := filepath.Abs(dir)
	//for _, file := range files {
	fmt.Println("path:", path)
	imageDir := filepath.Join(path, tempFile.Name())

	var i string

	i = imageDir
	imageDestination(i)

}

func imageDestination(i string) {

	product.Image = i

}

func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {

}

func (app *application) searchProducts(w http.ResponseWriter, r *http.Request) {

}

//this was test
func (app *application) editProducts(w http.ResponseWriter, r *http.Request) {

	var payload ProductPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	log.Println(payload.Title)
	log.Println(payload.Size)

	arraySize := strings.Split(payload.Size, ",")

	var category models.Category

	if payload.ID != "0" {
		id, _ := strconv.Atoi(payload.ID)
		p, _ := app.models.DB.Get(id)
		product = *p
		product.UpdatedAt = time.Now()
	}

	product.ID, _ = strconv.Atoi(payload.ID)
	product.Title = payload.Title
	product.Price, _ = strconv.Atoi(payload.Price)
	product.Size = arraySize
	product.Description = payload.Description
	product.Stock, _ = strconv.Atoi(payload.Stock)
	log.Println("ps:", product.Image)
	category.ID, _ = strconv.Atoi(payload.CategoryID)

	log.Println("Product price:", product.Price)

	if product.ID == 0 {
		newProductID, err := app.models.DB.InsertProduct(product)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		log.Println("category ID:", category.ID)

		productCategory := models.ProductCategory{
			ProductID:  newProductID,
			CategoryID: category.ID,
		}

		err = app.models.DB.InsertCategory(productCategory)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

	} else {
		err = app.models.DB.UpdateProduct(product)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
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
