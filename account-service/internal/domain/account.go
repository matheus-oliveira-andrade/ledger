package domain

import "time"

type Account struct {
	Id        string
	Name      string
	Document  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAccount(name, document string) *Account {
	return &Account{
		Id:        "",
		Name:      name,
		Document:  document,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
