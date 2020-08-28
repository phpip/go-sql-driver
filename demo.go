package main

import (
	"database/sql"
	DB "db-driver/driver"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func init() {
	var err error
	DB.Db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/ops")
	if err != nil {
		log.Panicln("err:", err.Error())
	}

	DB.Db.SetMaxOpenConns(0)
	DB.Db.SetMaxIdleConns(0)
}


func main() {
	datas := make(DB.SqlValues, 10)
	title := DB.SetField{"title", "今天天气不错"}
	uid := DB.SetField{"uid", 123}
	datas = append(datas, title)
	datas = append(datas, uid)

	DB.Insert("test", datas)
}
