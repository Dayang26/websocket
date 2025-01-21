package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
	}

	http.ServeFile(w, r, "./examples/chat/home.html")
}

func main() {
	flag.Parse()

	hub := newHub()

	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(hub, writer, request)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("failed!  %v", err)
	}

}
