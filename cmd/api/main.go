package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/geekilx/restaurantAPI/internal/mailer"
	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/redis/go-redis/v9"

	"github.com/kelseyhightower/envconfig"

	_ "github.com/lib/pq"
)

const Version = "1.0.0"

type config struct {
	Port int `envconfig:"PORT"`

	DB struct {
		DSN string `envconfig:"RESTAURANT_DB_DSN"`
	}
	Smtp struct {
		Host      string `envconfig:"SMTP_HOST"`
		Port      int    `envconfig:"SMTP_PORT"`
		Username  string `envconfig:"SMTP_USERNAME"`
		Passsword string `envconfig:"SMTP_PASSWORD"`
		Sender    string `envconfig:"SMTP_SENDER"`
	}
	Limiter struct {
		Rps     int  `envconfig:"LIMITER_RPS"`
		Burst   int  `envconfig:"LIMITER_BURST"`
		Enabled bool `envconfig:"LIMITER_ENABLED"`
	}
	Redis struct {
		Addr string `envconfig:"REDIS_ADDR"`
	}
}

type application struct {
	cfg    config
	logger *slog.Logger
	models models.Models
	mailer mailer.Mailer
	redis  *redis.Client
	wg     sync.WaitGroup
}

// @title Restaurant API
// @version 1.0
// @description This is a sample server for a restaurant management system.
// @host localhost:4000
// @BasePath /v1
func main() {
	var cfg config

	// flag.IntVar(&cfg.port, "port", 4000, "the specific port you want to run your program")

	// flag.StringVar(&cfg.db.DSN, "dsn", os.Getenv("RESTAURANT_DB_DSN"), "postgres dsn for database")

	// flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	// flag.IntVar(&cfg.smtp.port, "smtp-port", os.Getenv("SMTP_PORT"), "SMTP port")
	// flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	// flag.StringVar(&cfg.smtp.passsword, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP passsword")
	// flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	// flag.IntVar(&cfg.limiter.rps, "limiter-rps", 2, "rate limiter maximum request per second")
	// flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "rate limiter maximum burst")
	// flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "enable rate limiter")

	// flag.StringVar(&cfg.redis.addr, "redis-addr", os.Getenv("REDIS_ADDR"), "redis address")

	// flag.Parse()

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB(cfg.DB.DSN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	rdb, err := openRedis(cfg.Redis.Addr)
	if err != nil {
		logger.Error("failed to connect to redis", "Error", err)
		os.Exit(1)
	}

	app := application{
		cfg:    cfg,
		logger: logger,
		models: models.NewModels(db),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Passsword, cfg.Smtp.Sender),
		redis:  rdb,
	}

	err = app.serve()
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

func openRedis(addr string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, err
	}

	return rdb, nil
}
