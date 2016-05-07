package main

import (
	"github.com/codegangsta/negroni"
	"github.com/l-lin/pikachu/web"
	"log"
	"os"
)

func main() {
	app := negroni.Classic()

	router := web.NewRouter()

	app.UseHandler(router)
	app.Run(port())
}

func port() string {
	port := os.Getenv("PIKA_PORT")
	if port == "" {
		port = "3000"
		log.Println("[-] No PIKA_PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}
