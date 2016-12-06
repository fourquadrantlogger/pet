package pet

type DefaultController struct {
}

func (module *DefaultController) Init() error {
	return nil
}
func (module *DefaultController) Decode(bytes []byte) (body map[string]interface{}, e error) {
	return DefaultDecoder(bytes)
}
func (module *DefaultController) Encode(result map[string]interface{}) (res []byte, e error) {
	return DefaultEncoder(result)
}

func (module *DefaultController) CheckLogin(req *HttpRequest) bool {
	return true
}

func (module *DefaultController) Hello(req *HttpRequest, res map[string]interface{}) (e Error) {
	res["result"] = "World!"
	return
}
func (module *DefaultController) ErrorModule(req *HttpRequest, res map[string]interface{}) (e Error) {
	e.Msg = "Invalid Module Name"
	return
}
func (module *DefaultController) ErrorMethod(req *HttpRequest, res map[string]interface{}) (e Error) {
	e.Msg = "Invalid Method Name"
	return
}
