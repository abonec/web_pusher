package web_pusher

type UserSet struct {
	users     map[User]User
	id        string
	callbacks []func()
}

func NewUserSet(user User) *UserSet {
	set := &UserSet{make(map[User]User), user.Id(), make([]func(), 0)}
	set.users[user] = user
	return set
}

func (set *UserSet) AddUser(user User) {
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

func (set *UserSet) DeleteUser(user User) {
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
