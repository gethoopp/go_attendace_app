package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetDB() (*sql.DB, error) {
	// 1. PRIORITY: Heroku (JAWSDB)
	if jaws := os.Getenv("JAWSDB_URL"); jaws != "" {
		u, err := url.Parse(jaws)
		if err != nil {
			return nil, err
		}

		user := u.User.Username()
		pass, _ := u.User.Password()
		host := u.Host
		dbname := strings.TrimPrefix(u.Path, "/")

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}

		db.SetMaxIdleConns(20)
		db.SetMaxOpenConns(10)
		db.SetConnMaxLifetime(time.Minute * 5)

		return db, nil
	}

	// 2. fallback ke env lokal
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db, nil
}
