package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	*sql.DB
}

type DBConnectionParams struct {
	Host           string
	Port           int
	DBName         string
	User           string
	Password       string
	SSLMode        string
	ConnectTimeout int
}

func New(p *DBConnectionParams) (*PostgresDB, error) {
	source := fmt.Sprintf("host=%s port=%d dbname=%s connect_timeout=%d sslmode=%s", p.Host, p.Port, p.DBName, p.ConnectTimeout, p.SSLMode)
	log.Println("Connecting to Postgres:", source)
	source += fmt.Sprintf(" user=%s password=%s", p.User, p.Password)
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{DB: db}, nil
}

func main(){
	fmt.Println("Mock database connectivity with Postgres DB");
}
