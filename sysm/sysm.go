package sysm

type Sysm struct {
	port       int
	httpServer *HTTPServer
	stop       chan bool
}

//构造函数
func NewSysm(port int) *Sysm {
	return &Sysm{port, nil, make(chan bool)}
}

//开始运行
func (s *Sysm) Run() {
	s.httpServer = NewHTTPServer(s.port)
	s.httpServer.Listen()
	<-s.stop
}

func (s *Sysm) Stop() {
	s.httpServer.Stop()
	s.stop <- true
}
