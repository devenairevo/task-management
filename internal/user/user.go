package user

import (
	"errors"
)

type User struct {
	ID   int
	Name string
}

func New(id int, name string) (*User, error) {
	if id > 0 && name != "" {
		return &User{
			ID:   id,
			Name: name,
		}, nil
	}

	return nil, errors.New("couldn't create a new user")

}
