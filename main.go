package main

import (
	"database/sql"
	"flag"
	"fmt"
	"hex-arch/database/mongo"
	"hex-arch/database/psql"
	redisdb "hex-arch/database/redis"
	"hex-arch/domain/ticket"
	"hex-arch/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/globalsign/mgo"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func main() {

	dbType := flag.String("database", "redis", "database type [redis, psql, mongo]")
	flag.Parse()

	var ticketRepo ticket.Repository

	switch *dbType {
	case "psql":
		c := postgresConnection("postgresql://postgres@localhost/ticket?sslmode=disable")
		ticketRepo = psql.NewPostgresTicketRepository(c)
	case "redis":
		c := redisConnection("localhost:6379")
		ticketRepo = redisdb.NewRedisTicketRepository(c)
	case "mongo":
		c := mongoConnection("localhost:27017")
		ticketRepo = mongo.NewMongoTicketRepository(c)
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
		Password: "",
		DB:       0,
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

func mongoConnection(url string) *mgo.Session {
	fmt.Println("Connecting to Mongo DB")
	session, err := mgo.Dial(url)

	if err != nil {
		log.Fatalf("%s", err)
		panic(err)
	}

	return session
}
