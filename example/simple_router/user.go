package main

import (
	"math/rand"
	"strconv"
)

type User struct {
	id string
}

func (u *User) Id() string {
	return u.id
}
func (u *User) Send(msg []byte) bool {
	return true
}

func NewUser() *User {
	return &User{id: strconv.Itoa(rand.Int())}
}
