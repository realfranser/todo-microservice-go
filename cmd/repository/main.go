package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/realfranser/todo-microservice-go/internal/envvar"
	"github.com/realfranser/todo-microservice-go/internal/envvar/vault"
	"github.com/realfranser/todo-microservice-go/internal/postgresql"
	"github.com/realfranser/todo-microservice-go/internal/rest"
	"github.com/realfranser/todo-microservice-go/internal/service"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	if err := envvar.Load(env); err != nil {
		log.Fatalln("Couldn't load configuration", err)
	}

	conf := envvar.New(newVaultProvider())

	//-

	db := newDB(conf)
	defer db.Close()

	//-

	repo := postgresql.NewTask(db) // Task Repository
	svc := service.NewTask(repo)   // Task Application Service

	//-

	app := fiber.New(fiber.Config{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	})

	rest.NewTaskHandler(svc).Register(app)

	address := "0.0.0.0:9234"

	log.Println("Starting server", address)

	log.Fatal(app.Listen(address))
}

func newDB(conf *envvar.Configuration) *sql.DB {
	get := func(v string) string {
		res, err := conf.Get(v)
		if err != nil {
			log.Fatalf("Couldn't get configuration value for %s: %s", v, err)
		}

		return res
	}

	// XXX: We will revisit this code in future episodes replacing it with another solution
	databaseHost := get("DATABASE_HOST")
	databasePort := get("DATABASE_PORT")
	databaseUsername := get("DATABASE_USERNAME")
	databasePassword := get("DATABASE_PASSWORD")
	databaseName := get("DATABASE_NAME")
	databaseSSLMode := get("DATABASE_SSLMODE")
	// XXX: -

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%s", databaseHost, databasePort),
		Path:   databaseName,
	}

	q := dsn.Query()
	q.Add("sslmode", databaseSSLMode)

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		log.Fatalln("Couldn't open DB", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("Couldn't ping DB", err)
	}

	return db
}

func newVaultProvider() *vault.Provider {
	// XXX: We will revisit this code in future episodes replacing it with another solution
	vaultPath := os.Getenv("VAULT_PATH")
	vaultToken := os.Getenv("VAULT_TOKEN")
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	// XXX: -

	provider, err := vault.New(vaultToken, vaultAddress, vaultPath)
	if err != nil {
		log.Fatalln("Couldn't load provider", err)
	}

	return provider
}
