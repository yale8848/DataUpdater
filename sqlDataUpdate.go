// Create by Yale 2019/6/11 16:37
package dataupdater

import (
	"encoding/json"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type SQLDataUpdate struct {
	DataUpdate
	driverName string
	dataSourceName string
	sql string

	rowsEntity interface{}

}

func NewSQLDataUpdate(driverName ,dataSourceName ,sql string, rowsEntity interface{}) *SQLDataUpdate  {

	return &SQLDataUpdate{driverName:driverName,dataSourceName:dataSourceName,sql:sql,rowsEntity:rowsEntity}

}
func (du *SQLDataUpdate)do(){

	if du.updateFun == nil {
		return
	}

	var (
		err error
		engine *xorm.Engine
		values []byte
	)
	defer func() {
		if err!=nil && du.errFun !=nil{
			du.errFun(err)
		}
	}()
	engine, err = xorm.NewEngine(du.driverName, du.dataSourceName)
	if err!=nil {
		return
	}
	err = engine.SQL(du.sql).Find(du.rowsEntity)
	if err!=nil {
		return
	}
	values,err = json.Marshal(du.rowsEntity)
	if err!=nil {
		return
	}
	if du.canUpdate(values) {
		du.updateFun(du.rowsEntity)
	}
}
func (du *SQLDataUpdate)Start()   {

	du.do()

	go func() {

		defer du.Recover()

		start := false
		tm := time.NewTicker(du.duration)
		for {
			<-tm.C
			if start {
				continue
			}
			start = true
			du.do()
			start = false
		}
	}()
}

