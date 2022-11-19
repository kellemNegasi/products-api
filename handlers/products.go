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
	}

	fmt.Fprintln(w, "method not implemented ", http.StatusNotImplemented)
}

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJson(w)
	if err != nil {
		http.Error(w, "Unable to encode data to json.", http.StatusInternalServerError)
	}
	p.l.Print(" GET \n")
}
