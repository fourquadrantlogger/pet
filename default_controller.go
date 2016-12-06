package pet

type StatusController struct {
}

func (this *StatusController)Status(req *HttpRequest, res map[string]interface{})(e Error){
	res["key"]="value"
	return
}

func (module *StatusController) ErrorController(req *HttpRequest, res map[string]interface{}) (e Error) {
	e.Msg = "无效Controller名"
	e.Code = ERR_INVALID_PARAM
	return
}
func (module *StatusController) ErrorMethod(req *HttpRequest, res map[string]interface{}) (e Error) {
	e.Msg = "无效Method名"
	e.Code = ERR_INVALID_PARAM
	return
}