package itmo_trainer_api

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type config struct {
	dbname   string
	user     string
	password string
	host     string
}

func getConnection() (*sqlx.DB, error) {
	cfg := config{getDBNAME(), getDBUSER(), getDBPASSWORD(), getDBHOST()}
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s)/%s", cfg.user, cfg.password, cfg.host, cfg.dbname))
}
