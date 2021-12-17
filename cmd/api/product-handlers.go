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

type CartPayload struct {
	Product []Product `json:"product"`
	UserID  int       `json:"user"`
	Total   int       `json:"total"`
}

type Product struct {
	ID       string `json:"id"`
	Size     string `json:"size"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type OrderStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var dir string
var imageDir string

var cartUserID int
var cartOrderID int

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
	imageDir = filepath.Join(path, tempFile.Name())

}

func imageDestination(i string) {

	product.Image = i
	log.Println("product ID:", product.ID)
}

func (app *application) userCart(w http.ResponseWriter, r *http.Request) {

	var payload CartPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var cart models.CartProducts

	cart.ProductID = make([]string, len(payload.Product))
	cart.Size = make([]string, len(payload.Product))
	cart.Price = make([]string, len(payload.Product))
	cart.Quantity = make([]string, len(payload.Product))

	for i := 0; i < len(payload.Product); i++ {
		productID := strings.Split(payload.Product[i].ID, ",")
		cart.ProductID[i] = productID[0]
		cart.Size[i] = productID[1]
		cart.Price[i] = strconv.Itoa(payload.Product[i].Price)
		cart.Quantity[i] = strconv.Itoa(payload.Product[i].Quantity)
	}

	cart.UserID = payload.UserID
	cart.Total = payload.Total

	cartUserID, cartOrderID, err = app.models.DB.CartOrders(cart)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Println(cartUserID)
	log.Println(cartOrderID)

}

func (app *application) userBill(w http.ResponseWriter, r *http.Request) {

	var bill models.BillingInfo

	err := json.NewDecoder(r.Body).Decode(&bill)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	bill.UserID = cartUserID
	bill.OrderID = cartOrderID

	err = app.models.DB.BillingInfo(bill)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := app.models.DB.AllOrders()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, orders, "orders")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) orderStatus(w http.ResponseWriter, r *http.Request) {

	var payload OrderStatus

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	orderID, _ := strconv.Atoi(payload.ID)

	status := models.CartProducts{
		ID:     orderID,
		Status: payload.Status,
	}

	err = app.models.DB.UpdateStatus(status)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

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
		imageDestination(imageDir)
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
		err = os.Remove(product.Image)
		if err != nil {
			fmt.Println(err)
			return
		}

		imageDestination(imageDir)
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
