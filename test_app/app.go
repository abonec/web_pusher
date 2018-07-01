package test_app

import (
	"github.com/stretchr/testify/mock"
	"errors"
	"github.com/abonec/web_pusher/application"
)

type testApp struct {
	mock.Mock
}

type TestConn struct {
	mock.Mock
}

type testUser struct {
	id string
}

func (u *testUser) Id() string {
	return u.id
}

func (c *TestConn) Send(msg []byte) bool {
	c.Called(msg)
	return true
}
func (c *TestConn) Close(msg []byte) {
	//c.Called(msg)
}

const TestUserId = "testUser"
const MainChannel = "main"

var notFound = errors.New("not found")

func (app *testApp) Auth(msg []byte) (application.AuthUser, []string, error) {
	args := app.Called(msg)
	u := args.Get(0)
	ch := args.Get(1)
	var user testUser
	var channels []string
	if u != nil {
		user = u.(testUser)
	}
	if ch != nil {
		channels = ch.([]string)
	}
	return &user, channels, args.Error(2)
}

func NewTestApp() *testApp {
	app := &testApp{}
	app.On("Auth", []byte("valid")).Return(testUser{id: TestUserId}, []string{MainChannel}, nil)
	app.On("Auth", []byte("invalid")).Return(nil, nil, notFound)
	return app
}
