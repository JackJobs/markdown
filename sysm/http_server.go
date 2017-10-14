package sysm

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ListeningTestInterval = 500
	MaxListeningTestCount = 100
)

//定义一个HTTPServer结构体
type HTTPServer struct {
	port int
	//嵌入一个net.Listener接口，任何满足改接口的类型都可以嵌入改字段
	listener net.Listener
}

//构造方法
func NewHTTPServer(port int) *HTTPServer {
	return &HTTPServer{port, nil}
}

func (s *HTTPServer) Addr() string {
	return ":" + strconv.Itoa(s.port)
}

//使用http.Server类型建立一个web服务器，并默认使用ServeHTTP方法作为默认handler
func (s *HTTPServer) ListenAndServe() {
	var err error
	server := &http.Server{
		Addr:           s.Addr(),
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//使用net.listen进行端口监听
	s.listener, err = net.Listen("tcp", s.Addr())
	if err != nil {
		panic(err)
	}
	server.Serve(s.listener)
}

func (s *HTTPServer) Listen() {
	go s.ListenAndServe()

	isListening := make(chan bool)
	go func() {
		result := false
		ticker := time.NewTicker(time.Microsecond * ListeningTestInterval)
		for i := 0; i < MaxListeningTestCount; i++ {
			<-ticker.C
			resp, err := http.Get("http://localhost" + s.Addr() + "/ping")
			if err == nil && resp.StatusCode == 200 {
				result = true
				break
			}
		}
		ticker.Stop()
		isListening <- result
	}()

	if <-isListening {
		fmt.Println("Listening", s.Addr(), "...")
	} else {
		panic("Can't connect to server")
	}
}

//默认的请求处理函数
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:] //移除'/'
	if path == "ping" {
		w.Write([]byte("pong"))
		fmt.Println("accept connection")
	} else if isWebsocketRequest(r) {
		fmt.Println("websocket connect...")
		NewWebsocket().Serve(w, r) //创建websocket连接，发送数据
	} else {
		Template(w, s.port)
	}
}

func (s *HTTPServer) Stop() {
	s.listener.Close()
}

func contains(arr []string, needle string) bool {
	for _, v := range arr {
		if strings.Contains(v, needle) {
			return true
		}
	}
	return false
}

//判断是否是websocket请求
func isWebsocketRequest(r *http.Request) bool {
	upgrade := r.Header["Upgrade"]
	connection := r.Header["Connection"]
	return contains(upgrade, "websocket") && contains(connection, "Upgrade")
}
