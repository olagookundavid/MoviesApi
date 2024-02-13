/*
mkdir -p bin cmd/api internal migrations remote
cmd/api/main.go
*/
package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/olagookundavid/movie-api/internal/data"
	"github.com/olagookundavid/movie-api/internal/jsonlog"
	"github.com/olagookundavid/movie-api/internal/mailer"
	"github.com/olagookundavid/movie-api/internal/vcs"
)

// Declare a string containing the application version number. Later in the book we'll // generate this automatically at build time, but for now we'll just store the version // number as a hard-coded global constant.
var (
	version = vcs.Version()
)

// Define a config struct to hold all the configuration settings for our application. // For now, the only configuration settings will be the network port that we want the // server to listen on, and the name of the current operating environment for the // application (development, staging, production, etc.). We will read in these // configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env variable missing")
	}
	//env and port
	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	//db
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbUrl, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	//limiters
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	//smpt
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "eff608bdd9a4ff", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "e3d9c7596dda6a", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.goliath.net>", "SMTP sender")
	//cors
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()
	// If the version flag value is true, then print out the version number and immediately exit.
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)

	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
echo 'bin/' >> .gitignore put bin in gitignore
goose postgres postgres://greenlight:greenlight@localhost/greenlight up
can omit struct fields in json response completely or if zero pg 50

envelope json msgs: pg54

custom formatting for json if you need it for a field: pg61 and inv pg95
//err := json.NewEncoder(w).Encode(data) pg44
//this methods writes to output stream in a single step


db setup: pg105,108 pg cmds
postgres://greenlight:greenlight@localhost/greenlight



//for permission error
CREATE DATABASE EXAMPLE_DB;
CREATE USER EXAMPLE_USER WITH ENCRYPTED PASSWORD 'Sup3rS3cret';
GRANT ALL PRIVILEGES ON DATABASE EXAMPLE_DB TO EXAMPLE_USER;
\c EXAMPLE_DB postgres
# You are now connected to database "EXAMPLE_DB" as user "postgres".
GRANT ALL ON SCHEMA public TO EXAMPLE_USER;

race condition in updating pg 171

//user table

curl-
patch - curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4
same for del and put
can put -i for more detail
post - BODY='{"title":"Moana","year":2016,"runtime":"107 mins", "genres":["animation","adventure"]}' && curl -i -d "$BODY" localhost:4000/v1/movies
*/
