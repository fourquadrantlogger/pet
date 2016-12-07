package pet

import (
	"net/http"
	"strings"

	"strconv"
)

type HttpRequest struct {
	*http.Request
	JSONBody map[string]interface{}
	BodyRaw  []byte
}

func (this *HttpRequest) IP() string {
	ips := strings.Split(this.Header.Get("X-Forwarded-For"), ",")
	if len(ips[0]) > 3 {
		return ips[0]
	} else {
		addr := this.RemoteAddr[:strings.LastIndex(this.RemoteAddr,":")]

		return addr
	}
}
// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// othewise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
// 		GET /?name=Manu&lastname=
// 		("Manu", true) == c.GetQuery("name")
// 		("", false) == c.GetQuery("id")
// 		("", true) == c.GetQuery("lastname")
func (c *HttpRequest) GetQuery(key string) (string, bool) {
	req := c.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

// Query returns the keyed url query value if it exists,
// othewise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
// 		GET /path?id=1234&name=Manu&value=
// 		c.Query("id") == "1234"
// 		c.Query("name") == "Manu"
// 		c.Query("value") == ""
// 		c.Query("wtf") == ""
func (c *HttpRequest) Query(key string) string {
	value, _ := c.GetQuery(key)
	return value
}


func (c *HttpRequest) GetQueryInt64(key string) int64 {
	i, _ := strconv.ParseInt(c.Query(key), 10, 64)
	return i
}

func (c *HttpRequest) GetQueryUint(key string) uint {
	i, _ := strconv.ParseUint(c.Query(key), 10, 32)
	return uint(i)
}

func (c *HttpRequest) GetQueryString(key string) string {
	return c.Query(key)
}

func (c *HttpRequest) GetOffsetAndLimit() (uint, uint) {
	page := c.GetQueryUint("page")
	limit := c.GetQueryUint("limit")
	if limit == 0 {
		limit = 20
	}
	return page * limit, limit
}
func (c *HttpRequest)CheckJSONParam(params ...string)(error){
	for k,v:=range c.JSONBody{
		var k_exist=false
		for _,p:=range params{
			if(k==p){
				k_exist=true
			}
		}
		if(k_exist==false){
			return NewError(ERR_REQUIRE_PARAM,"json缺少字段"+k,nil)
		}
	}
	return nil
}