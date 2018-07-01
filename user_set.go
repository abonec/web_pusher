package web_pusher

import "github.com/abonec/web_pusher/application"

type UserSet struct {
	users     map[application.User]application.User
	id        string
	callbacks []func()
}

func NewUserSet(user application.User) *UserSet {
	set := &UserSet{make(map[application.User]application.User), user.Id(), make([]func(), 0)}
	set.users[user] = user
	return set
}

func (set *UserSet) AddUser(user application.User) {
	set.users[user] = user
}

func (set *UserSet) Size() int {
	return len(set.users)
}
func (set *UserSet) Send(msg []byte) {
	for _, user := range set.users {
		user.Send(msg)
	}
}

func (set *UserSet) DeleteUser(user application.User) {
	delete(set.users, user)
	if len(set.users) == 0 {
		for _, callback := range set.callbacks {
			callback()
		}
	}
}

func (set *UserSet) OnLastConnection(callback func()) {
	set.callbacks = append(set.callbacks, callback)
}
