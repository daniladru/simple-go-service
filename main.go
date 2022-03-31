package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/library")

	defer db.Close()

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
		return
	}

	// Execute the query
	results, err := db.Query("SELECT id, name FROM Books")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var book book
		err = results.Scan(&book.ID, &book.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		fmt.Println(book.Name)
	}

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/library")

	defer db.Close()

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
		return
	}

	res, err := db.Query("delete  FROM Books where id = '1'")

	defer res.Close()

	if err != nil {
		panic(err.Error())
	}

}

func createBook(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/library")

	defer db.Close()

	if err != nil {
		log.Print(err.Error())
		return
	}

	sql := "INSERT INTO books(name) VALUES ('New book name')"
	res, err := db.Exec(sql)

	if err != nil {
		panic(err.Error())
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The last inserted row id: %d\n", lastId)

}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/api/list", getBooks).Methods("GET")
	r.HandleFunc("/api/delete/{id}", deleteBook).Methods("DELETE")
	r.HandleFunc("/api/add", createBook).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", r))
}
