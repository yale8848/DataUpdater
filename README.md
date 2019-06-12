## DataUpdater

Golang remote data update with timer lib. You can update config auto.

## Feature

- http api to load data
- sql api to load data. use [xorm](github.com/go-xorm/xorm) support [Mysql](github.com/go-sql-driver/mysql),[Postgres](github.com/lib/pq),[SQLite3](github.com/mattn/go-sqlite3),[MsSql](github.com/denisenkom/go-mssqldb)

## Demo

 - http
 
 ```go
type ipData struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	Isp      string `json:"isp"`
}

func main()  {
 	data:=&ipData{}
 
 	ch := make(chan bool)
 	ud := NewHttpDataUpdate("get", "http://ip.tianqiapi.com/", nil,data)
 	ud.ErrListener(func(err error) {
       fmt.Println(err)
 	})
 	//if data changed, call back
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

```

- sql

```go

type tt struct {
	Id int64
	Name string
}

func main()  {
 	ch := make(chan bool)
 	to:=make([]*tt,0)
 	sdu:=NewSQLDataUpdate("mysql","xxxxx?charset=utf8","select * from xxxx",
 		&to)
 	sdu.ErrListener(func(err error) {
 		fmt.Println(err)
 	})
 	sdu.SetDuration(time.Second*5)
 	//if data changed, call back
 	sdu.UpdateListener(func(data interface{}) {
 		if v,ok:=data.(*[]*tt);ok {
 			fmt.Printf("%v",v)
 		}
 	})
 	sdu.Start()
 	<-ch
}

```



