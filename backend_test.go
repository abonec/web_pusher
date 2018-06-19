package web_pusher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"bytes"
)

func TestBackend_sendToUser(t *testing.T) {
	assert := assert.New(t)
	server := NewServer(newTestApp())
	server.Start()
	assert.NotNil(1)
	backend := NewBackend(server, nil)


	req := newSendToUserRequest(t, "testUser", []byte{})
	rr := httptest.NewRecorder()

	backend.SendToUser(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)

	req = newSendToUserRequest(t, "testUser", []byte("{}"))
	rr = httptest.NewRecorder()
	backend.SendToUser(rr, req)
	assert.Equal(rr.Code, http.StatusOK)

	testConn := &testConn{}
	conn, err := server.Auth(testConn, []byte("valid"))
	assert.NotNil(conn)
	assert.NoError(err)

	msg := []byte("{}")
	testConn.On("Send", msg)
	rr = httptest.NewRecorder()
	backend.SendToUser(rr, newSendToUserRequest(t, "testUser", msg))
	server.Close(conn)
	backend.SendToUser(rr, newSendToUserRequest(t, "testUser", msg))
	testConn.AssertNumberOfCalls(t, "Send", 1)
	testConn.AssertCalled(t, "Send", msg)
}

func newSendToUserRequest(t *testing.T, userId string, msg []byte) *http.Request {
	req, err := http.NewRequest("POST", "/send-to-user/"+userId, bytes.NewReader(msg))
	if err != nil {
		t.Fatal(err)
	}
	return req
}
