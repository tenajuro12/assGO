package main

import (
	"context"      // New import
	"database/sql" // New import
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/tenajuro12/assGO/internal/data"
	"github.com/tenajuro12/assGO/internal/mailer"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type application struct {
	config config
	logger *log.Logger
	models data.Models
	model2 data.TeacherModel
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:123456@localhost/t.SabyrovDB", "PostgreSQL DSN")

	// with your own Mailtrap credentials.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "e15bfd94a62406", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "3f62365712badd", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		model2: data.NewModel1(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {

	db, err := sql.Open("postgres", "postgres://postgres:123456@localhost/t.SabyrovDB?sslmode=disable")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil

}
