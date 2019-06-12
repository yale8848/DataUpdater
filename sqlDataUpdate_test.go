// Create by Yale 2019/6/11 18:00
package dataupdater

import (
	"fmt"
	"testing"
	"time"
)
type tt struct {
	Id int64
	Name string
}
func TestSQLDataUpdate_Start(t *testing.T) {

	ch := make(chan bool)

	to:=make([]*tt,0)
	sdu:=NewSQLDataUpdate("mysql","xxx?charset=utf8","select * from xxx",
		&to)
	sdu.ErrListener(func(err error) {
		fmt.Println(err)
	})
	sdu.SetDuration(time.Second*5)
	sdu.UpdateListener(func(data interface{}) {
		if v,ok:=data.(*[]*tt);ok {
			fmt.Printf("%v",v)
		}
	})
	sdu.Start()
	<-ch

}
