package web_pusher

import "errors"

type PushChannel interface {
	Id() string
	AddUser(user User)
	RemoveUser(userId string) bool
	GetUser(id string) (User, error)
}

type pushChannel struct {
	id    string
	users map[string]User
}

func (ch *pushChannel) Id() string {
	return ch.id
}

func (ch *pushChannel) AddUser(user User) {
	ch.users[user.Id()] = user
}

func (ch *pushChannel) RemoveUser(userId string) bool {
	delete(ch.users, userId)
	if len(ch.users) == 0 {
		return true
	} else {
		return false
	}
}

var UserNotFound = errors.New("user not found")

func (ch *pushChannel) GetUser(id string) (User, error) {
	user, ok := ch.users[id]
	if ok {
		return user, nil
	} else {
		return nil, UserNotFound
	}
}

func NewPushChannel(firstUser User) PushChannel {
	channel := &pushChannel{
		users: make(map[string]User),
	}
	channel.AddUser(firstUser)
	return channel
}
