package pet

import (
	"encoding/json"
)

type Controller interface {
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
	return JsonEncoder(result)
}
func JsonEncoder(result map[string]interface{}) (ret []byte, e error) {
	ret, e = json.Marshal(result)
	return
}
