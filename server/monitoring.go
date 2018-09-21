package server

type Monitoring interface {
	JoinActor()
	JoinConnection()
	LeaveActor()
	LeaveConnection()
	SendToActor(int)
	SendToConnection(int)
}

func (s *Server) SetMonitoring(m Monitoring) {
	s.monitoring = m
}

type blankMonitoring struct {
}

func (_ blankMonitoring) SendToActor(int) {
}

func (_ blankMonitoring) SendToConnection(int) {
}

func (_ blankMonitoring) JoinActor() {
}

func (_ blankMonitoring) JoinConnection() {
}

func (_ blankMonitoring) LeaveActor() {
}

func (_ blankMonitoring) LeaveConnection() {
}
