package conv

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
)

//Struct转换绑定
//model 合并的模型
//data 数据
func Struct(model interface{}, data interface{}) {
	result, err := json.Marshal(data)
	if err != nil {
		log.Print("类型不正确" + err.Error())
	}
	json.Unmarshal(result, model)
}

//Struct转换绑定 使用gob编码方式一般用于“相似”的两个结构体传输绑定或者RPC通讯
func StructToBytes(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		log.Println("Struct转Byte类型不正确" + err.Error())
	}
	return buf.Bytes()
}

//bytes转Struct绑定 使用gob编码方式一般用于“相似”的两个结构体传输绑定或者RPC通讯
func BytesToStruct(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
