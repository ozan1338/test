package postgresql

import (
	"database/sql"
	"os"
	logger "test/log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var Conn *sql.DB


var dsn string = os.Getenv("DSN")
var count int64
var queryCreateTableUsers = `
CREATE TABLE IF NOT EXISTS public.Users (
	id serial primary key,
	name varchar(50) not null,
	email varchar(250) unique not null,
	password varchar(250) not null
);`


var queryCreateTableJob = `CREATE TABLE IF NOT EXISTS public.job (
	id serial primary key,
	id_user int not null,
	city varchar(50) not null,
	full_time boolean not null,
	description varchar(50) not null,
	foreign key(id_user) references public.Users(id)
	on delete cascade
	on update cascade
);
`


func DatabaseInit() *sql.DB {
	for {
		var err error
		Conn, err = openDB(dsn)
		
		if err != nil {
			logger.Info("Potgres Not Yet Ready...")
			count++
		} else {
			logger.Info("Connected To Postgres!")
			if createTableUserErr := createTable(Conn,queryCreateTableUsers); createTableUserErr != nil {
				return nil
			}
			
			if createTableJobErr := createTable(Conn,queryCreateTableJob); createTableJobErr != nil {
				return nil
			}

			return Conn
		}

		if count > 30 {
			logger.Error("error connecting to DB",err)
			return nil
		}

		logger.Info("Backing off for two seconds..")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx",dsn)
	if err != nil {
		return nil,err
	}

	err = db.Ping()
	if err != nil {
		return nil,err
	}

	return db, nil
}

func createTable(conn *sql.DB, query string) error {
	stmt, err := Conn.Prepare(query)
	// log.Info(queryCreateTable)
	if err != nil {
		logger.Error("error when trying to prepare create table statement",err)
		return err
	}
	defer stmt.Close()

	stmt.Exec()
	return nil
}