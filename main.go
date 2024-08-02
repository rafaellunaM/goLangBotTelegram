package main

import (
	"database/sql"
	"fmt"
	"net/http"

	dbConfig "botTelegram/dbconfig"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app example")
	})

	port := ":8080"
	fmt.Printf("porta%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %s\n", err)
	}

	fmt.Printf("Acessing %s ...", dbConfig.DbName)
	db, err = sql.Open(dbConfig.PostgresDriver, dbConfig.DataSourceName)

	if err != nil {
		panic(err.Error())

	} else {
		fmt.Println("Connected")
	}

	defer db.Close()

}

func sqlSelect() {

	sqlStatement, err := db.Query("SELECT * FROM " + dbConfig.TableName)
	checkErr(err)

	for sqlStatement.Next() {

		var products dbConfig.Article

		err = sqlStatement.Scan(&products.ID, &products.Name, &products.Price)

		fmt.Printf("%d\t%s\t%s \n", products.ID, products.Name, products.Price)

	}
}
