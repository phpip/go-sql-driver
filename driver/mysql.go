package driver

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strings"
)

/*type SetField struct {
	FieldName string
	FieldData interface{}
}
type SqlValues []SetField
*/
type DataStruct map[string]interface{}
type DbConfig struct {
	db           *sql.DB
	DriverName   string
	Addr         string
	User         string
	Passwd       string
	Port         string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
}

func (config *DbConfig) Connect() (err error) {
	cfg := mysql.NewConfig()
	cfg.User = config.User
	cfg.Passwd = config.Passwd
	cfg.Net = "tcp"
	cfg.Addr = config.Addr
	cfg.DBName = config.DBName
	dsn := cfg.FormatDSN()
	config.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err := config.db.Ping(); err != nil {
		return err
	}
	maxOpenConns := 0
	if config.MaxOpenConns > 0 {
		maxOpenConns = config.MaxOpenConns
	}
	maxIdleConns := 0
	if config.MaxIdleConns > 0 {
		maxIdleConns = config.MaxIdleConns
	}
	config.db.SetMaxOpenConns(maxOpenConns)
	config.db.SetMaxIdleConns(maxIdleConns)
	return nil
}

func (S *DataStruct) parseData() (string, []interface{}, error) {
	keys := []string{}
	values := []interface{}{}
	for key, value := range *S {
		keys = append(keys, key)
		values = append(values, value)
	}
	return strings.Join(keys, ","), values, nil
}

//添加或者修改数据
func (d *DataStruct) Set(key string, value interface{}) {
	(*d)[key] = value
}

//获取数据
func (d DataStruct) Get(key string) interface{} {
	return d[key]
}

//配合update使用，生成 field=?
func (S *DataStruct) setData() (string, []interface{}, error) {
	keys := []string{}
	values := []interface{}{}
	for key, value := range *S {
		keys = append(keys, key+"=?")
		values = append(values, value)
	}
	return strings.Join(keys, ","), values, nil
}

//插入数据
func (config *DbConfig) Insert(table string, datas DataStruct) (id int64, err error) {
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
func (config *DbConfig) Update(table string, datas DataStruct, where string, args ...interface{}) (num int64, err error) {
	s, v, _ := datas.setData()
	sqlString := "UPDATE `" + table + "` SET " + s
	if where != "" {
		sqlString += " WHERE " + where
	}
	fmt.Println(sqlString)
	for _, value := range args {
		v = append(v, value)
	}
	result, err := config.db.Exec(sqlString, v...)
	num, err = result.RowsAffected()
	return
}

//获取一条
func (config *DbConfig) GetOne(table, fields, where string, args ...interface{}) (map[string]interface{}, error) {
	sqlString := "SELECT " + fields + " FROM `" + table + "`"
	if where != "" {
		sqlString += " WHERE " + where + " LIMIT 0,1"
	}
	fmt.Println(sqlString)
	rows, err := config.db.Query(sqlString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	item := make(map[string]interface{})
	for rows.Next() {
		_ = rows.Scan(cache...)
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
	}
	return item, nil
}

//批量查询，不带分页计算
func (config *DbConfig) Select(table string, fields string, where string, args ...interface{}) ([]map[string]interface{}, error) {
	sqlString := "SELECT " + fields + " FROM `" + table + "`"
	if where != "" {
		sqlString += " WHERE " + where
	}
	fmt.Println(sqlString)
	rows, err := config.db.Query(sqlString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	var results []map[string]interface{}
	item := make(map[string]interface{})
	for rows.Next() {
		_ = rows.Scan(cache...)
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		results = append(results, item)
	}
	return results, nil
}
func (config *DbConfig) Delete(table string, where string, args ...interface{}) (num int64, err error) {
	sqlString := "DELETE FROM `" + table + "`"
	if where != "" {
		sqlString += " WHERE " + where
	}
	fmt.Println(sqlString)
	stmt, err := config.db.Prepare(sqlString)

	result, err := stmt.Exec(args...)
	num, err = result.RowsAffected()
	return
}
func Format2String(datas map[string]interface{}, key string) string {
	if datas[key] == nil {
		return ""
	}
	ba := []byte{}
	data := datas[key].([]uint8)
	for _, b := range data {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
