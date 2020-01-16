package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewProduct(n string, p float64) Product {
	return Product{
		Name:  n,
		Price: p,
	}
}

func (p Product) String() string {
	return "This is a product"
}

type Products struct {
	inventory []Product
	lock      sync.Mutex
}

func (p *Products) add(product Product) {
	p.lock.Lock()
	p.inventory = append(p.inventory, product)
	p.lock.Unlock()
}

func main() {
	var soap Product = NewProduct("soap", 4.20)

	var shampoo Product = Product{
		Name:  "shampoo",
		Price: 420,
	}

	var products Products = Products{}

	products.add(soap)
	products.add(shampoo)

	fmt.Println(soap)

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		output, err := json.Marshal(products.inventory)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			fmt.Fprintln(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(output))
	})

	fmt.Println("Listening on ports 8080 (http) and 8081 (https)...")

	errorChan := make(chan error, 5)
	go func() {
		errorChan <- http.ListenAndServe(":8080", nil)
	}()
	go func() {
		errorChan <- http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	for {
		select {
		case err := <-errorChan:
			if err != nil {
				log.Fatalln(err)
			}

		case sig := <-signalChan:
			fmt.Println("\nShutting down due to", sig)
			os.Exit(0)
		}
	}
}
