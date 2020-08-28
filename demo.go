package main

import (
	DB "db-driver/driver"
	"fmt"
)


func main() {
	//连接池1
	db1 := DB.DbConfig{
		DriverName: "mysql",
		Address:    "127.0.0.1",
		User:       "root",
		Password:   "root",
		Port:       "3306",
		DbName:     "test1",
	}
	db1.Connect()

	datas := make(DB.SqlValues,0,10)
	title := DB.SetField{"title", "今天天气不错"}
	uid := DB.SetField{"uid", 123}
	datas = append(datas, title)
	datas = append(datas, uid)

	id,err:=db1.Insert("test", datas)
	if err!=nil {
		fmt.Println("insert err:",err.Error())
	}
	fmt.Println("连接池1: id = ", id)

	//连接池2

	db2 := DB.DbConfig{
		DriverName: "mysql",
		Address:    "127.0.0.1",
		User:       "root",
		Password:   "root",
		Port:       "3306",
		DbName:     "test2",
	}
	db2.Connect()

	datas2 := make(DB.SqlValues,0,10)
	title = DB.SetField{"title", "这里是连接池2"}
	uid = DB.SetField{"uid", 555}
	datas2 = append(datas2, title)
	datas2 = append(datas2, uid)

	id2,err:=db2.Insert("test", datas2)
	if err!=nil {
		fmt.Println("insert err:",err.Error())
	}
	fmt.Println("连接池2: id2 = ", id2)
}
