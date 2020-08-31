package main

import (
	DB "db-driver/driver"
	"fmt"
	"os"
)

func main() {

	//连接池1
	db1 := DB.DbConfig{
		DriverName: "mysql",
		Addr:       "127.0.0.1",
		User:       "root",
		Passwd:     "root",
		Port:       "3306",
		DBName:     "test1",
	}
	err := db1.Connect()
	if err != nil {
		fmt.Println("connect err:", err.Error())
		os.Exit(1)
	}
	datas := make(DB.DataStruct)
	datas["title"] = "今天天气不错"
	datas["uid"] = 666

	id, err := db1.Insert("test", datas)
	if err != nil {
		fmt.Println("insert err:", err.Error())
	}
	fmt.Println("连接池1: id = ", id)

	//连接池2
	db2 := DB.DbConfig{
		DriverName: "mysql",
		Addr:       "127.0.0.1",
		User:       "root",
		Passwd:     "root",
		Port:       "3306",
		DBName:     "test2",
	}
	err = db2.Connect()
	if err != nil {
		fmt.Println("connect err:", err.Error())
		os.Exit(1)
	}

	datas2 := make(DB.DataStruct)
	datas2["title"] = "还是不错的"
	datas2["uid"] = 999
	//datas2["adddate"] = "2010-01-01"
	datas2.Set("adddate", "2010-01-02")

	id2, err := db2.Insert("test", datas2)
	if err != nil {
		fmt.Println("insert err:", err.Error())
	}
	fmt.Println("连接池2: id2 = ", id2)

	//get one 获取一条

	data, err := db2.GetOne("test", "*", "WHERE id=", id2, "ORDER BY id DESC")
	fmt.Println(data)
	fmt.Println(data["title"])
	fmt.Println(data["uid"])
	fmt.Println(data["adddate"])
	/*t1:=db2.Format2String(data["title"].([]uint8))
	fmt.Printf(t1)
	t2:=db2.Format2String(data["uid"].([]uint8))
	fmt.Printf(t2)*/
	//fmt.Println(reflect.TypeOf(data["uid"]))
	fmt.Println(DB.Format2String(data, "adddate"))

}
