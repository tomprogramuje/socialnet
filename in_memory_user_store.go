package main

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{map[string]int{}, []User{}}
}

type InMemoryUserStore struct {
	store    map[string]int
	userbase []User
}

func (i *InMemoryUserStore) GetUserSqueakCount(name string) int {
	return i.store[name]
}

func (i *InMemoryUserStore) PostSqueak(name string) {
	i.store[name]++
}

func (i *InMemoryUserStore) GetUserbase() []User {
	var userbase []User
	for _, users := range i.userbase {
		userbase = append(userbase, User{users.Name, users.Squeaks})
	}
	return userbase
}
