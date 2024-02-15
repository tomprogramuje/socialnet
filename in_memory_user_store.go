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

func (i *InMemoryUserStore) PostSqueak(name, squeak string) {
	_, ok := i.store[name]
	if !ok {
		i.store[name] = []string{squeak}
	} else {
		i.store[name] = append(i.store[name], squeak)
	}
}

func (i *InMemoryUserStore) GetUserbase() []User {
	var userbase []User
	for name, squeaks := range i.store {
		userbase = append(userbase, User{name, squeaks})
	}
	return userbase
}
