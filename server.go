package web_pusher

import (
	_ "github.com/gorilla/websocket"
	_ "github.com/stretchr/objx"
	"errors"
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
	userId string
	result chan interface{}
}
type Server struct {
	initialized bool
	app         Application
	users       map[string]User
	channels    map[string]PushChannel
	joins       chan Join
	leaves      chan Leave
}

func NewServer(app Application) *Server {
	return &Server{
		initialized: true,
		app:         app,
		users:       make(map[string]User),
		channels:    make(map[string]PushChannel),
		joins:       make(chan Join),
		leaves:      make(chan Leave),
	}
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
	join := Join{NewUser(u.Id(), conn), &channels, result}
	s.joins <- join
	<-result
	return u, nil
}

func (s *Server) Close(user User) {
	result := make(chan interface{})
	s.leaves <- Leave{user.Id(), result}
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
	for _, channel := range *join.channels {
		pushChannel, ok := s.channels[channel]
		if ok {
			pushChannel.AddUser(join.user)
		} else {
			s.channels[channel] = NewPushChannel(join.user)
		}
	}
	s.users[join.user.Id()] = join.user
}

func (s *Server) leave(leave Leave) {
	defer close(leave.result)
	delete(s.users, leave.userId)
	deleteChannels := make([]string, len(s.channels))
	for name, channel := range s.channels {
		if channel.RemoveUser(leave.userId) {
			deleteChannels = append(deleteChannels, name)
		}
	}

	for _, channel := range deleteChannels {
		delete(s.channels, channel)
	}
}
