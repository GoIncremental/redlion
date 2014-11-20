package main

import (
	"log"
	"os"

	"github.com/goincremental/dal"
	"github.com/goincremental/web"
)

func main() {
	web.LoadEnv()
	d := dal.NewDAL()
	connection, err := d.Connect(os.Getenv("MONGO_HOSTS"))
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	port := "3200"
	p := os.Getenv("PORT")
	if p != "" {
		port = ":" + p
	}

	router := web.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/api/v1/checkout", Checkout).Methods("POST")

	server := web.NewServer()
	server.Use(web.Gzip())
	server.Use(Mongo(connection))
	server.Use(Render())
	server.Use(Site())
	server.UseHandler(router)
	server.Run(port)
}
