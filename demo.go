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

	datas := make(DB.DataStruct)
	datas["title"] = "今天天气不错"
	datas["uid"] = 666

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

	datas2 := make(DB.DataStruct)
	datas2["title"] = "还是不错的"
	datas2["uid"] = 999

	id2,err:=db2.Insert("test", datas)
	if err!=nil {
		fmt.Println("insert err:",err.Error())
	}
	fmt.Println("连接池2: id2 = ", id2)
}
