package pack

import (
	"reflect"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
)

var PackMap = make(map[byte]reflect.Type)

const (
	T_REG  byte = byte(iota)
	T_LINK
	T_PING
	T_PONE
)

func init() {
	PackMap[T_REG] = reflect.TypeOf(PackReg{})
	PackMap[T_LINK] = reflect.TypeOf(PackLink{})
	PackMap[T_PING] = reflect.TypeOf(PackPing{})
	PackMap[T_PONE] = reflect.TypeOf(PackPone{})
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
}

type PackPone struct {
}

// 解包，返回结构
func Parse(data []byte) (ret interface{}, err error) {
	tmp_io := bytes.NewBuffer(data)
	var t byte
	err = binary.Read(tmp_io, binary.BigEndian, t)
	if err != nil {
		return
	}
	dec := gob.NewDecoder(tmp_io)

	packType, ok := PackMap[t]
	if !ok {
		err = errors.New(fmt.Sprintf("type [%v] not exist", t))
		return
	}

	err = dec.DecodeValue(reflect.New(packType))
	if err != nil {
		return
	}

	return
}

func Package(pack interface{})(data []byte, err error){
	tmp_io := bytes.NewBuffer([]byte{})
	for k, v := range PackMap{
		if v == reflect.TypeOf(pack) {
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