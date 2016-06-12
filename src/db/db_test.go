package db

import (
	"flag"
	"github.com/fjl/go-couchdb"
	"testing"
)

var (
	username string
	host     string
	port     string
	path     string
)

//connect to couchdb server
func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "couchdb host")
	flag.StringVar(&port, "port", "5984", "couchdb port")
	flag.StringVar(&username, "user", "swissvoice", "user name of the server")
	flag.StringVar(&path, "path", "", "path for GetDoc")
	flag.Parse()

	rawurl := "http://" + host + ":" + port
	client, _ := couchdb.NewClient(rawurl, nil)
	DB(client)
}

//test GetDoc
func TestGetDoc(t *testing.T) {
	t.Log(path)
	doc := GetDoc(path)
	t.Logf("doc : %q", doc)
}

//test GetUser
func TestGetUser(t *testing.T) {
	user := GetUser(username)
	t.Logf("user :%q", user)
}

func TestGetPo(t *testing.T) {
	po := GetPo()
	t.Logf("po : %q", po)
}

//GetUser test
//user := GetUser(*column)
//if len(user) == 0 {
//	t.Errorf("GetUser error, the user is %q ,the len is %d", user, len(user))
//}

//inser testlog
/***
var value = make(map[string]interface{})
value["testId"] = "1111"
value["umid"] = "2000000000"
value["productId"] = "20478987"
value["testData"] = "9A"
value["testTime"] = 14
value["testResult"] = "pass"
value["createdTime"] = "2015-09-07 12:00:12"
value["createdBy"] = 1

id, affect := PostLog(value)
if affect != 1 {
	t.Errorf("Postlog error, the user is %d ,the len is %d", id, affect)
}

***/
/*
	ethaddr := GetEthaddr("1000000000")
	if ethaddr == "" {
		t.Errorf("GetEthaddr error, the ethaddr is %s ", ethaddr)
	} else {
		fmt.Println(ethaddr)
	}

	rfpi := GetRfpi("1000000000")
	if rfpi == "" {
		t.Errorf("GetRFPI error, the rfpi is %s ", rfpi)
	} else {
		fmt.Println(rfpi)
	}

	_, umid := GetUmid("rfpi", "0171A00001")
	if umid == "" {
		t.Errorf("GetUmid error, the umid is %s ", umid)
	} else {
		fmt.Println(umid)
	}

	isOK, _, _ := UpdateUmid("11:22:33:44:55:05", "0171A00005", 1)
	if isOK == false {
		t.Errorf("UpdateUmid error")
	}
*/
/*
	umid := GetUmidByIdentifier(TB_DIP_ETHADDR, "18:B7:9E:02:11:D0")
	if umid == "" {
		t.Errorf("GetUmidByIdentifier error, the umid is %s ", umid)
	} else {
		fmt.Println(umid)
	}

	isRetest := CheckRetest("18:B7:9E:02:11:D0")
	fmt.Println(isRetest)

	isOK, umid, serial := UpdateUmid("18:B7:9E:02:11:D0", "020130B198", 1)
	fmt.Println(isOK, umid, serial)
	fmt.Println(isOK)

}
*/
