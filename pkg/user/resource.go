package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type UsersResource struct {
	store *Store
}

func NewUserResource(store *Store) *UsersResource {
	return &UsersResource{store: store}
}

// Routes creates a REST router for the todos resource
func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /users - read a list of users
	r.Post("/", rs.Create) // POST /users - create a new user and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Use(rs.userCtx())      // lets have a users map, and lets actually load/manipulate
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
			user := rs.store.Get(userID)
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
	render.JSON(w, r, rs.store.GetAll())
}

func (rs UsersResource) Create(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	var u User

	err := json.NewDecoder(body).Decode(&u)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	rs.store.Add(u.Name, u.Email)

	w.Write([]byte("users create"))
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
	rs.store.Delete(user.UUID)
	w.WriteHeader(http.StatusNoContent)
}
