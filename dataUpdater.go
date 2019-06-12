// Create by Yale 2019/6/11 14:49
package dataupdater

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

type UpdateFunc func(data interface{})
type ErrFunc func(err error)

type HttpReq struct {
	Url    string
	Data   *url.Values
	Method string
}

type DataUpdater interface {
	UpdateListener(fun UpdateFunc)
	Start()
	SetDuration(dur time.Duration)
	ErrListener(errFun ErrFunc)
}

type DataUpdate struct {
	duration  time.Duration
	updateFun UpdateFunc
	oldSign [md5.Size]byte
	errFun ErrFunc
}
func (du *DataUpdate) SetDuration(dur time.Duration)  {
	du.duration = dur
	if du.duration <= time.Second {
		du.duration = time.Second
	}
}
func (du *DataUpdate)ErrListener(errFun ErrFunc){
	du.errFun =errFun
}
func (du *DataUpdate)Recover(){
	if err := recover(); err != nil &&du.errFun!=nil{
		du.errFun(errors.New(fmt.Sprintf("%v",err)))
	}
}
func (du *DataUpdate) canUpdate(data []byte) bool{

	if data == nil || len(data) == 0 {
		return  false
	}
    copyArray:= func(dist *[md5.Size]byte,src *[md5.Size]byte) {
		for i:=0;i<md5.Size;i++{
			dist[i] = src[i]
		}
	}
	nd:=md5.Sum(data)
	for i:=0;i<md5.Size;i++{
		if du.oldSign[i]!=nd[i] {
			copyArray(&du.oldSign,&nd)
			return true
		}
	}
	return false

}
func (du *DataUpdate) UpdateListener(fun UpdateFunc)  {
	du.updateFun = fun
}
