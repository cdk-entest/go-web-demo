package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// aws rds
const HOST = ""
const USER = "postgresql"
const PASS = "demo2024"
const DBNAME = "demo"

// local db
// const HOST = "localhost"
// const USER = "postgres"
// const DBNAME = "dvdrental"
// const PASS = ""

type Book struct {
	ID          uint
	Title       string
	Author      string
	Amazon      string
	Image       string
	Description string
}

func main() {

	// db init
	dns := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v", HOST, "5432", USER, PASS, DBNAME)

	db, error := gorm.Open(postgres.Open(dns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase:   false,
			SingularTable: true,
		},
	})

	if error != nil {
		fmt.Println(error)
	}

	mux := http.NewServeMux()

  // static files 
  mux.Handle("/demo/", http.StripPrefix("/demo/", http.FileServer(http.Dir("./static"))))

	// home page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

  // postgresql page 
	mux.HandleFunc("/postgresql", func(w http.ResponseWriter, r *http.Request) {

		// query a list of book []Book
		books := getBooks(db)

		// load template
		tmpl, error := template.ParseFiles("./static/book-template.html")

		if error != nil {
			fmt.Println(error)
		}

		// pass data to template and write to writer
		tmpl.Execute(w, books)
	})

	// create web server
	server := &http.Server{
		Addr:           ":80",
		Handler:        mux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// enable logging
	log.Fatal(server.ListenAndServe())

}

func getBooks(db *gorm.DB) []Book {
	var books []Book

	db.Limit(10).Find(&books)

	for _, book := range books {
		fmt.Println(book.Title)
	}

	return books
}
