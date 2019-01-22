package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	connConfMaxReadBufferLength = 0xffff
)

var updrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	configFilename = flag.String("file", "config.json", "config file name")
	configFilepath = flag.String("path", ".", "config file path")
)

type ClientSession struct {
	conn    *websocket.Conn
	tcpConn net.Conn
}

func (session *ClientSession) readWebSocket() {
	defer func() {
		session.conn.Close()
		session.tcpConn.Close()
	}()

	for {
		_, message, err := session.conn.ReadMessage()
		if nil != err {
			log.Println("websocket read error: ", err)
			return
		}
		_, err = session.tcpConn.Write(message)
		if nil != err {
			log.Println("tcpsocket send error: ", err)
			return
		}
	}
}

func (session *ClientSession) readTcpSocket() {
	defer func() {
		session.conn.Close()
		session.tcpConn.Close()
	}()

	readBuffer := make([]byte, connConfMaxReadBufferLength)
	for {
		readLen, err := session.tcpConn.Read(readBuffer)
		if nil != err {
			log.Println("tcpsocket read error: ", err)
			return
		}
		err = session.conn.WriteMessage(websocket.BinaryMessage, readBuffer[:readLen])
		if nil != err {
			log.Println("websocket send error: ", err)
			return
		}
	}
}

func websocketHandle(w http.ResponseWriter, r *http.Request) {

	conn, updraderErr := updrader.Upgrade(w, r, nil)
	if nil != updraderErr {
		log.Print("upgrade:", updraderErr)
		return
	}
	serverAddr := GetPorxyIP() + ":" + strconv.Itoa(GetPorxyPort())

	tcpConn, err := net.DialTimeout("tcp", serverAddr, time.Duration(5)*time.Second)
	if nil != err {
		log.Println("connect server error, disconnect websocket")
		return
	}

	session := &ClientSession{
		conn:    conn,
		tcpConn: tcpConn,
	}

	go session.readWebSocket()
	session.readTcpSocket()

}

func main() {
	flag.Parse()
	initConfig()
	ReadConfig(*configFilepath, *configFilename)

	http.HandleFunc("/websocket", websocketHandle)

	listenAddr := GetListenIP() + ":" + strconv.Itoa(GetListenPort())
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
