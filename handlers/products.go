package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kellemNegasi/product-api/data"
)

type Products struct {
	l *log.Logger
}
type KeyProduct struct{}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
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
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
	p.l.Printf("Prod: %#v", prod)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Printf("Unable to convert id %s", err)
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}
	p.l.Print(" Handle PUT Product")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		p.l.Println("product not found ")
		http.Error(w, "Product not found!", http.StatusNotFound)
		return
	}
	if err != nil {
		p.l.Println("internal error ")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJson(r.Body)
		if err != nil {
			p.l.Printf("failed to unmarshal %s", err)
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
