package main

import (
	"github.com/codegangsta/negroni"
	"github.com/l-lin/pikachu/web"
	"log"
	"os"
	"bufio"
)

func main() {
	displayBanner()

	app := negroni.Classic()

	router := web.NewRouter()

	app.UseHandler(router)
	app.Run(port())
}

func port() string {
	port := os.Getenv("PIKA_PORT")
	if port == "" {
		port = "3000"
		log.Println("[-] No PIKA_PORT environment variable detected. Setting to", port)
	}
	return ":" + port
}

func displayBanner() {
	file, err := os.Open("banner.txt")
	if err != nil {
		log.Fatal("Could not open file 'banner.txt'")
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	log.Println("-------------------------------------")
	log.Println("PIKACHU! I CHOOSE YOU!")
	log.Println("-------------------------------------")
	for scan.Scan() {
		log.Println(scan.Text())
	}
}
