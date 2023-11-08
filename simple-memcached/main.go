package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	mc := memcache.New("memcached:11211")

	r := chi.NewRouter()

	r.Get("/get", func(w http.ResponseWriter, r *http.Request) {
		// get number from query parameter
		numStr := r.URL.Query().Get("num")
		if numStr == "" {
			http.Error(w, "num parameter is required", http.StatusBadRequest)
			return
		}

		// get value from memcached
		item, err := mc.Get(numStr)
		if err != nil {
			http.Error(w, "value not found", http.StatusNotFound)
			return
		}

		// write response
		var res Response
		res.Message = string(item.Value)
		json.NewEncoder(w).Encode(res)
	})

	r.Post("/set", func(w http.ResponseWriter, r *http.Request) {
		// get number from body
		var reqBody struct {
			Num int `json:"num"`
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// set value in memcached
		err = mc.Set(&memcache.Item{Key: strconv.Itoa(reqBody.Num), Value: []byte("inserted")})
		if err != nil {
			log.Fatalf("Error setting value: %v", err)
		}

		// write response
		var res Response
		res.Message = "value inserted"
		json.NewEncoder(w).Encode(res)
	})

	http.ListenAndServe(":8080", r)
}
