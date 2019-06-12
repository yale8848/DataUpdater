// Create by Yale 2019/6/11 16:16
package dataupdater

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpDataUpdate struct {
	DataUpdate
	req *HttpReq
	dataEntity interface{}
}

func NewHttpDataUpdate(method, url string, data *url.Values,dataEntity interface{}) *HttpDataUpdate {
	return &HttpDataUpdate{req: &HttpReq{Method: method, Url: url, Data: data},dataEntity:dataEntity}
}

func (hdu *HttpDataUpdate) httpReq() (res *http.Response, err error) {
	if hdu.req == nil {
		return nil, errors.New("HttpDataUpdate must set req")
	}
	req := hdu.req
	var reqBody io.Reader
	if req.Data != nil {
		reqBody = strings.NewReader(req.Data.Encode())
	}
	httpReq, err := http.NewRequest(req.Method, req.Url, reqBody)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(httpReq)

}
func (hdu *HttpDataUpdate) do() {

	if hdu.updateFun == nil {
		return
	}
	var (
		res *http.Response
		err error
		bret []byte
	)
	defer func() {
		if err!=nil && hdu.errFun!=nil {
			hdu.errFun(err)
		}
	}()
	res, err = hdu.httpReq()
	if err != nil {
		return
	}
	bret, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	if hdu.canUpdate(bret) {
		if hdu.dataEntity!=nil {
             er := json.Unmarshal(bret,hdu.dataEntity)
			if er==nil {
				hdu.updateFun(hdu.dataEntity)
				return
			}
		}
		hdu.updateFun(bret)
	}

}
func (hdu *HttpDataUpdate) Start() {

	if hdu.req != nil {
		hdu.req.Method = strings.ToUpper(hdu.req.Method)
	}
	hdu.do()
	go func() {
		defer hdu.Recover()

		start := false
		tm := time.NewTicker(hdu.duration)
		for {
			<-tm.C
			if start {
				continue
			}
			start = true
			hdu.do()
			start = false
		}
	}()
}

