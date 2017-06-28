package server

import (
	"net/http"
	"strings"
	"fmt"
	"net/url"
	"github.com/antonholmquist/jason"
	"errors"
	"encoding/json"
)
type CloudRecord struct {
	objectId string `json:"objectId"`
	Id string `json:"id"`
	Addr string `json:"addr"`
}
func createCloudRecord(addr string, id string)(objId string, err error){
	json := fmt.Sprintf(`{"id":"%s","addr":"%s"}`, id, addr)
	resp, err := bmob.Post("user", json)
	if err != nil {
		return
	}
	obj, err := resp.Object()
	if err != nil {
		return "", err
	}
	objId, err = obj.GetString("objectId")
	return
}
func GetCloudRecord(id string)(rec CloudRecord, err error){
	jsonstr := fmt.Sprintf(`{"id":"%s"}`, id)
	obj, err := bmob.GET(fmt.Sprintf(`user?where=%s`, url.QueryEscape(jsonstr)))
	if err != nil {
		panic(err)
	}

	var results []*jason.Object
	if results, err = obj.GetObjectArray("results"); err != nil {
		panic(err)
	}

	if len(results) != 1 {
		err=errors.New(fmt.Sprintf("results count: %d", len(results)))
		return
	}

	err = json.Unmarshal([]byte(results[0].String()), &rec)
	if err != nil {
		panic(err)
	}
	return
}
func SyncNatToCloud(addr string, id string)(err error){
	objId := ""
	json := fmt.Sprintf(`{"id":"%s"}`, id)
	obj, err := bmob.GET(fmt.Sprintf(`user?where=%s`, url.QueryEscape(json)))
	if err != nil {
		panic(err)
	}

	var results []*jason.Object
	if results, err = obj.GetObjectArray("results"); err != nil {
		panic(err)
	}
	if len(results) == 0 {
		objId, err = createCloudRecord(addr, id)
		if err != nil {
			panic(err)
		}
	}else {
		if len(results) != 1 {
			panic(errors.New("results count tao duo"))
		}
		objId, err = results[0].GetString("objectId")
		if err != nil {
			panic(err)
		}
	}

	json = fmt.Sprintf(`{"addr":"%s"}`, addr)
	resp, err := bmob.Put("user/"+objId, json)
	if err != nil {
		panic(err)
	}
	if _, err := resp.GetString("updatedAt"); err != nil {
		return err
	}
	return nil
}
var bmob = Bmob{appid:"1951f0fe596b9a4b90c8c9d80c5072ae", restkey:"5a66cb402f854439fd70f735a7689683"}

const baseUrl string = "https://api.bmob.cn/1/classes/"
type Bmob struct {
	appid string
	restkey string
}

func (b *Bmob) GET(url string) (*jason.Object, error){
	req, err := http.NewRequest("GET", baseUrl + url, nil)
	if err != nil {
		return nil, err
	}
	return b.Do(req)
}
func (b *Bmob) Post(url string, json string) (*jason.Object, error){

	req, err := http.NewRequest("POST", baseUrl + url, strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	return b.Do(req)
}
func (b *Bmob) Put(url string, json string) (*jason.Object, error){
	req, err := http.NewRequest("PUT", baseUrl + url, strings.NewReader(json))
	if err != nil {
		return nil, err
	}
	return b.Do(req)
}

func (b * Bmob) Do(req *http.Request) (*jason.Object, error){
	req.Header.Set("X-Bmob-Application-Id", b.appid)
	req.Header.Set("X-Bmob-REST-API-Key", b.restkey)
	req.Header.Set("Content-Type", "application/json")
	resp, err :=  http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return jason.NewObjectFromReader(resp.Body)
}