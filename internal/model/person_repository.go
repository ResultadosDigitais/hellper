package model

type PersonRepository interface {
	UpdatePersonSlackID(user *User) error
}
