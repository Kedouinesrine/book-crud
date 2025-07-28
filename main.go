package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var Books = []Book{
	{ID: "1", Title: "Atomic Habits", Author: "James Clear"},
	{ID: "2", Title: "The Power of Now", Author: "Eckhart Tolle"},
	{ID: "3", Title: "The Alchemist Habits", Author: "Paulo coelho"},
	{ID: "4", Title: "Sapiens", Author: "Yuval Noah Harari"},
	{ID: "5", Title: "Think and Grow Rich", Author: "Napoleon hill"},
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var newBook Book

	err :=
		json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Books = append(Books, newBook)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	for index, book := range Books {
		if book.ID == id {
			Books = append(Books[:index], Books[index+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Book with ID %s delated", id)
			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", router)
}
