package login

import (
	"bufio"
	. "db"
	"flag"
	"fmt"
	"net"
	"testing"
)

var (
	schema string
)

func init() {
	flag.StringVar(&schema, "schema", "swissvoice", "Name of schema to connected.")
	flag.Parse()
}

func TestLogin(t *testing.T) {
	db, err := OpenDB(schema)
	if err != nil {
		return
	}
	defer db.Close()

	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8000")

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}
		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpHandle(t, tcpConn)
	}

}

/*
func tcpHandle(t *testing.T, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	username, password := Login(conn, reader)
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
	}()

	if username == "" {
		t.Errorf("Login error, the username, password is %q and %q ", username, password)
	}
}
*/
func tcpHandle(t *testing.T, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	isLogin, _ := Authorize(conn, reader)
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		conn.Close()
	}()

	if isLogin == true {
		t.Errorf("Login error")
	}
}
