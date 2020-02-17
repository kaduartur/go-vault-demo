package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	driveName         = "postgres"
	dataSourcePattern = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
)

var (
	host     = os.Getenv("DATABASE_HOST")
	port     = os.Getenv("DATABASE_PORT")
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PASSWORD")
	dbname   = os.Getenv("DATABASE_NAME")
)

type DbConnection interface {
	ConnectHandle() *sql.DB
}

type PgManager struct {
}

func NewPgManager() *PgManager {
	return &PgManager{}
}

func (p *PgManager) ConnectHandle() *sql.DB {
	db, err := sql.Open(driveName, p.dataSource())
	if err != nil {
		log.Panicln(err)
	}

	return db
}


func (p *PgManager) dataSource() string {
	dbPort, err := strconv.Atoi(port)
	if err != nil {
		log.Panicln(err)
	}
	return fmt.Sprintf(dataSourcePattern, host, dbPort, user, password, dbname)
}
