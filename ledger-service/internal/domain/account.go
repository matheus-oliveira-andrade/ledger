package domain

import "github.com/golang/protobuf/ptypes/timestamp"

type Account struct {
	Id        int64                `json:"id,omitempty"`
	Document  string               `json:"document,omitempty"`
	Name      string               `json:"name,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"createdAt,omitempty"`
}

func NewAccount(id int64, document, name string, createdAt *timestamp.Timestamp) *Account {
	return &Account{
		Id:        id,
		Document:  document,
		Name:      name,
		CreatedAt: createdAt,
	}
}
