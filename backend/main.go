package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type tconfig struct {
	Bind  string
	Debug bool
}

func main() {
	var config tconfig
	defaultBind := "127.0.0.1:8001"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		defaultBind = "127.0.0.1:" + envPort
	}
	flag.StringVar(&config.Bind, "bind", defaultBind, "Bind for HTTP requests, PORT env is also supported")
	flag.BoolVar(&config.Debug, "debug", false, "Log more")
	flag.Parse()
	log.Printf("bind %s", config.Bind)

	log.Fatal(http.ListenAndServe(config.Bind, NewRouter()))
}
