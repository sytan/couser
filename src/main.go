//Package main implements tcp/ip handle
//Date : 2016-05-08
package main

import (
	"bufio"
	. "db"
	"flag"
	"fmt"
	"github.com/fjl/go-couchdb"
	. "login"
	"net"
	. "parser"
	"strings"
	. "vars"
)

var (
	host      string
	port      string
	rawurl    string
	broadcast bool
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "couchdb host")
	flag.StringVar(&port, "port", "5984", "couchdb port")
	flag.BoolVar(&broadcast, "b", false, "broadcast message to other clients")
	flag.Parse()
	rawurl = "http://" + host + ":" + port
}

func main() {
	client, _ := couchdb.NewClient(rawurl, nil) //As fawurl has a format of "http://xxx:yyy" , there will be no err from NewClient
	DB(client)

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", ":8000")

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpHandle(tcpConn)
	}

}

func tcpHandle(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	isLogin, username := Authorize(conn, reader)
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
		delete(UsersOnline, username)
		fmt.Println(len(UsersOnline), " users remaining")
	}()
	for isLogin == true {
		message, err := reader.ReadString(ENDLINE)
		if err != nil {
			return
		}
		message = strings.Trim(message, NEWLINE)
		if (message) != "" {
			MsgParser(message, UsersOnline[username])
		}
		if broadcast == true && (message) != "" { //broadcast message to all other clients
			msg := username + ":" + message + "\n"
			broadcastMessage(msg, conn)
		}

	}
}

func broadcastMessage(message string, conn *net.TCPConn) {
	b := []byte(message)
	for _, user := range UsersOnline { //broadcast message to all other clients
		if user.Conn != conn {
			user.Conn.Write(b)
		}
	}
}
