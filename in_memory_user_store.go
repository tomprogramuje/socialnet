package main

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{map[string]int{}}
}

type InMemoryUserStore struct {
	store map[string]int
}

func (i *InMemoryUserStore) GetUserSqueakCount(name string) int {
	return i.store[name]
}

func (i *InMemoryUserStore) PostSqueak(name string) {
	i.store[name]++
}
