package web_pusher

import (
	"github.com/stretchr/testify/mock"
	"errors"
)

type testApp struct {
	mock.Mock
}

type testConn struct {
	mock.Mock
}

type testUser struct {
	id string
}

func (u *testUser) Id() string {
	return u.id
}

func (c *testConn) Send(msg []byte) bool {
	c.Called(msg)
	return true
}
func (c *testConn) Close(msg []byte) {
	//c.Called(msg)
}

const testUserId = "testUser"
const mainChannel = "main"

var notFound = errors.New("not found")

func (app *testApp) Auth(msg []byte) (AuthUser, []string, error) {
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

func newTestApp() *testApp {
	app := &testApp{}
	app.On("Auth", []byte("valid")).Return(testUser{id: testUserId}, []string{mainChannel}, nil)
	app.On("Auth", []byte("invalid")).Return(nil, nil, notFound)
	return app
}
