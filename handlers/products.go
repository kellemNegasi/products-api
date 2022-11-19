package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/kellemNegasi/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}
	if r.Method == http.MethodPost {
		p.AddProduct(w, r)
		return
	}
	if r.Method == http.MethodPut {
		path := r.URL.Path
		re := "/([0-9]+)"
		exp := regexp.MustCompile(re)
		matchingGroup := exp.FindAllStringSubmatch(path, -1)
		p.l.Println("matching group ", matchingGroup)
		if len(matchingGroup) != 1 {
			p.l.Printf("Invalid Uri in path %s", path)
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(matchingGroup[0]) != 2 {
			p.l.Printf("Invalid Uri in path %s", path)
			http.Error(w, "Invalid URI ", http.StatusBadRequest)
			return
		}
		idStr := matchingGroup[0][1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			p.l.Printf("Invalid Uri in path %s", path)
			http.Error(w, "Invalid URI ", http.StatusBadRequest)
		}
		p.UpdateProduct(id, w, r)
		return
	}

	fmt.Fprintln(w, "method not implemented ", http.StatusNotImplemented)
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to encode data to json.", http.StatusInternalServerError)
	}
	p.l.Print(" Handle GET Request \n")
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Print(" Handle POST Product")
	prod := &data.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		p.l.Printf("failed to unmarshal %s", err)
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Prod: %#v", prod)
}

func (p *Products) UpdateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Print(" Handle PUT Product")
	prod := &data.Product{}
	err := prod.FromJson(r.Body)
	if err != nil {
		p.l.Printf("failed to unmarshal %s", err)
		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		p.l.Println("product not found ")
		http.Error(w, "Product not found!", http.StatusNotFound)
	}
	if err != nil {
		p.l.Println("internal error ")
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
	return
}
