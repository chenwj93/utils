package utils

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"runtime/debug"
	"strings"
)

/**
create by cwj on j2018-05-30
used by router
 */
type RouterInterface interface {
	Handle(operation string, paramInput map[string]interface{}) (*Response, error)
}

type Handler struct {
	RootPath     string
	GetDataFuncs map[string]func(r *http.Request) interface{}
	F            func() RouterInterface
}

type Response struct {
	Code int
	Json []byte
}

func NewResponse() *Response {
	return &Response{}
}

func ErrHandle() {
	if err := recover(); err != nil {
		//v := fmt.Sprintf("ERROR!!\n%s--\n  stack \n%s", err,string(debug.Stack()))

		ULog.Println("ERROR :", err)
		ULog.Println(string(debug.Stack()))
	}
}

func (c *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ULog.Println(r.URL.Path)
	defer ErrHandle()
	r.ParseForm()
	var e error
	var paramInput = make(map[string]interface{})
	if r.Method == "POST" {
		body, e := ioutil.ReadAll(r.Body)
		if e == nil {
			e = json.Unmarshal(body, &paramInput)
			if e != nil {
				ULog.Println("post data format error : %s\n", e)
			}
		} else {
			ULog.Println("read post data error :", e)
		}
		setRawQuery(paramInput, r.URL.RawQuery)
	} else if r.Method == "GET" {
		for k, v := range r.Form {
			paramInput[k] = v[0]
		}
	} else if r.Method == "OPTIONS" {
		setAccessControl(w)
		w.WriteHeader(200)
		paramOutput := map[string]interface{}{"status": 200}
		res, _ := json.Marshal(paramOutput)
		w.Write(res)
		return
	}
	c.getDatas(paramInput, r)
	operation := strings.TrimPrefix(r.URL.Path, c.RootPath)
	var res *Response
	res, e = c.F().Handle(operation, paramInput)
	setAccessControl(w)
	if e != nil {
		ULog.Println(e)
		w.WriteHeader(500)
	} else {
		w.WriteHeader(res.Code)
		w.Write(res.Json)
	}
}

func (c *Handler) getDatas(paramInput map[string]interface{}, r *http.Request) {
	for k, f := range c.GetDataFuncs {
		paramInput[k] = f(r)
	}
}

func setAccessControl(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")               //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "*")   //header的类型
	w.Header().Set("content-type", "application/json;charset=utf-8") //返回数据格式是json
}

func setRawQuery(m map[string]interface{}, rawQuery string) {
	rawList := strings.Split(rawQuery, "&")
	for _, raw := range rawList {
		kAndV := strings.Split(raw, "=")
		if len(kAndV) == 2 {
			m[kAndV[0]] = kAndV[1]
		}
	}
}
