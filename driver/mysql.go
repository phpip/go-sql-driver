package driver

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)




type SetField struct {
	FieldName string
	FieldData interface{}
}
type SqlValues []SetField


type DbConfig struct {
	db *sql.DB
	DriverName string
	Address string
	User string
	Password string
	Port string
	DbName string
}


func (config *DbConfig) Connect() {
	var err error
	config.db, err = sql.Open("mysql", config.User+":"+config.Password+"@tcp("+config.Address+":"+config.Port+")/"+config.DbName)
	if err != nil {
		log.Panicln("err:", err.Error())
	}

	config.db.SetMaxOpenConns(0)
	config.db.SetMaxIdleConns(0)
}

//插入数据
func (config *DbConfig) Insert(table string, datas SqlValues) (id int64, err error) {
	fieldString := ""
	placeString := ""
	var fieldValues []interface{}
	for _, data := range datas {
		if data.FieldName == "" {
			continue
		}
		if fieldString != "" {
			fieldString += ","
			placeString += ","
		}
		fieldString += "`" + data.FieldName + "`"
		placeString += "?"
		fieldValues = append(fieldValues, data.FieldData)
	}
	sql := "INSERT INTO `" + table + "` (" + fieldString + ") VALUES (" + placeString + ")"

	result, err := config.db.Exec(sql, fieldValues...)
	if err != nil {
		return 0, err
	}
	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//更新
func Update(table string, datas SqlValues, where string) {

}
