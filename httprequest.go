package pet

import (
	"net/http"
	"strings"
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
		addr := strings.Split(this.RemoteAddr, ":")
		return addr[0]
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