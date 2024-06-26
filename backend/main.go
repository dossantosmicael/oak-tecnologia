package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/gorilla/mux"
)

type Product struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   string  `json:"available"`
}

type Data struct {
	Products []Product `json:"products"`
}

var dataFile = "backend/data.json"

func getProducts(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(dataFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Slice(data.Products, func(i, j int) bool {
		return data.Products[i].Price < data.Products[j].Price
	})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data.Products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile(dataFile, os.O_RDWR, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Products = append(data.Products, newProduct)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := file.Truncate(0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(file).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/products", getProducts).Methods("GET")
	router.HandleFunc("/api/products", createProduct).Methods("POST")
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend/"))))

	log.Println("Server is running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
