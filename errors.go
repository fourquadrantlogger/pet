package pet

import "fmt"

type Error struct {
	Code   uint
	Msg    string //客户端显示的内容
	Detail interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%n, msg=%s, detail=%v", e.Code, e.Msg, e.Detail)
}

func NewError(ecode uint, msg string, detail interface{}) (err Error) {
	err = Error{Code: ecode, Msg: msg, Detail: detail}
	return
}
