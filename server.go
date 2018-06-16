package web_pusher

import (
	"errors"
	_ "github.com/gorilla/websocket"
	_ "github.com/stretchr/objx"
)

type Application interface {
	Auth(msg []byte) (User, []string, error)
}

type Join struct {
	user     User
	channels *[]string
	result   chan interface{}
}
type Leave struct {
	user   User
	result chan interface{}
}

type Server struct {
	initialized      bool
	app              Application
	users            map[string]*UserSet
	channels         map[string]PushChannel
	joins            chan Join
	leaves           chan Leave
	onlineConnection int
}

func NewServer(app Application) *Server {
	return &Server{
		initialized: true,
		app:         app,
		users:       make(map[string]*UserSet),
		channels:    make(map[string]PushChannel),
		joins:       make(chan Join),
		leaves:      make(chan Leave),
	}
}

func (s *Server) OnlineUsers() int {
	return len(s.users)
}

func (s *Server) OnlineConnections() int {
	return s.onlineConnection
}

func (s *Server) sendToUser(userId string, msg []byte) {
}

func (s *Server) sendToChannel(channelId string, msg []byte) {
}

func (s *Server) SendById(id string) {
}

func (s *Server) Auth(conn Connection, msg []byte) (User, error) {
	u, channels, err := s.app.Auth(msg)
	if err != nil {
		conn.Close(nil)
		return nil, err
	}

	result := make(chan interface{})
	user := NewUser(u.Id(), conn)
	join := Join{user, &channels, result}
	s.joins <- join
	<-result
	return user, nil
}

func (s *Server) Close(user User) {
	result := make(chan interface{})
	s.leaves <- Leave{user, result}
	<-result
}

func (s *Server) Start() error {
	if !s.initialized {
		return errors.New("server should be initialized")
	}

	go func() {
		for {
			select {
			case join := <-s.joins:
				s.join(join)
			case leave := <-s.leaves:
				s.leave(leave)
			}
		}
	}()
	return nil
}

func (s *Server) join(join Join) {
	defer close(join.result)
	userSet, ok := s.users[join.user.Id()]
	if !ok {
		userSet = NewUserSet(join.user)
		userSet.OnLastConnection(func(){
			delete(s.users, join.user.Id())
		})
		s.users[join.user.Id()] = userSet
	}else{
		userSet.AddUser(join.user)
	}
	s.onlineConnection += 1
	for _, channel := range *join.channels {
		pushChannel, ok := s.channels[channel]
		if ok {
			pushChannel.AddUserSet(userSet)
		} else {
			pushChannel = NewPushChannel(userSet)
			s.channels[channel] = pushChannel
			pushChannel.OnLastUser(func(ch string) func() {
				return func() {
					delete(s.channels, ch)
				}
			}(channel))
		}
	}
}

func (s *Server) leave(leave Leave) {
	defer close(leave.result)
	if set, ok := s.users[leave.user.Id()]; ok {
		set.DeleteUser(leave.user)
	}
	s.onlineConnection -= 1
}
