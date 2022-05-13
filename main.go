package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"

	"gorm.io/gorm"
	"log"
	"net/http"
)

type Person struct {
	gorm.Model
	Name  string
	Email string `gorm:"type:varchar(100);unique_index"`
	Books []Book
}
type Book struct {
	gorm.Model
	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}

var db *gorm.DB
var err error

func main() {
	// get all .env variables
	user := "postgres"
	password := "example"
	dbName := "postgres"
	host := "localhost"
	port := "5432"
	// cdatabase connection sting
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, port)

	//db connection
	db, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Person{}, &Book{})
	// Api routes
	r := mux.NewRouter()
	// People routes
	r.HandleFunc("/api/people", getPeople).Methods("GET")
	r.HandleFunc("/api/people/{id}", getPerson).Methods("GET")
	r.HandleFunc("/api/people", createPerson).Methods("POST")
	r.HandleFunc("/api/people/{id}", updatePerson).Methods("PUT")
	r.HandleFunc("/api/people/{id}", deletePerson).Methods("DELETE")
	// Books routes
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))

}

// person controllers
func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)
_:
	json.NewEncoder(w).Encode(people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	var books []Book
	db.First(&person, params["id"])
	db.Model(&person).Association("Books").Find(&books)
	json.NewEncoder(w).Encode(&person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	db.Create(&person)
_:
	json.NewEncoder(w).Encode(person)
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	db.Model(&person).Where("id = ?", params["id"]).Updates(&person)
	_ = json.NewDecoder(r.Body).Decode(&person)

}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	db.Where("id = ?", params["id"]).Delete(&person)
}

// book controllers
func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	db.Find(&books)
_:
	json.NewEncoder(w).Encode(books)
}
func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.First(&book, params["id"])
_:
	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	db.Create(&book)
	_ = json.NewDecoder(r.Body).Decode(&book)

}
func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.Where("id = ?", params["id"]).Updates(book)
	_ = json.NewDecoder(r.Body).Decode(&book)

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.Where("id = ?", params["id"]).Delete(&book)
}
