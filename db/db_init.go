package db

import (
	"database/sql"
	"news/config"
	"news/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

var ReformDB *reform.DB

func InitDB() {

	cfg := config.LoadConfig()

	db, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		utils.Log.Fatal("Cannot connect to database:", err)
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DBSource,
	)
	if err != nil {
		utils.Log.Fatal("Failed to migrate", err)
	}

	// Выполнение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		utils.Log.Fatal("Failed to up migrations", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	ReformDB = reform.NewDB(db, postgresql.Dialect, reform.NewPrintfLogger(utils.Log.Printf))
}
