package driver

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"strconv"
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
func (config *DbConfig) Update(table string, datas DataStruct, where string) {

}

//获取一条
func (config *DbConfig) GetOne(table, fields string, args ...interface{}) (map[string]interface{}, error) {
	sql := ""
	sql = "SELECT " + fields + " FROM `" + table + "`"
	for _, arg := range args {
		switch arg.(type) {
		case int:
			sql += strconv.Itoa(arg.(int))
		case int64:
			sql += strconv.FormatInt(arg.(int64), 10)
		case string:
			sql += " " + arg.(string)
		case float32:
			sql += "'" + strconv.FormatFloat(float64(arg.(float32)), 'f', -1, 32) + "'"
		case float64:
			sql += "'" + strconv.FormatFloat(arg.(float64), 'f', -1, 64) + "'"
		}
	}
	sql += " LIMIT 0,1"
	fmt.Println(sql)
	rows, err := config.db.Query(sql)
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
func (config *DbConfig) Select(table string, fields string, where string) (map[string]interface{}, error) {
	rows, err := config.db.Query("SELECT " + fields + " FROM `" + table + "` WHERE " + where + " LIMIT 0,1")
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
