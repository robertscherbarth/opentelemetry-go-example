package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/robertscherbarth/opentelemetry-go-example/pkg/opentelemetry"
)

const serviceName = "consumer"

func main() {
	ctx := context.Background()

	//precondition to set a random number
	rand.Seed(time.Now().UnixNano())

	//init exporter
	tp := opentelemetry.InitTraceProvider(ctx)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	defer func() { _ = tp.Shutdown(ctx) }()

	// Wrap zap logger to extend Zap with API that accepts a context.Context.
	log := otelzap.New(zap.NewExample(), otelzap.WithStackTrace(true))

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	addressURI := "http://localhost:8081/users"

	urlENV, ok := os.LookupEnv("URL")
	if ok {
		addressURI = urlENV
		log.Info("fund new url", zap.String("url", urlENV))
	}

	//rnd get
	go func(context.Context) {
		ticker := time.NewTicker(time.Duration(generateRndInt()) * time.Second)
		for {
			select {
			case <-ticker.C:
				ticker.Reset(time.Duration(generateRndInt()) * time.Second)

				err := rndUserList(ctx, addressURI)
				if err != nil {
					log.Ctx(ctx).Error("oops something went wrong", zap.Error(err))
					break
				}
				log.InfoContext(ctx, "user get all")
			}
		}
	}(ctx)

	go func(context.Context) {
		ticker := time.NewTicker(time.Duration(generateRndInt()) * time.Second)
		for {
			select {
			case <-ticker.C:
				ticker.Reset(time.Duration(generateRndInt()) * time.Second)

				err := rndUserCreate(ctx, addressURI)
				if err != nil {
					log.Ctx(ctx).Error("oops something went wrong", zap.Error(err))
					break
				}
				log.InfoContext(ctx, "user create")
			}
		}
	}(ctx)

	go func(context.Context) {
		ticker := time.NewTicker(time.Duration(generateRndInt()) * time.Second)
		for {
			select {
			case <-ticker.C:
				ticker.Reset(time.Duration(generateRndInt()) * time.Second)

				err := rndUserDelete(ctx, addressURI)
				if err != nil {
					log.Ctx(ctx).Error("oops something went wrong", zap.Error(err))
					break
				}
				log.InfoContext(ctx, "user delete")
			}
		}
	}(ctx)

	log.Info("started consumer application")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Error("error while running server", zap.Error(err))
	}
}

func generateRndInt() int {
	max := 5
	min := 1
	return rand.Intn(max-min) + min
}

// HTTPClientTransporter is a convenience function which helps to attach tracing
// functionality to conventional HTTP clients.
func HTTPClientTransporter(rt http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(rt)
}

func rndUserList(ctx context.Context, addressURI string) error {
	ctx, span := otel.Tracer("").Start(ctx, "user.list", trace.WithTimestamp(time.Now().UTC()))
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addressURI, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}

	res, err := HTTPClientTransporter(http.DefaultTransport).RoundTrip(req)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}
	span.SetStatus(codes.Ok, res.Status)
	return nil
}

func rndUserCreate(ctx context.Context, addressURI string) error {
	ctx, span := otel.Tracer("").Start(ctx, "user.create")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	user := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  "test-user-" + fmt.Sprint(rand.Int()),
		Email: "test-user@example.com",
	}
	jsonData, err := json.Marshal(&user)
	if err != nil {
		span.RecordError(err)
		span.End()
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addressURI, bytes.NewBuffer(jsonData))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}

	res, err := HTTPClientTransporter(http.DefaultTransport).RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.End()
		return err
	}

	span.SetStatus(codes.Ok, res.Status)
	span.End()

	return nil
}

func rndUserDelete(ctx context.Context, addressURI string) error {
	ctx, span := otel.Tracer("").Start(ctx, "user.delete")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addressURI, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}
	res, err := HTTPClientTransporter(http.DefaultTransport).RoundTrip(req)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}

	type user struct {
		UUID uuid.UUID `json:"uuid"`
	}
	users := make([]user, 0)
	err = json.NewDecoder(res.Body).Decode(&users)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}

	//get random id
	id := uuid.New().String()
	if len(users) != 0 {
		id = users[rand.Intn(len(users))].UUID.String()
	}

	delReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, addressURI+"/"+id, nil)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}
	delRes, err := HTTPClientTransporter(http.DefaultTransport).RoundTrip(delReq)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
		return err
	}
	span.AddEvent("user.delete", trace.WithAttributes(
		attribute.String("id", id)))
	span.SetStatus(codes.Ok, delRes.Status)
	return nil

}
