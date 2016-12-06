package pet

import (
	"encoding/json"
)

type Controller interface {
	Init() (err error)
	Decode(bytes []byte) (body map[string]interface{}, e error)
	Encode(result map[string]interface{}) (ret []byte, e error)
}

func DefaultDecoder(bytes []byte) (map[string]interface{}, error) {
	return JsonDecoder(bytes)
}

func JsonDecoder(bytes []byte) (body map[string]interface{}, e error) {
	if len(bytes) == 0 {
		//可能是Get请求
		body = make(map[string]interface{})
	} else {
		e = json.Unmarshal(bytes, &body)
		if e != nil {
			e = Error{
				Msg: "read body error",
			}
			return
		}
	}
	return
}

func DefaultEncoder(result map[string]interface{}) (ret []byte, e error) {
	resBytes, err := json.Marshal(result)
	if err != nil {
		e = err
		return
	}
	ret = []byte(resBytes)
	return
}
func JsonEncoder(result map[string]interface{}) (ret []byte, e error) {
	ret, e = json.Marshal(result)
	return
}
