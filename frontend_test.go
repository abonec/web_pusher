package web_pusher

import (
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"github.com/gorilla/websocket"
	"strings"
)

type testHandler struct {
	front *Frontend
}

type testServer struct {
	*httptest.Server
	url       string
	appServer *Server
}

func (h *testHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.front.Handle(rw, r)
}

func newFrontendServer(assert *assert.Assertions) *testServer {
	testSrv := &testServer{}
	app := newTestApp()
	appServer := NewServer(app)
	err := appServer.Start()
	assert.NoError(err)
	front := NewFrontend(appServer, nil)
	testSrv.Server = httptest.NewServer(&testHandler{front})
	testSrv.URL = makeWsProto(testSrv.Server.URL)
	testSrv.appServer = appServer
	return testSrv
}

var testDialer = websocket.Dialer{}

func TestSuccessAuth(t *testing.T) {
	assert := assert.New(t)
	server := newFrontendServer(assert)

	ws, _, err := testDialer.Dial(server.URL, http.Header{})
	assert.NoError(err)
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("valid"))
	assert.NoError(err)
	var msg FrontendMessage
	err = ws.ReadJSON(&msg)
	assert.NoError(err)
	assert.Equal("auth", msg.MsgType)
	assert.Equal("success", msg.AuthStatus)
	assert.Equal(1, server.appServer.OnlineUsers())
	err = ws.WriteMessage(websocket.CloseMessage, []byte{})
	assert.NoError(err)
	ws.ReadMessage()
	assert.Equal(0, server.appServer.OnlineUsers())
}

func TestFailureAuth(t *testing.T) {
	assert := assert.New(t)
	server := newFrontendServer(assert)

	ws, _, err := testDialer.Dial(server.URL, http.Header{})
	assert.NoError(err)
	defer ws.Close()
	err = ws.WriteMessage(websocket.TextMessage, []byte("invalid"))
	assert.NoError(err)
	var msg FrontendMessage
	err = ws.ReadJSON(&msg)
	assert.NoError(err)
	assert.Equal("auth", msg.MsgType)
	assert.Equal("failure", msg.AuthStatus)
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}
