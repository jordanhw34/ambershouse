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

func NewPostgresRepo(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &dbPostresRepo{
		App: app,
		DB:  db,
	}
}

// func NewMySqlRep(db *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
// 	return &dbPostresRepo{
// 		App: app,
// 		DB:  db,
// 	}
// }
