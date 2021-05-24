package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

// PrepareDBConnection は dbstring を受け取って、 mysql の *sql.DB を返します。 接続に失敗した場合は panic します
func PrepareDBConnection(user, pass, host, port, name string, maxOpen, maxIdle int, conLifetime string) *sqlx.DB {
	d, err := time.ParseDuration(conLifetime)
	if err != nil {
		panic(fmt.Sprintf("mysql connection lifetime setting is wrong (%v).", conLifetime))
	}

	c := sqlx.MustConnect("mysql", PrepareDBString(user, pass, host, port, name))
	c.SetMaxOpenConns(maxOpen)
	c.SetMaxIdleConns(maxIdle)
	c.SetConnMaxLifetime(d)
	return c
}

// PrepareDBString は dbuser, dbpass, dbhost, dbport, dbname を受け取って、 mysql driver の connection string を返します
func PrepareDBString(user, pass, host, port, name string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		pass,
		host,
		port,
		name,
	)
}
