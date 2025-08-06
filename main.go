package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var collection *mongo.Collection

func connectDB() {
	client, err :=
		mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel :=
		context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	collection = client.Database("testdb").Collection("books")
	fmt.Println("connected to MongoDB")
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

	cur, err :=
		collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var books []Book
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var book Book
		cur.Decode(&book)
		books = append(books, book)
	}
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
	_, err =
		collection.InsertOne(context.TODO(), newBook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	res, err :=
		collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.DeletedCount == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "book with ID %s deleted", id)
}

func main() {
	connectDB()
	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", router)
}
