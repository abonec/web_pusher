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
	user := args.Get(0).(testUser)
	return &user, args.Get(1).([]string), args.Error(2)
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
	user, err := server.Auth(&testConn{}, []byte("valid"))
	assert.NoError(err)
	assert.Equal(server.users[testUserId].Size(), 1)
	assert.Equal(server.OnlineConnections(), 1)
	assert.Equal(server.OnlineUsers(), 1)
	set, err := server.channels[mainChannel].GetUserSet(testUserId)
	assert.NoError(err)
	assert.Equal(set.id, testUserId)
	server.Close(user)
	assert.Zero(len(server.channels))
	assert.Zero(server.OnlineUsers())
	assert.Zero(server.OnlineConnections())
}
