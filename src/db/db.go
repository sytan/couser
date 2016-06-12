//Package db implements utility routines for manipulating database.
//Date : 2016-05-08
package db

import (
	. "common"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/fjl/go-couchdb"
	"strings"
)

const (
	USER_GET_BY_USERNAME = "_design/user/_view/get_by_username" //url id for view of invoxia
	PO_GET_BY_STATUS     = "_design/po/_view/get_by_status"     //url id for view of po
	ETHADDR_GET_BY_ID    = "_design/ethaddr/_view/get_by_id"    //url id for view of ethaddr
	RFPI_GET_BY_ID       = "_design/rfpi/_view/get_by_id"       //url id for view of rfpi
	SN_GET_BY_ID         = "_design/sn/_view/get_by_id"         //url id for view of sn

)

const (
	DB_PO          = "po"
	DB_USER        = "user"
	DB_CUSTOMER    = "customer"
	DB_PRODUCT     = "product"
	DB_SOFTWARE    = "software"
	DB_ETHADDR     = "ethaddr"
	DB_RFPI        = "rfpi"
	DB_SN          = "sn"
	DB_PRO_DIP     = "pro_dip"
	DB_PRO_DIP_LOT = "pro_dip_log"
)

var (
	DbMap = make(map[string]*couchdb.DB)
	//DbName = []string{"po", "user", "customer", "software", "ethaddr", "rfpi", "sn", "pro_dip"}
	DbName = []string{DB_PO, DB_USER, DB_CUSTOMER, DB_PRODUCT, DB_SOFTWARE, DB_ETHADDR, DB_RFPI, DB_SN, DB_PRO_DIP, DB_PRO_DIP_LOT}
)

//toStr : convert string like object to string
func toStr(v interface{}) string {
	if m, ok := v.(string); ok {
		return m
	}
	return ""
}

//formatOtps : to make options like purly string to be "string" for couchdb
func formatOtps(otps couchdb.Options) couchdb.Options {
	var object interface{}
	for k, v := range otps {
		if str, ok := v.(string); ok {
			err := json.Unmarshal([]byte(str), &object)
			if err != nil { //while err , it means str is not invalid json struct, it's a purly string
				otps[k] = "\"" + str + "\""
			} else { //while str is string of number "112233445566", it will be unmarshall to float64
				if _, ok = object.(float64); ok {
					otps[k] = "\"" + str + "\""
				}
			}
		}
	}
	return otps
}

func Put(path string, doc interface{}) {

}

func SaveDoc(db string, id string, doc interface{}, rev string) (newrev string, err error) {
	return DbMap[db].Put(id, doc, rev)
}

//similar as GetDoc
func GetDocGabs(path *gabs.Container) *gabs.Container {
	pathStr := toStr(path.Data())
	return GetDoc(pathStr)
}

//simply get a doc
//path is like "dbname/_id"
func GetDoc(path string) *gabs.Container {
	p := strings.SplitN(path, "/", 2)
	db := p[0]
	id := p[1]

	if _, ok := DbMap[db]; !ok {
		return nil
	}
	var doc interface{}

	err := DbMap[db].Get(id, &doc, nil)
	if err != nil {
		fmt.Print("Error while GET :")
		fmt.Println(err)
	}
	docJson, _ := json.Marshal(doc)
	docParsed, _ := gabs.ParseJSON(docJson)

	return docParsed
}

//simply get a view
func GetView(db string, id string, opts couchdb.Options) *gabs.Container {
	var doc interface{}
	opts = formatOtps(opts)
	err := DbMap[db].Get(id, &doc, opts)
	if err != nil {
		fmt.Print("Error while GET :")
		fmt.Println(err)
	}
	docJson, _ := json.Marshal(doc)
	docParsed, _ := gabs.ParseJSON(docJson)

	return docParsed
}

//DB : open the database
func DB(client *couchdb.Client) {
	for _, v := range DbName {
		DbMap[v] = client.DB(v)
	}
}

//GetUser: get user information according to user name
//a user map will be returned.
func GetUser(username string) (user map[string]interface{}) {
	id := USER_GET_BY_USERNAME
	opts := make(couchdb.Options)
	opts["key"] = username

	doc := GetView(DB_USER, id, opts)
	data, _ := doc.ArrayElementP(0, "rows")
	userid := data.Path("value._id")
	password := data.Path("value.password")
	catalog := data.Path("value.catalog")
	previlige := data.Path("value.previlige")
	lastLogin := data.Path("value.last_login")
	user = map[string]interface{}{
		"id":        toStr(userid.Data()),
		"username":  username,
		"password":  toStr(password.Data()),
		"catalog":   toStr(catalog.Data()),
		"previlige": toStr(previlige.Data()),
		"lastLogin": toStr(lastLogin.Data()),
	}
	return user
}

//GetPo: get po information
func GetPo() (poInfor []map[string]interface{}) {
	id := PO_GET_BY_STATUS
	opts := make(couchdb.Options)
	opts["key"] = "undefined"
	doc := GetView(DB_PO, id, opts)
	length, _ := doc.ArrayCountP("rows")

	for j := 0; j < length; j++ {
		data, _ := doc.ArrayElementP(j, "rows")

		po := toStr(data.Path("value.po").Data())
		product := toStr(data.Path("value.product._id").Data())
		temp := strings.Split(product, "/")
		product = temp[len(temp)-1]

		customerDoc := GetDocGabs(data.Path("value.customer._id"))
		softwareDoc := GetDocGabs(data.Path("value.software._id"))

		software := toStr(softwareDoc.Path("version").Data())
		l, _ := customerDoc.ArrayCountP("product")

		var customerId string
		for i := 0; i < l; i++ {
			p, _ := customerDoc.ArrayElementP(i, "product")
			_id := toStr(p.Path("_id").Data())
			alias := toStr(p.Path("alias").Data())
			_id = LastSplit(_id, "/")

			if _id == product {
				customerId = alias
				break
			}
		}
		infor := map[string]interface{}{
			"value":      po,
			"productId":  product,
			"software":   software,
			"customerId": customerId,
		}
		poInfor = append(poInfor, infor)
	}

	return poInfor
}

/*
	rows, err := DB_DIP.Query(fmt.Sprintf(`select value, productId, software, customerId from %s
	 										where status = '%s'`, TB_DIP_PO, "undefined"))
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	var value, productId, software, customerId string
	for rows.Next() {
		err := rows.Scan(&value, &productId, &software, &customerId)
		if err != nil {
			return
		}
		po := map[string]interface{}{
			"value":      value,
			"productId":  productId,
			"software":   software,
			"customerId": customerId,
		}
		poInfor = append(poInfor, po)
	}
	return poInfor
}
*/
/*
//PostLot implements insert test log to database
func PostLog(table string, msg map[string]interface{}) (id, affect int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sql := fmt.Sprintf(`insert into %s(testId,description, productId, umid, testData,spec,unit,
						testTime, testResult, createdTime, createdBy)values(?,?,?,?,?,?,?,?,?,?,?)`, table)
	stmt, err := DB_DIP.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}
	result, err := stmt.Exec(msg["testId"], msg["description"], msg["productId"], msg["umid"], msg["testData"],
		msg["spec"], msg["unit"], msg["testTime"], msg["testResult"], msg["createdTime"], msg["createdBy"])
	if err != nil {
		fmt.Println(err)
	}
	id, _ = result.LastInsertId()
	affect, _ = result.RowsAffected()

	return id, affect
}

//PostPanel implements insert a record to panel
func PostPanel(table string, msg map[string]interface{}) (id, affect int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sql := fmt.Sprintf(`insert into %s(productId, umid,testTime, createdTime,
	 					createdBy, status)values(?,?,?,?,?,?)`, table)
	stmt, err := DB_DIP.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}
	result, err := stmt.Exec(msg["productId"], msg["umid"], msg["testTime"],
		msg["createdTime"], msg["createdBy"], msg["status"])
	if err != nil {
		fmt.Println(err)
	}
	id, _ = result.LastInsertId()
	affect, _ = result.RowsAffected()

	return id, affect
}

//PostPcba implements insert a record to pcba
func PostPcba(table string, msg map[string]interface{}) (id, affect int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sql := fmt.Sprintf(`insert into %s(productId, umid, ethaddr,rfpi,testTime, flag, createdTime,
	 					createdBy, status)values(?,?,?,?,?,?,?,?,?)`, table)
	stmt, err := DB_DIP.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}
	result, err := stmt.Exec(msg["productId"], msg["umid"], msg["ethaddr"], msg["rfpi"], msg["testTime"],
		msg["flag"], msg["createdTime"], msg["createdBy"], msg["status"])
	if err != nil {
		fmt.Println(err)
	}
	id, _ = result.LastInsertId()
	affect, _ = result.RowsAffected()

	return id, affect
}

//GetPo implements get po information
func GetPo() (poInfor []map[string]interface{}) {
	rows, err := DB_DIP.Query(fmt.Sprintf(`select value, productId, software, customerId from %s
	 										where status = '%s'`, TB_DIP_PO, "undefined"))
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	var value, productId, software, customerId string
	for rows.Next() {
		err := rows.Scan(&value, &productId, &software, &customerId)
		if err != nil {
			return
		}
		po := map[string]interface{}{
			"value":      value,
			"productId":  productId,
			"software":   software,
			"customerId": customerId,
		}
		poInfor = append(poInfor, po)
	}
	return poInfor
}

//GetRelativeId implements get the relative id of a record
func GetRelativeId(table, value string) (relativeId int) {
	var minId, id int
	rows, _ := DB_DIP.Query(fmt.Sprintf("select id from %s where value = '%s'", table, value))
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return
		}
	}
	rows, _ = DB_DIP.Query(fmt.Sprintf("select min(id) from %s", table))
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&minId)
		if err != nil {
			return
		}
	}
	relativeId = id - minId
	return relativeId
}

//GetColumn implements get a column with a reference column
func GetColumn(table, column, referColumn, referValue string) (value string) {
	sql := fmt.Sprintf("select %s from %s where %s = '%s'", column, table, referColumn, referValue)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return
		}
	}
	return value
}

//GetColumns implements get a column with two reference column
func GetColumns(table, column, columnOne, valueOne, columnTwo, valueTwo string) (value string) {
	sql := fmt.Sprintf("select %s from %s where %s = '%s' and %s = '%s'", column, table, columnOne, valueOne, columnTwo, valueTwo)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return
		}
	}
	return value
}

//GetColumns implements get a column with three reference column
func GetColumnsThree(table, column, columnOne, valueOne, columnTwo, valueTwo, columnThree, valueThree string) (value string) {
	sql := fmt.Sprintf("select %s from %s where %s = '%s' and %s = '%s' and %s = '%s'", column, table, columnOne, valueOne, columnTwo, valueTwo, columnThree, valueThree)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return
		}
	}
	return value
}

//PostTracking implements insert a record to tracking
func PostTracking(table string, msg map[string]interface{}) (id, affect int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sql := fmt.Sprintf(`insert into %s(po,productId, ethaddr, rfpi, trackingLot,createdTime,
						createdBy)values(?,?,?,?,?,?,?)`, table)
	stmt, _ := DB_DIP.Prepare(sql)
	defer stmt.Close()
	result, _ := stmt.Exec(msg["po"], msg["productId"], msg["ethaddr"], msg["rfpi"],
		msg["trackingLot"], msg["createdTime"], msg["createdBy"])
	id, _ = result.LastInsertId()
	affect, _ = result.RowsAffected()
	return id, affect
}

//GetCount implements get count of records
func GetCount(table, column, value string) (count int) {
	sql := fmt.Sprintf("SELECT count(*) FROM %s where %s = '%s'", table, column, value)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return
		}
	}
	return count
}

//PostPacking implements insert a record to packing
func PostPacking(table string, msg map[string]interface{}) (id, affect int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sql := fmt.Sprintf(`insert into %s(po,productId,trackingLot,packingLot,createdTime,
						createdBy)values(?,?,?,?,?,?)`, table)
	stmt, _ := DB_DIP.Prepare(sql)
	defer stmt.Close()
	result, _ := stmt.Exec(msg["po"], msg["productId"], msg["trackingLot"], msg["packingLot"],
		msg["createdTime"], msg["createdBy"])
	id, _ = result.LastInsertId()
	affect, _ = result.RowsAffected()
	return id, affect
}

//Getlastcolumn implements get last record of table with one column
func GetLastColumn(table, column string) (value string) {
	sql := fmt.Sprintf("select %s from %s order by `id` desc limit 1", column, table)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return
		}
	}
	return value
}

//Getlastcolumns implements get last record of table with two column
func GetLastColumns(table, columnOne, columnTwo string) (valueOne, valueTwo string) {
	sql := fmt.Sprintf("select %s,%s from %s order by `id` desc limit 1", columnOne, columnTwo, table)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&valueOne, &valueTwo)
		if err != nil {
			return
		}
	}
	return valueOne, valueTwo
}

//GetFirstcolumn implements get first record of table with one column
func GetFirstColumn(table, column string) (value string) {
	sql := fmt.Sprintf("select %s from %s order by `id` asc limit 1", column, table)
	rows, _ := DB_DIP.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			return
		}
	}
	return value
}

//UpdateColumn implements update a column of one record
func UpdateColumn(table, column, value, referColumn, referValue string) (affect int64) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	sql := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=?", table, column, referColumn)
	stmt, _ := DB_DIP.Prepare(sql)
	defer stmt.Close()
	result, _ := stmt.Exec(value, referValue)
	affect, _ = result.RowsAffected()
	return affect
}


*/
