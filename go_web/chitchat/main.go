package main

import "net/http"

func main() {
	p("ChitChat", version(), "Start at", config.Address)

	mux = http.NewServeMux()
	files = http.FileServer(http.Dir(config.static))
	mux.Handle("/static", http.StripPrefix("/static/", files))

	mux.HandleFunc("/")
}
