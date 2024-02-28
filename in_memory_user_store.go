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

func (i *InMemoryUserStore) PostSqueak(name, squeak string) int {
	_, ok := i.store[name]
	if !ok {
		i.store[name] = []string{squeak}
		return 0
	} else {
		i.store[name] = append(i.store[name], squeak)
		return 1
	}
}

func (i *InMemoryUserStore) GetUserbase() []User {
	var userbase []User
	for name, squeaks := range i.store {
		userbase = append(userbase, User{name, squeaks})
	}
	return userbase
}
