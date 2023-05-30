package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store interface {
	CreateTransfer() error
}

type Storage struct {
	DB *sql.DB
}

func NewStorage(dbUser,
	dbPassword,
	dbHost,
	dbPort,
	dbName string,
) (*Storage, error) {
	dataSource := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: db,
	}, nil
}

func (store *Storage) CloseDBConnection() error {
	if err := store.DB.Close(); err != nil {
		return fmt.Errorf("error close database: %w", err)
	}

	return nil
}

func (store *Storage) CreateTransfer() error {

	return nil
}
