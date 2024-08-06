package crud

import (
	dbConfig "botTelegram/dbconfig"
	"database/sql"
	"fmt"
)

var db *sql.DB
var err error

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func acessDatabase() {

	fmt.Printf("Acessing %s ...", dbConfig.DbName)
	db, err = sql.Open(dbConfig.PostgresDriver, dbConfig.DataSourceName)

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
