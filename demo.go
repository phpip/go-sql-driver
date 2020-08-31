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
		Debug:     true,
	}
	err := db1.Connect()
	if err != nil {
		fmt.Println("connect err:", err.Error())
		os.Exit(1)
	}
	defer db1.Close()

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
		Debug:     true,
	}
	err = db2.Connect()
	if err != nil {
		fmt.Println("connect err:", err.Error())
		os.Exit(1)
	}
	defer db2.Close()

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

	data, err := db2.GetOne("test", "*", "id > ? ORDER BY id DESC", 159)
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

	//select
	data2, err := db2.Select("test", "*", "id > ? ORDER BY id DESC Limit 0,10", 10)
	for i, i2 := range data2 {
		fmt.Println(i, i2)
		fmt.Println(DB.Format2String(i2, "title"))
	}
	//UPDATE
	datas3 := make(DB.DataStruct)
	datas3["title"] = "修改后的结果2"
	datas3["uid"] = 8

	rows, err := db2.Update("test", datas3, "id>?", 86)
	fmt.Println("影响行数:", rows)

	//delete
	deleteid := 86
	rows2, err := db2.Delete("test", "id<?", deleteid)
	fmt.Println("影响行数:", rows2)

	title := "----test---"
	rows3, err := db2.Delete("test", "title=?", title)
	fmt.Println("影响行数:", rows3)

}
