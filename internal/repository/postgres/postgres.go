package postgres

import (
	"fmt"
	"log/slog"
	"task/internal/config"
	"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustInitDatabase(cfg *config.Config, log *slog.Logger) *sqlx.DB {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password, cfg.Db.Database, cfg.Db.Sslmode))
	if err != nil {
		panic(err)
	}
	if err := WaitToDB(db, 15, log); err != nil {
		panic(err)
	}
	return db
}

func WaitToDB(db *sqlx.DB, maxRetries int, log *slog.Logger) error {
	var err error
	for range maxRetries {
		log.Info("Attempting to connect to database")
		if err = db.Ping(); err == nil {
			log.Info("Successfully connection to database")
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}