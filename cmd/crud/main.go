package main

import (
	"context"
	"github.com/robertscherbarth/opentelemetry-go-example/pkg/opentelemetry"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"

	"github.com/robertscherbarth/opentelemetry-go-example/pkg/user"
)

const serviceName = "user-management"

func main() {

	//init exporter
	tp := opentelemetry.InitJaegerTracerProvider(serviceName)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	store := initUserStore()
	userResource := user.NewUserResource(store)

	router := chi.NewRouter()
	router.Use(otelchi.Middleware("", otelchi.WithChiRoutes(router)))
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := trace.SpanContextFromContext(r.Context())
			log.Println(s.TraceID())
			for k, v := range r.Header {
				log.Printf("header - key: %s value: %s\n", k, v)
			}
			handler.ServeHTTP(w, r)
		})
	})
	router.Use(middleware.Logger, middleware.StripSlashes)
	router.Mount("/users", userResource.Routes())

	log.Println("started crud application")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Printf("error while running server (%s)\n", err.Error())
	}
}

func initUserStore() *user.Store {
	store := user.NewStore()
	ctx := context.Background()

	store.Add(ctx, "tester", "tester@example.com")
	store.Add(ctx, "tester-1", "tester-1@example.com")

	log.Printf("initialized store user_count: %d", len(store.GetAll(ctx)))
	return store
}
