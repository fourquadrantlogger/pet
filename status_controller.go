package pet

type StatusController struct {
}

func (this *StatusController)Status(req *HttpRequest, res map[string]interface{})(e error){
	res["key"]="value"
	return nil
}
