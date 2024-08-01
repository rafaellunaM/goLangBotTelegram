package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "app example")
	})

	port := ":8080"
	fmt.Printf("porta%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %s\n", err)
	}
}
