package dbrepo

import (
	"database/sql"

	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/repository"
)

type dbPostresRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type dbTestRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates a new repository to hold the connection to our Postgres Database
func NewPostgresRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &dbPostresRepo{
		App: app,
		DB:  db,
	}
}

// NewDBTestRepo creates a new repository to hold the connection for testing
func NewDBTestRepo(app *config.AppConfig) repository.DatabaseRepo {
	return &dbTestRepo{
		App: app,
	}
}

// func NewMySqlRep(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
// 	return &dbPostresRepo{
// 		App: app,
// 		DB:  db,
// 	}
// }
