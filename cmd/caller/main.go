package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//precondition to set a random number
	rand.Seed(time.Now().UnixNano())
	client := http.DefaultClient

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	//rnd get
	go func() {
		ticker := time.NewTicker(time.Duration(generateRndInt()) * time.Second)
		for {
			select {
			case <-ticker.C:
				addressURI := "http://localhost:8081/users"
				res, err := client.Get(addressURI)
				if err != nil {
					log.Printf("err: %s\n", err.Error())
					continue
				}
				log.Printf("called uri: %s method: GET response_code: %d\n", addressURI, res.StatusCode)
				ticker.Reset(time.Duration(generateRndInt()) * time.Second)
			}
		}
	}()

	//rnd post
	go func() {
		ticker := time.NewTicker(time.Duration(generateRndInt()) * time.Second)
		for {
			select {
			case <-ticker.C:
				addressURI := "http://localhost:8081/users"
				user := struct {
					Name  string `json:"name"`
					Email string `json:"email"`
				}{
					Name:  "test-user-" + fmt.Sprint(rand.Int()),
					Email: "test-user@example.com",
				}
				jsonData, err := json.Marshal(&user)
				if err != nil {
					log.Printf("err: %s\n", err.Error())
					return
				}
				res, err := client.Post(addressURI, "application/json; charset=utf-8", bytes.NewBuffer(jsonData))
				if err != nil {
					log.Printf("err: %s\n", err.Error())
					continue
				}
				log.Printf("called uri: %s method: POST response_code: %d\n", addressURI, res.StatusCode)
				ticker.Reset(time.Duration(generateRndInt()) * time.Second)
			}
		}
	}()

	//rnd delete
	go func() {
		ticker := time.NewTicker(time.Duration(generateRndInt()) * time.Second)
		for {
			select {
			case <-ticker.C:
				addressURI := "http://localhost:8081/users"
				res, err := client.Get(addressURI)
				if err != nil {
					log.Printf("err: %s\n", err.Error())
					continue
				}
				type user struct {
					UUID uuid.UUID `json:"uuid"`
				}
				users := make([]user, 0)
				err = json.NewDecoder(res.Body).Decode(&users)
				if err != nil {
					log.Printf("err: %s\n", err.Error())
					continue
				}
				id := uuid.New().String()
				if len(users) != 0 {
					id = users[rand.Intn(len(users))].UUID.String()
				}
				req, _ := http.NewRequest(http.MethodDelete, addressURI+"/"+id, nil)
				delRes, err := client.Do(req)
				if err != nil {
					log.Printf("err: %s\n", err.Error())
					continue
				}

				log.Printf("called uri: %s method: DELETE response_code: %d\n", addressURI, delRes.StatusCode)
				ticker.Reset(time.Duration(generateRndInt()) * time.Second)
			}
		}
	}()

	log.Println("started caller application")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Printf("error while running server (%s)\n", err.Error())
	}
}

func generateRndInt() int {
	max := 5
	min := 1
	return rand.Intn(max-min) + min
}
