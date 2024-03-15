package main

import (
	"fmt"
	"time"
)

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{map[string][]SqueakPost{}}
}

type InMemoryUserStore struct {
	store map[string][]SqueakPost
}

func (i *InMemoryUserStore) GetUserSqueaks(name string) ([]SqueakPost, error) {
	squeaks, ok := i.store[name]
	if !ok {
		return nil, fmt.Errorf("no squeaks found for %s", name)
	}

	return squeaks, nil
}

func (i *InMemoryUserStore) PostSqueak(name, squeak string) (int, error) {
	_, ok := i.store[name]
	if !ok {
		i.store[name] = []SqueakPost{{squeak, time.Now()}}
		return 0, nil
	} else {
		i.store[name] = append(i.store[name], SqueakPost{squeak, time.Now()})
		return 1, nil
	}
}

func (i *InMemoryUserStore) GetUserbase() ([]User, error) {
	var userbase []User
	for name, squeaks := range i.store {
		userbase = append(userbase, User{1, name, "", "", squeaks, time.Now()})
	}
	return userbase, nil
}

func (i *InMemoryUserStore) CreateUser(name, password string) (int, error) { return 0, nil }
