package main

import "fmt"

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{map[string][]string{}}
}

type InMemoryUserStore struct {
	store map[string][]string
}

func (i *InMemoryUserStore) GetUserSqueaks(name string) ([]string, error) {
	squeaks, ok := i.store[name]
	if !ok {
		return nil, fmt.Errorf("no squeaks found for %s", name)
	}

	return squeaks, nil
}

func (i *InMemoryUserStore) PostSqueak(name, squeak string) (int, error) {
	_, ok := i.store[name]
	if !ok {
		i.store[name] = []string{squeak}
		return 0, nil
	} else {
		i.store[name] = append(i.store[name], squeak)
		return 1, nil
	}
}

func (i *InMemoryUserStore) GetUserbase() ([]User, error) {
	var userbase []User
	for name, squeaks := range i.store {
		userbase = append(userbase, User{1, name, "", "", squeaks})
	}
	return userbase, nil
}

func (i *InMemoryUserStore) CreateUser(name, password string) (int, error) { return 0, nil }
