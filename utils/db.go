package utils

import (
	"fmt"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// GetMySQLDb returns MySQL connection
func GetMySQLDb() (*sql.DB, error) {
	user := viper.GetString("mysql-user")
	pass := viper.GetString("mysql-pass")
	port := viper.GetString("mysql-port")
	host := viper.GetString("mysql-host")
	protocol := viper.GetString("mysql-protocol")

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/", user, pass, protocol, host, port)
	return sql.Open("mysql", dsn)
}

// GetPostgresDb returns PostgreSQL connection
func GetPostgresDb() (*sql.DB, error) {
	user := viper.GetString("postgres-user")
	pass := viper.GetString("postgres-pass")
	port := viper.GetString("postgres-port")
	host := viper.GetString("postgres-host")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/?sslmode=disable", user, pass, host, port)
	return sql.Open("postgres", dsn)
}
