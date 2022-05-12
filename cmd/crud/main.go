package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robertscherbarth/opentelemetry-go-example/pkg/user"
)

func main() {
	store := initUserStore()
	userResource := user.NewUserResource(store)

	router := chi.NewRouter()
	router.Use(middleware.Logger, middleware.StripSlashes)
	router.Mount("/users", userResource.Routes())

	log.Println("started crud application")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Printf("error while running server (%s)\n", err.Error())
	}
}

func initUserStore() *user.Store {
	store := user.NewStore()
	store.Add("tester", "tester@example.com")
	store.Add("tester-1", "tester-1@example.com")

	log.Printf("initialized store user_count: %d", len(store.GetAll()))
	return store
}
