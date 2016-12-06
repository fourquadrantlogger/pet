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
