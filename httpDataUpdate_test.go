// Create by Yale 2019/6/11 16:27
package dataupdater

import (
	"fmt"
	"testing"
	"time"
)

type ipData struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Isp      string `json:"isp"`
}

func TestHttpDataUpdate_Start(t *testing.T) {

	data:=&ipData{}

	ch := make(chan bool)
	ud := NewHttpDataUpdate("get", "http://ip.tianqiapi.com/", nil,data)
	ud.ErrListener(func(err error) {
       fmt.Println(err)
	})
	ud.UpdateListener(func(data interface{}) {

		if v,ok:=data.(*ipData);ok {
			fmt.Printf("%v\r\n",v)
		}
		if v,ok:=data.([]byte);ok {
			fmt.Println(string(v))
		}
	})
	ud.SetDuration(5 * time.Second)
	ud.Start()
	<-ch
}
