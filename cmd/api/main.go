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

	"github.com/geekilx/restaurantAPI/internal/mailer"
	"github.com/geekilx/restaurantAPI/internal/models"

	_ "github.com/lib/pq"
)

const Version = "1.0.0"

type config struct {
	port int

	db struct {
		DSN string
	}
	smtp struct {
		host      string
		port      int
		username  string
		passsword string
		sender    string
	}
	limiter struct {
		rps     int
		burst   int
		enabled bool
	}
}

type application struct {
	cfg    config
	logger *slog.Logger
	models models.Models
	mailer mailer.Mailer
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "the specific port you want to run your program")

	flag.StringVar(&cfg.db.DSN, "dsn", os.Getenv("RESTAURANT_DB_DSN"), "postgres dsn for database")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "372d553c29c9c6", "SMTP username")
	flag.StringVar(&cfg.smtp.passsword, "smtp-password", "3fb3fd1b008ee2", "SMTP passsword")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "restaurantAPI <no-reply@restaurantAPI.ilx.net>", "SMTP sender")

	flag.IntVar(&cfg.limiter.rps, "limiter-rps", 2, "rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "enable rate limiter")

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
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.passsword, cfg.smtp.sender),
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
