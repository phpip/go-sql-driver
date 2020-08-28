package driver

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type SetField struct {
	FieldName string
	FieldData interface{}
}
type SqlValues []SetField

type DbConfig struct {
	db         *sql.DB
	DriverName string
	Address    string
	User       string
	Password   string
	Port       string
	DbName     string
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

func (S *SqlValues) parseData() (string, []interface{}, error) {
	keys := []string{}
	values := []interface{}{}
	for _, key := range *S {
		keys = append(keys, key.FieldName)
		values = append(values, key.FieldData)
	}
	return strings.Join(keys, ","), values, nil
}

//插入数据
func (config *DbConfig) Insert(table string, datas SqlValues) (id int64, err error) {
	s, v, _ := datas.parseData()
	placeString := fmt.Sprintf("%s", strings.Repeat("?,", len(v)))
	placeString = placeString[:len(placeString)-1]
	sql := "INSERT INTO `" + table + "` (" + s + ") VALUES (" + placeString + ")"

	result, err := config.db.Exec(sql, v...)
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
