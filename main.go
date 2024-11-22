package main

import (
	"log"
	"net/http"

	"github.com/dgorb/wikitranslate/web"
)

func main() {
    r := web.NewRouter()
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
