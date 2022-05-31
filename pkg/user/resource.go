package user

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/robertscherbarth/opentelemetry-go-example/pkg/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UsersResource struct {
	store     *Store
	analytics string
}

func NewUserResource(store *Store, analyticsURL string) *UsersResource {
	return &UsersResource{store: store, analytics: analyticsURL}
}

// Routes creates a REST router for the todos resource
func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)    // GET /users - read a list of users
	r.Post("/", rs.Create) // POST /users - create a new user and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Use(rs.userCtx())      // let's have a users map, and let's actually load/manipulate
		r.Get("/", rs.Get)       // GET /users/{id} - read a single user by :id
		r.Delete("/", rs.Delete) // DELETE /users/{id} - delete a single user by :id
	})

	return r
}

func (rs UsersResource) userCtx() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDParam := chi.URLParam(r, "id")
			userID, err := uuid.Parse(userIDParam)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user := rs.store.Get(r.Context(), userID)
			if user.Name == "" {
				http.Error(w, http.StatusText(404), http.StatusNotFound)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (rs UsersResource) List(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, rs.store.GetAll(r.Context()))
}

func (rs UsersResource) Create(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	var u User

	err := json.NewDecoder(body).Decode(&u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id := rs.store.Add(r.Context(), u.Name, u.Email)

	ctx, span := otel.Tracer("user.resource").Start(r.Context(), "provide.analytics", trace.WithSpanKind(trace.SpanKindClient))
	payload, _ := json.Marshal(&u)
	analyticsReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://analytics:8082/", bytes.NewBuffer(payload))
	defer span.End()
	if err != nil {
		span.RecordError(err)
	}
	_, err = opentelemetry.HTTPClientTransporter(http.DefaultTransport).RoundTrip(analyticsReq)
	if err != nil {
		span.RecordError(err)
	}

	render.JSON(w, r, &id)
}

func (rs UsersResource) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(User)
	if !ok {
		http.Error(w, http.StatusText(422), http.StatusUnprocessableEntity)
		return
	}
	render.JSON(w, r, user)
}

func (rs UsersResource) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(User)
	if !ok {
		http.Error(w, http.StatusText(422), http.StatusUnprocessableEntity)
		return
	}
	rs.store.Delete(r.Context(), user.UUID)
	w.WriteHeader(http.StatusNoContent)
}
