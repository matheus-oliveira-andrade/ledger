package repositories

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func NewDBConnection() *sql.DB {
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	name := viper.GetString("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	result, err := sql.Open("postgres", connString)
	if err != nil {
		slog.Error("Error creating db connection", " error", err.Error())
		panic(err)
	}

	return result
}
