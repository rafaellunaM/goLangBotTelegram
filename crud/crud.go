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

func acessDatabase() error {

	var err error
    db, err = sql.Open(dbConfig.PostgresDriver, dbConfig.DataSourceName)
    if err != nil {
        return fmt.Errorf("falha ao abrir a conexão com o banco de dados: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return fmt.Errorf("falha ao conectar com o banco de dados: %w", err)
    }
    return nil
}

/* func SqlSelect() dbConfig.Article{

	sqlStatement, err := db.Query("SELECT * FROM " + dbConfig.TableName)
	checkErr(err)

	var products dbConfig.Article

	for sqlStatement.Next() {

		err = sqlStatement.Scan(&products.ID, &products.Name, &products.Price)

		fmt.Printf("%d\t%s\t%s \n", products.ID, products.Name, products.Price)

		return products
	}
	return products
} */

func GetProducts() ([]dbConfig.Article, error) {
    
    if err := acessDatabase(); err != nil {
        return nil, fmt.Errorf("erro ao acessar o banco de dados: %w", err)
    }

    rows, err := db.Query("SELECT * FROM " + dbConfig.TableName)
    if err != nil {
        return nil, fmt.Errorf("erro ao executar a consulta: %w", err)
    }
    defer rows.Close() 

    var products []dbConfig.Article
    for rows.Next() {
        var product dbConfig.Article
        if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
            return nil, fmt.Errorf("erro ao ler os dados: %w", err)
        }
        products = append(products, product)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erro durante a iteração dos resultados: %w", err)
    }

    return products, nil
}

