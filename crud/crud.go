package crud

import (
	dbConfig "botTelegram/dbconfig"
	"database/sql"
	"fmt"
	"log"
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

func GetUsers() ([]dbConfig.Users, error) {
	if err := acessDatabase(); err != nil {
		return nil, fmt.Errorf("erro ao acessar o banco de dados: %w", err)
	}

	rows, err := db.Query("SELECT * FROM " + dbConfig.TableUser)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar a consulta: %w", err)
	}
	defer rows.Close()

	var users []dbConfig.Users
	for rows.Next() {
		var user dbConfig.Users
		if err := rows.Scan(&user.Name, &user.Cpf, &user.Phone); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração dos resultados: %w", err)
	}

	return users, nil
}

func GetIssues() ([]dbConfig.Issues, error) {
	if err := acessDatabase(); err != nil {
		return nil, fmt.Errorf("erro ao acessar o banco de dados: %w", err)
	}

	rows, err := db.Query("SELECT * FROM " + dbConfig.TableIssues)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar a consulta: %w", err)
	}
	defer rows.Close()

	var issues []dbConfig.Issues
	for rows.Next() {
		var issue dbConfig.Issues
		if err := rows.Scan(&issue.Cpf, &issue.Name); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados: %w", err)
		}
		issues = append(issues, issue)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração dos resultados: %w", err)
	}

	return issues, nil
}

func SetUsers(cpf string, name string, phone string, issues string) {

	acessDatabase()
	insertUser := `INSERT INTO users (cpf, name, phone) VALUES ($1, $2, $3)`
	_, err := db.Exec(insertUser, cpf, name, phone)
	if err != nil {
		log.Printf("Erro ao salvar a mensagem no banco de dados: %v", err)
	}

	insertIssues := `INSERT INTO issues (issues, cpf) VALUES ($1, $2);`
	_, err = db.Exec(insertIssues, issues, cpf)
	if err != nil {
		log.Printf("Erro ao salvar a mensagem no banco de dados: %v", err)
	}

}

// func GetOrders(pedido string) ([]dbConfig.Order, error) {
// 	acessDatabase()

// 	// Consulta apenas a tabela orders
// 	getOrder := `SELECT orders.pedido AS OrderID, orders.produto AS ProductID FROM orders WHERE orders.pedido = $1`

// 	rows, err := db.Query(getOrder, pedido)
// 	if err != nil {
// 		log.Printf("Erro ao consultar pedido: %v", err)
// 		return nil, err
// 	}

// 	var orders []dbConfig.Order
// 	for rows.Next() {
// 		var order dbConfig.Order
// 		if err := rows.Scan(&order.OrderID, &order.ProductID); err != nil {
// 			return nil, fmt.Errorf("erro ao ler os dados: %w", err)
// 		}
// 		orders = append(orders, order)
// 	}

// 	return orders, nil
// }

func GetOrders(pedido string) (string, error) {
	acessDatabase()

	// Consulta apenas a tabela orders
	getOrder := `SELECT products.name AS ProductName, products.price  FROM orders INNER JOIN products ON orders.produto = products.ID WHERE orders.pedido = $1;`

	rows, err := db.Query(getOrder, pedido)
	if err != nil {
		log.Printf("Erro ao consultar pedido: %v", err)
		return "", err
	}

	var orderStr string
	for rows.Next() {
		var order dbConfig.Order
		if err := rows.Scan(&order.OrderID, &order.ProductID); err != nil {
			return "", fmt.Errorf("erro ao ler os dados: %w", err)
		}
		// Adiciona a linha formatada à string
		orderStr += fmt.Sprintf("%s | R$ %s\n", order.OrderID, order.ProductID)
	}

	return orderStr, nil
}

// func GetOrders(productid string) ([]dbConfig.Order, error) {
// 	acessDatabase()

// 	getOrder := "SELECT products.ID AS ProductID, products.Name AS ProductName, products.Price AS ProductPrice FROM orders INNER JOIN products ON orders.produto = $1"
// 	rows, err := db.Query(getOrder, productid)
// 	if err != nil {
// 		log.Printf("Erro ao consultar pedido %v", err)
// 	}
// 	var orders []dbConfig.Order
// 	for rows.Next() {
// 		var order dbConfig.Order
// 		if err := rows.Scan(&order.OrderID, &order.ProductID); err != nil {
// 			return nil, fmt.Errorf("erro ao ler os dados: %w", err)
// 		}
// 		orders = append(orders, order)
// 	}

// 	return orders, nil
// }
