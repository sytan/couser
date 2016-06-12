//Pakage login implements server login
//Date : 2015-12-26
package login

import (
	"bufio"
	. "db"
	"net"
	"strings"
	. "vars"
)

type UserInfor struct {
	UserId    string
	Username  string
	Catalog   string
	Previlige string
	LastLogin string
	Conn      *net.TCPConn
}

var UsersOnline map[string]*UserInfor = make(map[string]*UserInfor)

var msgLogin string = "Please login: "
var msgPassword string = "Password: "
var MsgLoginFailed string = "Fail to login , please check your user name or password!"
var MsgLoginSucceed string = "Welcome to swissvoice!"
var MsgLoginOverTimes string = "Sorry, Game over!"
var MsgAlreadyOnline = "Already online !"

//var enableEcho = []byte("\xff\xfc\x01")  //escape sequence -> IAC:\xff enable: \xfb echo :\x01  //debug
//var disableEcho = []byte("\xff\xfb\x01") //escape sequence -> IAC:\xff disable: \xfb echo :\x01 //debug
//debug , please refer to https://github.com/Cristofori/kmud/tree/master/telnet

func Login(conn *net.TCPConn, reader *bufio.Reader) (username, password string) {
	msgPrompt := msgLogin
	conn.Write([]byte(msgPrompt))
	/*
		if msgPrompt != msgPassword {
			conn.Write(enableEcho)
		}
	*/ //debug
	for {
		message, err := reader.ReadString(ENDLINE)
		if err != nil {
			return
		}
		if message != NEWLINE {
			if msgPrompt == msgLogin {
				username = strings.Trim(message, NEWLINE)
				msgPrompt = msgPassword
				conn.Write([]byte(msgPrompt))
				//if msgPrompt == msgPassword { //debug
				//	conn.Write(disableEcho)
				//}
			} else {
				password = strings.Trim(message, NEWLINE)
				break
			}
		} else {
			conn.Write([]byte(msgPrompt))
		}
	}
	return username, password
}

func Authorize(conn *net.TCPConn, reader *bufio.Reader) (isLogin bool, username string) {
	tryTimes := 3
	isLogin = false
	var password string

LOGIN: //loop for login
	for {
		username, password = Login(conn, reader)
		tryTimes -= 1
		user := GetUser(username)
		if password == user["password"] {
			_, ok := UsersOnline[username]
			if !ok {
				UsersOnline[username] = new(UserInfor)
				UsersOnline[username].Conn = conn
				UsersOnline[username].UserId, _ = user["id"].(string)
				UsersOnline[username].Catalog, _ = user["catalog"].(string)
				UsersOnline[username].Username, _ = user["username"].(string)
				UsersOnline[username].Previlige, _ = user["previlige"].(string)
				UsersOnline[username].LastLogin, _ = user["lastLogin"].(string)
				isLogin = true
				conn.Write([]byte(MsgLoginSucceed + NEWLINE))
				break LOGIN
			} else {
				conn.Write([]byte(MsgAlreadyOnline + NEWLINE))
				continue
			}
		}
		if tryTimes == 0 {
			conn.Write([]byte(MsgLoginOverTimes + NEWLINE))
			break LOGIN
		} else {
			conn.Write([]byte(MsgLoginFailed + NEWLINE))
		}
	}

	return isLogin, username
}
