package dbconfig

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Article struct {
	ID    int
	Name  string
	Price float32
}

type Order struct {
	OrderID   string
	ProductID string
}

type Users struct {
	Name  string
	Cpf   string
	Phone string
}

type Issues struct {
	Name string
	Cpf  string
}

var (
	PostgresDriver = "postgres"

	User     = getEnv("DB_USER", "bot")
	Password = getEnv("DB_PASSWORD", "bot")
	DbName   = getEnv("DB_NAME", "bot")
	Host     = getEnv("DB_HOST", "localhost")
	Port     = getEnv("DB_PORT", "5432")

	TableName   = "products"
	TableUser   = "users"
	TableIssues = "issues"
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var DataSourceName = fmt.Sprintf(
	"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	Host,
	Port,
	User,
	Password,
	DbName,
<<<<<<< HEAD
=======
	)
>>>>>>> 4549413cb2033ee4453dbd1ab4fdb5c9d5e5bb26
