package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/taxintt/go-playgrounds/ent-with-atlas/ent"
	"github.com/taxintt/go-playgrounds/ent-with-atlas/ent/user"

	_ "github.com/lib/pq"
)

type server struct {
	client *ent.Client
}

func newServer(client *ent.Client) *server {
	return &server{client: client}
}

func main() {
	client, err := ent.Open("postgres", "host=postgres port=5432 user=postgres dbname=testdb password=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool.
	// https://entgo.io/ja/docs/getting-started/#automatic-migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	srv := newServer(client)
	r := chi.NewRouter()

	r.Route("/user", func(r chi.Router) {
		r.Get("/{userID}", srv.queryUser)
		r.Post("/", srv.createUser)
	})

	fmt.Println("Server running")
	defer fmt.Println("Server stopped")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func (s *server) queryUser(w http.ResponseWriter, r *http.Request) {
	userIDString := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := s.client.User.Query().Where(user.ID(userID)).Only(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "user returned: %v", u)
}

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var user ent.User
	json.Unmarshal(body, &user)
	u, err := s.client.User.Create().SetAge(user.Age).SetName(user.Name).Save(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "user was created: %v", u)
}
