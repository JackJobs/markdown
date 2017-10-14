package previewer

import (
	"fmt"
	goWs "github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	WriteTimeout = 5 * time.Second
	BufferSize   = 2048
)

var upgrader = goWs.Upgrader{
	ReadBufferSize:  BufferSize,
	WriteBufferSize: BufferSize,
}

type Websocket struct {
	watcher *Watcher
}

func NewWebsocket(path string) *Websocket {
	return &Websocket{NewWatcher(path)}
}

func (ws *Websocket) Reader(c *goWs.Conn, closed chan<- bool) {
	defer c.Close()
	for {
		messageType, _, err := c.NextReader()
		if err != nil || messageType == goWs.CloseMessage {
			break
		}
	}
	closed <- true
}

//将从watcher中获取的数据，通过websocket发送到浏览器端
func (ws *Websocket) Writer(c *goWs.Conn, closed <-chan bool) {
	ws.watcher.Start()
	defer ws.watcher.Stop()
	defer c.Close()
	for {
		select {
		//等待接收markdown文件原始内容
		case data := <-ws.watcher.C.Raw:
			c.SetWriteDeadline(time.Now().Add(WriteTimeout))
			//将markdown内容转换为html，然后发送给浏览器
			err := c.WriteMessage(goWs.TextMessage, MdConverter.Convert(*data))
			if err != nil {
				return
			}
		case <-closed:
			return
		}
	}
}

func (ws *Websocket) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	//将http协议转换成websocket协议，并开始数据读写操作
	sock, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Can't connect to websocket")
		return
	}
	closed := make(chan bool)

	go ws.Reader(sock, closed)
	ws.Writer(sock, closed)
}
