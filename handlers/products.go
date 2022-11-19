package handlers

import (
	"fmt"
	"log"
	"net/http"

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
