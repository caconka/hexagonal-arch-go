package main

import (
	"database/sql"
	"flag"
	"fmt"
	"hex-arch/database/psql"
	redisdb "hex-arch/database/redis"
	"hex-arch/domain/ticket"
	"hex-arch/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func main() {

	dbType := flag.String("database", "redis", "database type [redis, psql]")
	flag.Parse()

	var ticketRepo ticket.Repository

	switch *dbType {
	case "psql":
		ticketRepo = psql.NewPostgresTicketRepository(postgresConnection("postgresql://postgres@localhost/ticket?sslmode=disable"))
	case "redis":
		ticketRepo = redisdb.NewRedisTicketRepository(redisConnection("localhost:6379"))
	default:
		panic("Unknown database")
	}

	ticketService := ticket.NewTicketService(ticketRepo)

	svr := server.New(ticketService)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :3000")
		errs <- http.ListenAndServe(":3000", svr)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("terminated %s", <-errs)
}

func redisConnection(url string) *redis.Client {
	fmt.Println("Connecting to Redis DB")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping().Err()

	if err != nil {
		panic(err)
	}
	return client
}

func postgresConnection(database string) *sql.DB {
	fmt.Println("Connecting to PostgreSQL DB")
	db, err := sql.Open("postgres", database)
	if err != nil {
		log.Fatalf("%s", err)
		panic(err)
	}
	return db
}
