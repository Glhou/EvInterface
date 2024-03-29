package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "8080"

func main() {
	fileServer := http.FileServer(http.Dir("./static")) // New code
	http.Handle("/", fileServer)

	http.HandleFunc("/bid", bidHandle)
	http.HandleFunc("/ev", evHanlde)
	http.HandleFunc("/auction", auctionHandle)
	http.HandleFunc("/token", tokenHandle)
	http.HandleFunc("/data", dataHandle)
	http.HandleFunc("/reset", initDB)

	fmt.Printf("Starting server at port %v\n", PORT)
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal(err)
	}
}
