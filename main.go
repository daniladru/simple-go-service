package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/ini.v1"
)

type book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func getconnect() *sql.DB {
	cfg, err := ini.Load("setup.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return nil
	}

	//db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/library")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s",
		cfg.Section("mysql").Key("user"),
		cfg.Section("mysql").Key("password"),
		cfg.Section("mysql").Key("host"),
		cfg.Section("mysql").Key("database")))

	if err != nil {
		log.Print(err.Error())
		return nil
	}

	return db

}

func getBooks(w http.ResponseWriter, r *http.Request) {
	db := getconnect()

	if db == nil {
		return
	}

	defer db.Close()

	results, err := db.Query("SELECT id, name FROM Books")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var book book
		err = results.Scan(&book.ID, &book.Name)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(book.Name)
	}

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	db := getconnect()

	if db == nil {
		return
	}

	defer db.Close()

	res, err := db.Query("delete  FROM Books where id = '1'")

	defer res.Close()

	if err != nil {
		panic(err.Error())
	}

}

func createBook(w http.ResponseWriter, r *http.Request) {
	db := getconnect()

	if db == nil {
		return
	}

	defer db.Close()

	sql := "INSERT INTO Books(name) VALUES ('New book name')"
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
