package server

import (
	"github.com/abonec/web_pusher/test_app"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"testing"
)

func TestServer_Auth(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(test_app.NewTestApp())
	server.Start()

	firstConnection, err := server.Auth(&test_app.TestConn{}, []byte("valid"))
	assert.NoError(err)
	assertOnline(assert, server, 1, 1)
	channel := server.channels[test_app.MainChannel]
	assert.Equal(channel.Id(), test_app.MainChannel)
	userSet, err := channel.GetUserSet(test_app.TestUserId)
	assert.NoError(err)
	assert.Equal(userSet.id, test_app.TestUserId)
	assert.Equal(userSet.Size(), 1)

	invalidUser, err := server.Auth(&test_app.TestConn{}, []byte("invalid"))
	assert.Error(err)
	assert.Nil(invalidUser)
	assertOnline(assert, server, 1, 1)

	secondConnection, err := server.Auth(&test_app.TestConn{}, []byte("valid"))
	assert.NoError(err)
	assert.NotNil(secondConnection)
	assertOnline(assert, server, 1, 2)
	assert.Equal(userSet.Size(), 2)

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
