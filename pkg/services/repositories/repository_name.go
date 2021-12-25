package repositories

import (
	"database/sql"
)

type Name interface {
	GetData() ([]struct{}, error)
	SaveData([]struct{}) error
}

type name struct{}

func NewName(db *sql.DB) Name {
	return &name{}
}

func (n *name) GetData() ([]struct{}, error) {
	return nil, nil
}

func (n *name) SaveData([]struct{}) error {
	return nil
}
