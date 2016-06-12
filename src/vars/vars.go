//Package vars implements all the const , declare and global vars
//Date : 2016-01-08
package vars

const (
	NEWLINE = "\r\n"
	ENDLINE = '\n'
)

/*
	case "log_pcba": //for test log
	{"Cmd":"log_pcba","Msg":{"testId":"4568","description":"version","productId":"20407974","umid":"01123456789","testData":"01.24.89","spec":"PASS","unit":"--","testTime":"90","testResult":"PASS"}}
	case "log_qa": //for qa test log
	{"Cmd":"log_qa","Msg":{"testId":"4568","description":"version","productId":"20407974","umid":"01123456789","testData":"01.24.89","spec":"PASS","unit":"--","testTime":"90","testResult":"PASS"}}
	case "pcba_test": //for pcba test
	{"Cmd":"pcba_test","Msg":{"productId":"20407974","umid":"A1123456789","ethaddr":"18:B7:9E:02:6D:BA","rfpi":"0201324AA8","serial":"2040797455202991","testTime":"90","flag":"55550000","status":"valid"}}
	case "qa_test": //for pcba qa
    {"Cmd":"qa_test","Msg":{"productId":"20407974","umid":"A1123456789","ethaddr":"18:B7:9E:02:6D:BA","rfpi":"0171A1235","testTime":"90","flag":"55550000","status":"valid"}}
	case "get_po": //for get po information
	{"Cmd":"get_po"}
	case "get_umid": //for get umid
	{"Cmd":"get_umid","Msg":{"prefix":"FF"}}
	case "log_umid": //for attach umid to ethaddr, rfpi and serailnumber
	{"Cmd":"log_umid","Msg":{"ethaddr":"18:B7:9E:02:6D:BC","rfpi":"0201324AB8","serial":"2040797455202993","umid":"AD0123456789"}}
	case "verify_barcode_test", "verify_barcode_qa": //for

	{"Cmd":"verify_barcode_pcba","Msg":{"ethaddr":"112233445567","rfpi":"0171A00002","po":"475708951"}}

	{"Cmd":"verify_barcode_function","Msg":{"ethaddr":"18:B7:9E:02:79:21","rfpi":"020132A5E0","po":"475593083"}}
	{"Cmd":"verify_barcode_qa","Msg":{"ethaddr":"18:B7:9E:02:79:21","rfpi":"020132A5E0","po":"475593083"}}
	case "log_tracking":
	{"Cmd":"log_tracking","Msg":{"po":"475593083","productId":"20407974","trackingLot":"00000003T02","ethaddr":"18:B7:9E:02:6D:BA","rfpi":"0201324AA8"}}
	case "log_packing":
	{"Cmd":"log_packing","Msg":{"po":"475593083","productId":"20407974","trackingLot":"00000004T02","packingLot":"0002B02"}}
	case "get_history_packing":
	{"Cmd":"get_history_packing"}

*/

// verify_barcode
//1 {"Cmd":"get_barcode","Msg":{"ethaddr":"112233445566","rfpi":"0171A00001","po":"475708951"}}
//2 {"Cmd":"get_barcode","Msg":{"ethaddr":"112233445567","rfpi":"0171A00002","po":"475708951"}}
//3 {"Cmd":"get_barcode","Msg":{"ethaddr":"112233445568","rfpi":"0171A00003","po":"475708952"}}
//4 {"Cmd":"get_barcode","Msg":{"ethaddr":"112233445569","rfpi":"0171A00004","po":"475708952"}}
//5 {"Cmd":"log_panel","Msg":{"device":"001450092128","doc":{"po":{"_id":"po/475708951"},"product":{"_id":"product/20407974"},"test_log":"hello"}}}
//6	{"Cmd":"log_panel","Msg":{"device":"101450092128","doc":{"po":{"_id":"po/475708951"},"product":{"_id":"product/20407974"},"test_log":"hello"}}}

//7	{"Cmd":"log_panel","Msg":{"device":"101450092128","doc":{"po":{"_id":"po/475708951"},"product":{"_id":"product/20407974"},"test_log":{"panel": {"_id": "pro_dip_log/001450092128","test_spec": {"_id": "test_spec/0001"},"pass_flag": "FFFFFFFFF0"},"pcba_rf": {"_id": "pro_dip_log/001450092128","test_spec": {"_id": "test_spec/0002"},"station": {"_id": "station/0002"}}}}}}

//{"Cmd":"log_panel","Msg":{"device":"071450092128","doc":{"1234":"tang"}}}

//{"cmd":"save_device","msg":{"device":"071450092128"}}

//	log_panel

/*func docUpdate(subDoc interface{}, targetDoc *gabs.Container) *gabs.Container {
	var doc interface{}
	doc = subDoc
	if mmap, ok := doc.(map[string]interface{}); ok {
		for k, v := range mmap {
			targetDoc, _ = targetDoc.Set(v, k)
			fmt.Println(k, v)
		}
	} else if marray, ok := doc.([]interface{}); ok {
		fmt.Println(marray)
	}

	return targetDoc
}


device := toStr(cmdMsg.Msg["device"])
		msgDoc := cmdMsg.Msg["doc"]
		if _, ok := Docs[device]; !ok {
			path := DB_PRO_DIP + "/" + device
			fmt.Println("the path is ", path)
			doc := GetDoc(path)
			fmt.Println("the doc is")
			doc, _ = doc.Set(device, "_id")
			Docs[device] = doc
			fmt.Println("new device : ", Docs[device])
		}
		fmt.Println(Docs[device])
		Docs[device] = docUpdate(msgDoc, Docs[device])
		fmt.Println(Docs[device])
*/
