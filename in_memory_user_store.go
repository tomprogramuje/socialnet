package main

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{map[string][]string{}}
}

type InMemoryUserStore struct {
	store map[string][]string
}

func (i *InMemoryUserStore) GetUserSqueaks(name string) []string {
	squeaks, ok := i.store[name]
	if !ok {
		return []string{}
	}

	return squeaks
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

func (i *InMemoryUserStore) GetUserbase() []User {
	var userbase []User
	for name, squeaks := range i.store {
		userbase = append(userbase, User{name, squeaks})
	}
	return userbase
}

func (i *InMemoryUserStore) CreateUser(name string) (int, error) {return 0, nil}