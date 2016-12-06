package pet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Server struct {
	controllers     map[string]Controller
	conf            *Config
	otherHandleFunc map[string]func(http.ResponseWriter, *http.Request)
}

func NewPet(conf *Config) (server *Server, err error) {

	if err != nil {
		return nil, err
	}
	server = &Server{make(map[string]Controller), conf, make(map[string]func(http.ResponseWriter, *http.Request))}
	server.AddController("status", &StatusController{})
	return server, nil
}

func (server *Server) AddController (name string, controller Controller) (err error) {
	log.Printf("add controller %s ", name)

	if err != nil {
		log.Print("failed")
		return err
	}
	server.controllers[name] = controller
	return
}

func (server *Server) AddHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) (e error) {
	log.Println("add otherHandleFunc pattern %s... ", pattern)
	server.otherHandleFunc[pattern] = handler
	return nil
}

func (server *Server) StartService() error {
	handler := http.NewServeMux()
	handler.HandleFunc("/", server.allHandler)

	for key, value := range server.otherHandleFunc {
		handler.HandleFunc(key, value)
	}

	s := &http.Server{
		Addr:           server.conf.IpPort,
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}
	l := fmt.Sprint("service start at ", server.conf.IpPort, " ...")
	fmt.Println(l)
	return s.ListenAndServe()
}
func (server *Server) allHandler(w http.ResponseWriter, r *http.Request) {
	var result map[string]interface{} = make(map[string]interface{})
	result["status"] = "error"
	var err error

	fields := strings.Split(r.URL.Path[1:], "/")
	var body []byte
	C:=fields[0]
	M:=fields[1]
	M=strings.ToUpper(M[:1])+M[1:]
	body, err = server.handleRequest(C,M,r,result)

	server.processError(w, r, err, body, result)
}
func (server *Server) processError(w http.ResponseWriter, r *http.Request, err error, reqBody []byte, result map[string]interface{}) {
	var re Error
	switch e := err.(type) {
	case nil:
	case Error:
		re = e
	default:
		re = NewError(ERR_INTERNAL, e.Error(), "")
	}

	if re.Code == ERR_NOERR {
		result["status"] = "ok"
		result["time"] = time.Now().Unix()

		server.writeBack(r, w, reqBody, result, true)
	} else {
		server.writeBackErr(r, w, reqBody, re)
	}
}
func (server *Server) writeBack(r *http.Request, w http.ResponseWriter, reqBody []byte, result map[string]interface{}, success bool) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resBytes, err := server.encode(r, result)
	if err != nil {
		result["status"] = "fail"
		result["code"] = ERR_INTERNAL
		result["msg"] = ""
		result["detail"] = err.Error()

		resBytes, _ = DefaultEncoder(result)
	}

	w.Write(resBytes)

	var l string
	var response string
	if reqBody != nil {
		response = string(reqBody)
	}
	format := "uri: %s\n"
	format += "param: %s\n"
	format += "response:\n%s\n"
	format += "\n"

	l = fmt.Sprintf(format, r.URL.String(), response, string(resBytes))
	fmt.Println(l)
}
func (server *Server) writeBackErr(r *http.Request, w http.ResponseWriter, reqBody []byte, err Error) {
	var result map[string]interface{} = make(map[string]interface{})
	result["status"] = "fail"

	fmt.Println("err", err)
	result["code"] = err.Code
	result["msg"] = err.Msg
	result["detail"] = fmt.Sprint(err.Detail)
	server.writeBack(r, w, reqBody, result, true)
}

//根据不同模块encode
func (server *Server) encode(r *http.Request, result map[string]interface{}) (ret []byte, e error) {
	fields := strings.Split(r.URL.Path[1:], "/")
	controllerName := fields[0]
	if len(fields) >= 3 {
		controllerName = fields[1]
	}

	_, ok := server.controllers[controllerName]
	if !ok {
		l := fmt.Sprintf("invalid controller name when write back,%s", controllerName)
		e = errors.New(l)
		return
	}
	return DefaultEncoder(result)
}

//todo::过滤对内部方法（如Init，CheckLogin）的请求
func (server *Server) handleRequest(controllerName string, methodName string, r *http.Request, result map[string]interface{}) ([]byte, error) {
	if methodName == "Init" || methodName == "Decode" || methodName == "Encode" {
		return nil, NewError(ERR_PATH, "非法的Controller方法:"+methodName, "")
	}
	bodyBytes, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, NewError(ERR_INTERNAL, "read http data error : "+e.Error(), "")
	}

	var values []reflect.Value
	controller, ok := server.controllers[controllerName]
	if !ok {
		return nil, NewError(ERR_INVALID_PARAM, "Invalid controller Name", "msg??")
	}
	body, err := DefaultDecoder(bodyBytes)
	if err != nil {
		return nil, NewError(ERR_INVALID_PARAM, err.Error(), "")
	}

	if ok {
		method := reflect.ValueOf(controller).MethodByName(methodName)

		req := &HttpRequest{r, body, bodyBytes}

		if method.IsValid() {
			values = method.Call([]reflect.Value{reflect.ValueOf(req), reflect.ValueOf(result)})
		} else {
			method = reflect.ValueOf(server.controllers["default"]).MethodByName("ErrorMethod")
			values = method.Call([]reflect.Value{reflect.ValueOf(req), reflect.ValueOf(result)})
		}
	} else {
		method := reflect.ValueOf(server.controllers["default"]).MethodByName("ErrorController")
		values = method.Call([]reflect.Value{reflect.ValueOf(&HttpRequest{r, body, bodyBytes}), reflect.ValueOf(result)})
	}
	if len(values) != 1 {
		return bodyBytes, NewError(ERR_INTERNAL, fmt.Sprintf("method %s.%s return value is not 2.", controllerName, methodName), "")
	}
	switch x := values[0].Interface().(type) {
	case nil:
		return bodyBytes, nil
	default:
		return bodyBytes, x.(error)
	}
}
