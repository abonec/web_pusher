package web_pusher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
)

func TestServer_Auth(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(newTestApp())
	server.Start()

	firstConnection, err := server.Auth(&testConn{}, []byte("valid"))
	assert.NoError(err)
	assertOnline(assert, server, 1, 1)
	channel := server.channels[mainChannel]
	assert.Equal(channel.Id(), mainChannel)
	userSet, err := channel.GetUserSet(testUserId)
	assert.NoError(err)
	assert.Equal(userSet.id, testUserId)
	assert.Equal(userSet.Size(), 1)

	invalidUser, err := server.Auth(&testConn{}, []byte("invalid"))
	assert.Error(err)
	assert.Nil(invalidUser)
	assertOnline(assert, server, 1, 1)

	secondConnection, err := server.Auth(&testConn{}, []byte("valid"))
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
