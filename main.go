package main

import "net/http"

func main() {
	serverHandler := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: serverHandler,
	}

	serverHandler.Handle("/", http.FileServer(http.Dir("./")))

	server.ListenAndServe()
}
