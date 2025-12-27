package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/geekilx/restaurantAPI/internal/models"

	_ "github.com/lib/pq"
)

const Version = "1.0.0"

type config struct {
	port int

	db struct {
		DSN string
	}
}

type application struct {
	cfg    config
	logger *slog.Logger
	models models.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "the specific port you want to run your program")

	flag.StringVar(&cfg.db.DSN, "dsn", os.Getenv("RESTAURANT_DB_DSN"), "postgres dsn for database")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB(cfg.db.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	app := application{
		cfg:    cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	var server = http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.route(),
	}

	log.Printf("listening server on :%d", cfg.port)
	err = server.ListenAndServe()
	log.Fatalln(err)

}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}
