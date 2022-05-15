package user

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sync"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type Store struct {
	sync.RWMutex
	user map[uuid.UUID]User
}

func NewStore() *Store {
	return &Store{
		user: make(map[uuid.UUID]User, 0),
	}
}

func (s *Store) Add(ctx context.Context, name, email string) uuid.UUID {
	_, span := otel.Tracer("").Start(ctx, "store.add", trace.WithAttributes(attribute.String("store", "internal")))
	defer span.End()

	s.Lock()
	defer s.Unlock()

	span.AddEvent("creating a user")

	uuid := uuid.New()
	user := User{
		UUID:  uuid,
		Name:  name,
		Email: email,
	}

	span.AddEvent("user successfully created")

	s.user[uuid] = user

	span.AddEvent("added a user to the store", trace.WithAttributes(
		attribute.String("id", uuid.String()),
		attribute.String("name", name),
		attribute.String("email", email),
	))

	return uuid
}

func (s *Store) Get(ctx context.Context, uuid uuid.UUID) User {
	_, span := otel.Tracer("").Start(ctx, "store.get")
	defer span.End()

	s.RLock()
	defer s.RUnlock()

	return s.user[uuid]
}

func (s *Store) GetAll(ctx context.Context) []User {
	_, span := otel.Tracer("").Start(ctx, "store.get.all")
	defer span.End()

	s.RLock()
	defer s.RUnlock()

	users := make([]User, 0)
	for _, v := range s.user {
		users = append(users, v)
	}

	return users
}

func (s *Store) Delete(ctx context.Context, uuid uuid.UUID) {
	_, span := otel.Tracer("").Start(ctx, "store.delete")
	defer span.End()

	s.Lock()
	defer s.Unlock()

	delete(s.user, uuid)
}
