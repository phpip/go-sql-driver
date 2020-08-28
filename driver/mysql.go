package driver

import (
	"database/sql"
	"fmt"
	"log"
)
var Db *sql.DB

type SetField struct {
	FieldName string
	FieldData interface{}
}
type SqlValues []SetField

func Insert(table string, datas SqlValues) {
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

	result, e := Db.Exec(sql, fieldValues...)
	if e != nil {
		log.Panicln("user insert error", e.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Panicln("user insert id error", err.Error(), id)
	}
	fmt.Printf("id=%d", id)
}

