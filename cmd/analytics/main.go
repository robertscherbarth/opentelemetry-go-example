package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/baggage"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"github.com/robertscherbarth/opentelemetry-go-example/pkg/opentelemetry"
	"github.com/robertscherbarth/opentelemetry-go-example/pkg/user"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const serviceName = "analytics"

func main() {

	//init exporter
	tp := opentelemetry.InitJaegerTracerProvider(serviceName)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	router := chi.NewRouter()
	router.Use(otelchi.Middleware("analytics-server", otelchi.WithChiRoutes(router)), middleware.Logger)
	router.Post("/", analyzeUserHandler)

	log.Println("started analytics application")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Printf("error while running server (%s)\n", err.Error())
	}
}

// analyzeUserHandler add additional information to trace object
func analyzeUserHandler(w http.ResponseWriter, r *http.Request) {
	bag := baggage.FromContext(r.Context())
	m := bag.Member("user-id-baggage")
	log.Printf("request user id from baggage: %s\n", m.String())

	_, span := otel.Tracer("analytics").Start(r.Context(), "analytics.user", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	body := r.Body
	var u user.User

	err := json.NewDecoder(body).Decode(&u)
	defer body.Close()
	if err != nil {
		span.RecordError(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	//create a random error
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(10-0) + 0
	if v > 5 {
		span.RecordError(fmt.Errorf("can't create new user in system"))
		span.SetStatus(codes.Error, "")
	} else {
		span.AddEvent("user.new", trace.WithAttributes(attribute.String("name", u.Name)))
	}

	w.WriteHeader(http.StatusOK)
}
