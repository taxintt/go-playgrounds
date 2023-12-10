package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	//r.Use(middleware.DefaultCompress) //using this produces the same result
	r.Use(middleware.Compress(5))

	r.Get("/", Hello)
	http.ListenAndServe(":3333", r)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") //according to the documentation this must be here to enable gzip

	randomString := generateRandomString(1000) // 2000 bytes
	w.Write([]byte(randomString))
}

func generateRandomString(length int) string {
	randomString := make([]byte, length)
	rand.Read(randomString)
	return fmt.Sprintf("%x", randomString)
}
