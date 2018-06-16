package web_pusher

import (
	"errors"
)

type PushChannel interface {
	Id() string
	AddUserSet(set *UserSet)
	GetUserSet(id string) (*UserSet, error)
	OnLastUser(callback func())
}

type pushChannel struct {
	id               string
	users            map[string]*UserSet
	lastUserCallback func()
}

func (ch *pushChannel) Id() string {
	return ch.id
}

func (ch *pushChannel) AddUserSet(set *UserSet) {
	if _, present := ch.users[set.id]; present {
		return
	} else {
		ch.users[set.id] = set
	}
	set.OnLastConnection(func() {
		delete(ch.users, set.id)
		if len(ch.users) == 0 {
			if ch.lastUserCallback != nil {
				ch.lastUserCallback()
			}
		}
	})
}

var UserNotFound = errors.New("user not found")

func (ch *pushChannel) GetUserSet(id string) (*UserSet, error) {
	user, ok := ch.users[id]
	if ok {
		return user, nil
	} else {
		return nil, UserNotFound
	}
}

func (ch *pushChannel) OnLastUser(callback func()) {
	ch.lastUserCallback = callback
}

func NewPushChannel(firstUser *UserSet) PushChannel {
	channel := &pushChannel{
		users: make(map[string]*UserSet),
	}
	channel.AddUserSet(firstUser)
	return channel
}
