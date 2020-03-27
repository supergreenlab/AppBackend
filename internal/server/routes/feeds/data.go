package feeds

import (
	"github.com/spf13/viper"
	"upper.io/db.v3/postgresql"
)

var settings postgresql.ConnectionURL

func initDB() {
	settings = postgresql.ConnectionURL{
		Host:     "postgres",
		Database: "sglapp",
		User:     "postgres",
		Password: viper.GetString("PGPassword"),
	}
}
