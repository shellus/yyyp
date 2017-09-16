package pack

import (
	"reflect"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
)

var PackMap = make(map[byte]reflect.Type)

const (
	T_REG  byte = byte(iota)
	T_LINK
	T_PING
	T_PONE
	T_CONNECT
)

func init() {
	PackMap[T_REG] = reflect.TypeOf(PackReg{})
	PackMap[T_LINK] = reflect.TypeOf(PackLink{})
	PackMap[T_PING] = reflect.TypeOf(PackPing{})
	PackMap[T_PONE] = reflect.TypeOf(PackPone{})
	PackMap[T_CONNECT] = reflect.TypeOf(PackConnect{})

	gob.Register(PackReg{})
	gob.Register(PackLink{})
	gob.Register(PackPing{})
	gob.Register(PackPone{})
	gob.Register(PackConnect{})
	//for _,v := range PackMap {
	//	gob.Register(reflect.New(v))
	//}
}

type PackReg struct {
	Name       string
}

type PackLink struct {
	Name       string
}

// 命令一个node去连接一个地址
// 或是答应一个node的link请求
type PackConnect struct {
	RemoteAddr       string
}

type PackPing struct {
	Expansion []byte
}

type PackPone struct {
	Expansion []byte
}

// 对方通知过来的错误
type PackErr struct {
	Message string
}

// 解包，返回结构
func Parse(data []byte) (ret interface{}, err error) {
	tmp_io := bytes.NewBuffer(data)
	var k byte
	err = binary.Read(tmp_io, binary.BigEndian, &k)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(tmp_io)

	bodyType, ok := PackMap[k]
	if !ok {
		err = errors.New(fmt.Sprintf("type [%v] not exist", k))
		return
	}



	retReflect := reflect.New(bodyType).Elem()
	log.Printf("recv pack type [%s][%#v]", retReflect.String(), k)
	err = dec.DecodeValue(retReflect)
	if err != nil {
		return
	}
	ret = retReflect.Addr().Interface()

	return
}

func Package(pack interface{})(data []byte, err error){
	tmp_io := bytes.NewBuffer([]byte{})
	for k, v := range PackMap{
		if v == reflect.TypeOf(pack) {
			log.Printf("send pack type [%s][%#v]", reflect.New(v).Elem().String(), k)
			binary.Write(tmp_io, binary.BigEndian, k)
			enc := gob.NewEncoder(tmp_io)
			err = enc.Encode(pack)
			if err != nil {
				return
			}
			data = tmp_io.Bytes()
			return
		}
	}
	err = errors.New(fmt.Sprintf("pack pack [%v] invalid", pack))
	return
}