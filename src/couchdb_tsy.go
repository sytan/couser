package main

import (
	"fmt"
	//"log"
	//. "common"
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/fjl/go-couchdb"
)

//formatOtps : to make options like purly string to be "string" for couchdb
func formatOtps(otps couchdb.Options) couchdb.Options {
	var object interface{}
	for k, v := range otps {
		if str, ok := v.(string); ok {
			err := json.Unmarshal([]byte(str), &object)
			if err != nil { //while err , it means str is not invalid json struct, it's a purly string
				otps[k] = "\"" + str + "\""
			}
		} else {
			fmt.Println(otps)
		}
	}
	return otps
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

func main() {

	type MDoc map[string]*gabs.Container

	var Docs = make(MDoc)
	subJson := []byte(`{
			"name":"tang",
			"age":"100",
			"job":{
				"level":"B",
				"ip":"127.0.0.1",
				"manager":{
					"name":"xiao",
					"age":"50"
				},
				"member":["1","2"],
				"target":[{"fruit":"apple"},{"phone":"huawei"}]
			},
			"school":"zn"

		}`)

	targetJson := []byte(`{
			"name":"tiny",
			"age":"18",
			"sex":"male",
			"job":{
				"title":"engineer",
				"level":"A",
				"manager":{
					"_id":"0001",
					"name":"tt"
				},
				"member":["3","4"],
				"target":[{"fruit":"oringe"},{"phone":"huawei"}]
			}
		}`)
	subDoc, _ := gabs.ParseJSON(subJson)
	targetDoc, _ := gabs.ParseJSON(targetJson)
	Docs["hi"] = targetDoc
	fmt.Println(Docs["hi"])
	docMerge(subDoc.Data(), Docs["hi"])
	fmt.Println(Docs["hi"])

	/*
		client, _ := couchdb.NewClient("http://127.0.0.1:5984", nil)
		db := client.DB("songs")
		opt := make(couchdb.Options)
	*/
	/*
		var result interface{}
		var opt couchdb.Options
		opt = make(map[string]interface{})
		opt["include_docs"] = "false"
		err := db.AllDocs(&result, opt)
		if err != nil {
			//error handler
		}
		re, _ := json.Marshal(result)
		jsonParsed, _ := gabs.ParseJSON(re)
		fmt.Println(jsonParsed.String() + "\n")
	*/ /*
		id := "_design/songs/_view/all"
		var doc interface{}

		//j, _ := json.Marshal("hunter")
		//fmt.Println(string(j))
		//var x interface{}
		//json.Unmarshal([]byte(`{    "tik tok": "dd"     }`), &x)
		//fmt.Println(x)
		opt["key"] = "dying in the song"
		opt = formatOtps(opt)
		err := db.Get(id, &doc, opt)
		if err != nil {
			//error handler
		}
		docJson, _ := json.Marshal(doc)
		fmt.Println(string(docJson))
		/*
			docJson, _ := json.Marshal(doc)
			docParsed, _ := gabs.ParseJSON(docJson)
			rev, _ := ToString(docParsed.Path("_rev").Data())
			docParsed.SetP("123456", "ti.ty.hu")
			myarrParsed := docParsed.Path("fruit")
			myarrParsed.SetIndex(rev, 0)
			docParsed.DeleteP("ti")

			//docParsed.S("fruit").SetIndex("app", 0)
			//docParsed.ArrayAppend("apple", "fruit")
			//docParsed.ArrayAppend("oringe", "fruit")
			//son.Unmarshal([]byte(docParsed.String()), &doc)
			newRev, err := db.Put(id, docParsed.Data(), rev)

			fmt.Println(newRev, err)
	*/
	/*
		v := Values{}
		v.Add("keys", []string{"sy", "sy2"})
		msg, _ := json.Marshal(v)

		b := bytes.NewReader(json)
		client := &http.Client{}
		method := "PUT"
		req, _ := http.NewRequest(method, "http://localhost:5984/ck/tsy10", body)
		//req.Header.Set("Accept", "text/plain")             //Content-Type must be application/json for couchdb
		req.Header.Set("Content-Type", "application/json") //Content-Type must be application/json for couchdb
		resp, _ := client.Do(req)                          //发送
		defer resp.Body.Close()                            //一定要关闭resp.Body
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
	*/
}
