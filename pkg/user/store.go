package user

import (
	"log"
	"sync"

	"github.com/google/uuid"
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

func (s *Store) Add(name, email string) uuid.UUID {
	s.Lock()
	defer s.Unlock()

	uuid := uuid.New()
	user := User{
		UUID:  uuid,
		Name:  name,
		Email: email,
	}

	s.user[uuid] = user

	log.Printf("%v", s.user)
	return uuid
}

func (s *Store) Get(uuid uuid.UUID) User {
	s.RLock()
	defer s.RUnlock()

	return s.user[uuid]
}

func (s *Store) GetAll() []User {
	s.RLock()
	defer s.RUnlock()

	users := make([]User, 0)
	for _, v := range s.user {
		users = append(users, v)
	}

	return users
}

func (s *Store) Delete(uuid uuid.UUID) {
	s.Lock()
	defer s.Unlock()

	delete(s.user, uuid)
}
