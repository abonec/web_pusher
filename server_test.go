package web_pusher

import (
	"testing"
	//_ "github.com/stretchr/testify/"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/mock"
)

type testApp struct {
	mock.Mock
}

type testConn struct {
}

type testUser struct {
	id string
}

func (u *testUser) Id() string {
	return u.id
}

func (c *testConn) Send(msg []byte) {
	//c.Called(msg)
}
func (c *testConn) Close(msg []byte) {
	//c.Called(msg)
}

const testUserId = "testUser"
const mainChannel = "main"

var notFound = errors.New("not found")

func (app *testApp) Auth(msg []byte) (User, []string, error) {
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

func TestServer_Auth(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(newTestApp())
	server.Start()

	firstConnection, err := server.Auth(&testConn{}, []byte("valid"))
	assert.NoError(err)
	assertOnline(assert, server, 1, 1)
	userSet, err := server.channels[mainChannel].GetUserSet(testUserId)
	assert.NoError(err)
	assert.Equal(userSet.id, testUserId)
	assert.Len(userSet.users, 1)

	invalidUser, err := server.Auth(&testConn{}, []byte("invalid"))
	assert.Error(err)
	assert.Nil(invalidUser)
	assertOnline(assert, server, 1, 1)


	secondConnection, err := server.Auth(&testConn{}, []byte("valid"))
	assert.NoError(err)
	assert.NotNil(secondConnection)
	assertOnline(assert, server, 1, 2)
	assert.Len(userSet.users, 2)

	server.Close(firstConnection)
	assertOnline(assert, server, 1, 1)
	server.Close(secondConnection)
	assert.Zero(len(server.channels))
	assertOnline(assert, server, 0, 0)
}

func assertOnline(assert *assert.Assertions, server *Server, usersOnline, connectionsOnline int) {
	assert.Equal(server.OnlineConnections(), connectionsOnline)
	assert.Equal(server.OnlineUsers(), usersOnline)
}
