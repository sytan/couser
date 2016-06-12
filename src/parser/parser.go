//Package parser implements message parser .
//date : 2016-01-08
package parser

import (
	. "common"
	. "db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/fjl/go-couchdb"
	. "login"
	"strconv"
	"time"
)

type CmdMsg struct {
	Cmd string
	Msg map[string]interface{}
}

type ReplyMsg struct {
	Cmd string
	Msg interface{}
}

//toStr : convert string like object to string
func toStr(v interface{}) string {
	if m, ok := v.(string); ok {
		return m
	}
	return ""
}

//toStr : convert int like object to int
func toFloat(v interface{}) float64 {
	if m, ok := v.(float64); ok {
		return m
	}
	return -1
}

type MDoc map[string]*gabs.Container

var Docs = make(MDoc)

//MsgParser implements message parser
func MsgParser(message string, currentUser *UserInfor) {
	var cmdMsg CmdMsg
	var replyMsg ReplyMsg
	var byteMsg []byte
	json.Unmarshal([]byte(message), &cmdMsg)
	if cmdMsg.Msg == nil {
		cmdMsg.Msg = make(map[string]interface{})
	}
	cmdMsg.Msg["createdTime"] = time.Now().Format("2006-01-02 15:04:05")
	cmdMsg.Msg["createdBy"] = currentUser.UserId
	replyMsg.Cmd = cmdMsg.Cmd
	if cmdMsg.Cmd != "" {
		fmt.Println(cmdMsg.Cmd)
	}
	switch cmdMsg.Cmd {
	case "get_po": //for get po information
		poInfor := GetPo()
		replyMsg.Msg = poInfor
	case "get_umid": //for get umid
		umid := generateUmid(currentUser, cmdMsg.Msg)
		replyMsg.Msg = map[string]interface{}{
			"umid": umid,
		}
	case "get_barcode": // get barcode information
		infor := GetBarcodeStatus(cmdMsg.Msg)
		replyMsg.Msg = infor

	case "log":
		id := toStr(cmdMsg.Msg["id"])
		db := toStr(cmdMsg.Msg["db"])
		doc := cmdMsg.Msg["doc"]
		rev := log(db, id, doc, Docs)
		replyMsg.Msg = map[string]interface{}{
			"rev": rev,
		}
		fmt.Println(Docs[db+"/"+id])

	case "save": //commit data to couchdb
		id := toStr(cmdMsg.Msg["id"])
		db := toStr(cmdMsg.Msg["db"])
		rev, err := commit(db, id, Docs)
		fmt.Println(err) //log here
		replyMsg.Msg = map[string]interface{}{
			"rev": rev,
		}

	/*

		d	evice := toStr(cmdMsg.Msg["device"])
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

	//
	case "log_pcba": //log pcba data
		//
	case "log_function": //log function data
		//
	case "log_qa": //log qa data
		//
	case "log_packing": //log packing data
		//
	case "get_packing": //get packing history
		//
	default:
		replyMsg.Msg = "Unkown message"
	}
	byteMsg, _ = json.Marshal(replyMsg)
	currentUser.Conn.Write(byteMsg)
}

//GetUmid implements generate umid
func generateUmid(currentUser *UserInfor, msg map[string]interface{}) (umid string) {
	prefix, _ := msg["prefix"].(string)
	userId := currentUser.UserId
	if prefix == "" {
		prefix = SubStrByLen(userId, -2, "0")
	}
	umid = prefix + strconv.FormatInt(time.Now().Unix(), 10)
	return umid
}

func GetBarcodeStatus(msg map[string]interface{}) (infor map[string]interface{}) {
	var sn, device string
	var status, passFlag, packing string
	var opts couchdb.Options

	ethaddr, _ := msg["ethaddr"].(string)
	rfpi, _ := msg["rfpi"].(string)
	po, _ := msg["po"].(string)
	opts = make(couchdb.Options)
	opts["key"] = ethaddr
	ethaddrDoc := GetView(DB_ETHADDR, ETHADDR_GET_BY_ID, opts)
	opts = make(couchdb.Options)
	opts["key"] = rfpi
	rfpiDoc := GetView(DB_RFPI, RFPI_GET_BY_ID, opts)

	ethaddrOffset := toFloat(ethaddrDoc.Path("offset").Data())
	rfpiOffset := toFloat(rfpiDoc.Path("offset").Data())

	opts = make(couchdb.Options)
	opts["skip"] = ethaddrOffset
	opts["limit"] = 1
	snDoc := GetView(DB_SN, SN_GET_BY_ID, opts)
	snOffset := toFloat(snDoc.Path("offset").Data())
	snData, _ := snDoc.ArrayElementP(0, "rows")
	sn = toStr(snData.Path("id").Data())

	ethaddrData, _ := ethaddrDoc.ArrayElementP(0, "rows")
	rfpiData, _ := rfpiDoc.ArrayElementP(0, "rows")

	ethaddrPo := toStr(ethaddrData.Path("value.po._id").Data()) //get po : "po/475708951"
	rfpiPo := toStr(rfpiData.Path("value.po._id").Data())
	snPo := toStr(snData.Path("value.po._id").Data())
	ethaddrPo = LastSplit(ethaddrPo, "/")
	rfpiPo = LastSplit(rfpiPo, "/")
	snPo = LastSplit(snPo, "/")

	ethaddrDevice := toStr(ethaddrData.Path("value.device._id").Data()) //get device :"pro_dip/001450092128"
	rfpiDevice := toStr(rfpiData.Path("value.device._id").Data())
	snDevice := toStr(snData.Path("value.device._id").Data())
	ethaddrDevice = LastSplit(ethaddrDevice, "/")
	rfpiDevice = LastSplit(rfpiDevice, "/")
	snDevice = LastSplit(snDevice, "/")

	ethaddrStatus := toStr(ethaddrData.Path("value.status").Data()) //get status
	rfpiStatus := toStr(rfpiData.Path("value.status").Data())
	snStatus := toStr(snData.Path("value.status").Data())

	if ethaddrPo != po || rfpiPo != po || ethaddrPo == "" { //check whether the barcode is belong to current po
		status = "invalid_po"
		goto RETURN
	}

	if ethaddrOffset != rfpiOffset || snOffset != rfpiOffset {
		status = "invalid_match" //ethaddr and rfpi doesn't match
		goto RETURN
	}

	if ethaddrStatus != rfpiStatus || snStatus != rfpiStatus {
		status = "invalid_status"
		goto RETURN
	} else {
		status = ethaddrStatus
	}

	if ethaddrDevice == rfpiDevice && snDevice == rfpiDevice {
		device = ethaddrDevice
		if device != "" {
			devicePath := ethaddrData.Path("value.device._id")
			deviceDoc := GetDocGabs(devicePath)
			passFlag = toStr(deviceDoc.Path("pass_flag").Data())

			packing = toStr(deviceDoc.Path("packing.reel").Data())
		}
	} else {
		status = "invalid_devie"
		goto RETURN
	}

RETURN:
	infor = map[string]interface{}{
		"status":   status,
		"serial":   sn,
		"passFlag": passFlag,
		"packing":  packing,
		"device":   device,
	}
	return infor
}

func log(db, id string, doc interface{}, Docs MDoc) (rev string) { //add log to doc
	path := db + "/" + id
	if _, ok := Docs[path]; !ok {
		doc := GetDoc(path)
		doc.SetP(id, "_id")
		Docs[path] = doc
	}
	rev = toStr(Docs[path].Path("_rev").Data())
	docMerge(doc, Docs[path])

	return rev
}

func commit(db, id string, Docs MDoc) (rev string, err error) { //save doc to couchdb
	path := db + "/" + id
	if _, ok := Docs[path]; !ok {
		return rev, errors.New("unkown doc id")
	}
	rev = toStr(Docs[path].Path("_rev").Data())
	newRev, err := SaveDoc(db, id, Docs[path].Data(), rev)
	if err == nil {
		Docs[path].SetP(newRev, "_rev") //update rev after save doc
	}

	return newRev, err
}

func docMerge(subDoc interface{}, targetDoc *gabs.Container) *gabs.Container {
	if mmap, ok := subDoc.(map[string]interface{}); ok { //a couchdb doc is always a map first
		for k, v := range mmap {
			sub := targetDoc.Search(k) //check whether key exist
			if sub.Data() != nil {     //when key exit , revise the key value
				if _, ok := v.(map[string]interface{}); ok { //when sub json is a map , recursion
					docMerge(v, sub)
				} else if arr, ok := v.([]interface{}); ok { //when sub json is a list, append the element
					for _, a := range arr {
						targetDoc.ArrayAppendP(a, k)
					}
				} else { //whe sub json is not an object , set the value
					targetDoc.SetP(v, k)
				}
			} else { //when key is new , add the key and set the value
				targetDoc.SetP(v, k)
			}
		}
	}
	return targetDoc
}

/*

	var isPo, status, serial, idStr, serialPo string
	ethaddr, _ := msg["ethaddr"].(string)
	rfpi, _ := msg["rfpi"].(string)
	po, _ := msg["po"].(string)


	ethaddrDoc := GeView()


	ethaddrOffset := GetOffset()
	rfpiId := GetOffset()

	ethaddrId := GetRelativeId(TB_DIP_ETHADDR, ethaddr)
	rfpiId := GetRelativeId(TB_DIP_RFPI, rfpi)
	ethaddrPo := GetColumn(TB_DIP_ETHADDR, "po", "value", ethaddr)
	rfpiPo := GetColumn(TB_DIP_RFPI, "po", "value", rfpi)


	if ethaddrId == rfpiId && ethaddrId >= 0 {
		firstIdStr := GetFirstColumn(TB_DIP_SERIAL, "id")
		firstId, _ := strconv.Atoi(firstIdStr)
		serialId := firstId + ethaddrId
		serialIdStr := strconv.Itoa(serialId)
		serial = GetColumn(TB_DIP_SERIAL, "value", "id", serialIdStr)
		serialPo = GetColumn(TB_DIP_SERIAL, "po", "value", serial)
	}

	if ethaddrPo != po || rfpiPo != po || serialPo != po {
		isPo = "false"
	} else {
		isPo = "true"
	}

	switch cmd {
	case "verify_barcode_qa": //debug
		ethaddrUmid := GetColumn(TB_DIP_ETHADDR, "umid", "value", ethaddr)
		rfpiUmid := GetColumn(TB_DIP_RFPI, "umid", "value", rfpi)
		serialUmid := GetColumn(TB_DIP_SERIAL, "umid", "value", serial)
		if ethaddrUmid == rfpiUmid && serialUmid == rfpiUmid {
			idStr = GetColumns(TB_DIP_DEVICE_PCBA, "id", "umid", ethaddrUmid, "status", "valid") //debug
			if idStr == "" {
				idStr = "-1"
			}
			id, _ := strconv.Atoi(idStr)
			if id >= 0 {
				status = "valid"
			} else {
				status = "untest"
			}
		} else {
			status = "invalid"
		}
		infor = map[string]interface{}{
			"status": status,
			"serial": serial,
			"umid":   ethaddrUmid,
			"isPo":   isPo,
		}
	case "verify_barcode_function":
		ethaddrUmid := GetColumn(TB_DIP_ETHADDR, "umid", "value", ethaddr)
		rfpiUmid := GetColumn(TB_DIP_RFPI, "umid", "value", rfpi)
		serialUmid := GetColumn(TB_DIP_SERIAL, "umid", "value", serial)
		if ethaddrUmid == rfpiUmid && serialUmid == rfpiUmid {
			idStr = GetColumns(TB_DIP_DEVICE_PCBA, "id", "umid", ethaddrUmid, "status", "valid") //debug
			if idStr == "" {
				idStr = "-1"
			}
			id, _ := strconv.Atoi(idStr)
			if id >= 0 {
				status = "valid"
			} else {
				status = "untest"
			}
		} else {
			status = "invalid"
		}
		infor = map[string]interface{}{
			"status": status,
			"serial": serial,
			"umid":   ethaddrUmid,
			"isPo":   isPo,
		}
	case "verify_barcode_pcba":
		if ethaddrId == rfpiId && ethaddrId >= 0 {
			idStr = GetColumnsThree(TB_DIP_DEVICE_PCBA, "id", "ethaddr", ethaddr, "rfpi", rfpi, "status", "valid")
			if idStr == "" {
				idStr = "-1"
			}
			id, _ := strconv.Atoi(idStr)
			if id >= 0 {
				status = "retest"
			} else {
				status = "valid"
			}
		} else {
			status = "invalid"
		}
		infor = map[string]interface{}{
			"status": status,
			"serial": serial,
			"isPo":   isPo,
		}
	default:
		status = ""
		infor = map[string]interface{}{
			"status": status,
			"serial": serial,
			"isPo":   isPo,
		}
	}

	return infor
}

/*
func logUmid(msg map[string]interface{}) (infor map[string]interface{}) {
	var status string
	umid, _ := msg["umid"].(string)
	ethaddr, _ := msg["ethaddr"].(string)
	rfpi, _ := msg["rfpi"].(string)
	serial, _ := msg["serial"].(string)
	affectEthaddr := UpdateColumn(TB_DIP_ETHADDR, "umid", umid, "value", ethaddr)
	affectRfpi := UpdateColumn(TB_DIP_RFPI, "umid", umid, "value", rfpi)
	affectSerial := UpdateColumn(TB_DIP_SERIAL, "umid", umid, "value", serial)
	ethaddrUmid := GetColumn(TB_DIP_ETHADDR, "umid", "value", ethaddr)
	rfpiUmid := GetColumn(TB_DIP_RFPI, "umid", "value", rfpi)
	serialUmid := GetColumn(TB_DIP_SERIAL, "umid", "value", serial)
	if umid != "" && ethaddrUmid == umid && rfpiUmid == umid && serialUmid == umid {
		status = "valid"
		if affectEthaddr != 1 || affectRfpi != 1 || affectSerial != 1 {
			status = "repeat"
		}
	} else {
		status = "invalid"
	}
	infor = map[string]interface{}{
		"status":  status,
		"ethaddr": ethaddrUmid,
		"rfpi":    rfpiUmid,
		"serial":  serialUmid,
	}
	return infor
}

func logTracking(msg map[string]interface{}) (infor map[string]interface{}) {
	ethaddr, _ := msg["ethaddr"].(string)
	rfpi, _ := msg["rfpi"].(string)
	trackingLot, _ := msg["trackingLot"].(string)
	var id, affect int64
	statusInfor := GetBarcodeStatus("verify_barcode_qa", msg)
	status := statusInfor["status"]
	isPo := statusInfor["isPo"]
	if status == "valid" && isPo == "true" {
		idStr := GetColumns(TB_DIP_PACKING_T, "id", "ethaddr", ethaddr, "rfpi", rfpi)
		idTemp, _ := strconv.Atoi(idStr)
		if idTemp > 0 {
			status = "repeat"
		} else {
			id, affect = PostTracking(TB_DIP_PACKING_T, msg)
		}
	}
	count := GetCount(TB_DIP_PACKING_T, "trackingLot", trackingLot)

	infor = map[string]interface{}{
		"status": status,
		"id":     id,
		"affect": affect,
		"count":  count,
		"isPo":   isPo,
	}
	return infor
}

func logPacking(msg map[string]interface{}) (infor map[string]interface{}) {
	trackingLot, _ := msg["trackingLot"].(string)
	packingLot, _ := msg["packingLot"].(string)
	status := GetColumn(TB_DIP_PACKING_T, "status", "trackingLot", trackingLot)
	var id, affect int64
	if status == "valid" {
		idStr := GetColumn(TB_DIP_PACKING_B, "id", "trackingLot", trackingLot)
		idTemp, _ := strconv.Atoi(idStr)
		if idTemp > 0 {
			status = "repeat"
		} else {
			id, affect = PostPacking(TB_DIP_PACKING_B, msg)
		}
	} else {
		status = "invalid"
	}
	count := GetCount(TB_DIP_PACKING_B, "packingLot", packingLot)
	infor = map[string]interface{}{
		"status": status,
		"id":     id,
		"affect": affect,
		"count":  count,
	}
	return infor

}

func getPackingHistory() (infor map[string]interface{}) {
	trackingLot_t := GetLastColumn(TB_DIP_PACKING_T, "trackingLot")
	trackingLot_b, packingLot := GetLastColumns(TB_DIP_PACKING_B, "trackingLot", "packingLot")
	var status, trackingLot string
	if trackingLot_t != trackingLot_b {
		status = "oldlot"
		trackingLot = trackingLot_t
	} else {
		status = "newlot"
		lotStr := trackingLot_t[0:8] //debug
		lot, _ := strconv.Atoi(lotStr)
		newLot := lot + 1
		newLotStr := strconv.Itoa(newLot)
		new_trackingLot_t := SubStrByLen(newLotStr, 8, "0") + trackingLot_t[8:len(trackingLot_t)]
		trackingLot = new_trackingLot_t
	}

	infor = map[string]interface{}{
		"status":      status,
		"trackingLot": trackingLot,
		"packingLot":  packingLot,
		//"shippingLot": shippingLot,
	}
	return infor
}

*/
